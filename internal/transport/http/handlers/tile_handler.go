package handlers

import (
	"fmt"
	"github.com/alexcesaro/statsd"
	"github.com/dgraph-io/ristretto"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"tbtt-heatmaps-service/internal/config/settings"
	"tbtt-heatmaps-service/internal/enums"
	"tbtt-heatmaps-service/internal/transport/http/dto"
	"tbtt-heatmaps-service/pkg/logging"
)

type CachedTile struct {
	Data []byte
	ETag string
}

type TileHandler struct {
	logger   *zap.Logger
	cache    *ristretto.Cache
	settings *settings.Settings
	stats    *statsd.Client
}

func NewTileHandler(
	logger *zap.Logger,
	cache *ristretto.Cache,
	settings *settings.Settings,
	stats *statsd.Client,
) *TileHandler {
	return &TileHandler{
		logger:   logger,
		cache:    cache,
		settings: settings,
		stats:    stats,
	}
}

func (h *TileHandler) TileHandle(c *gin.Context) {
	world, ok := enums.ParseWorld(c.Param("world"))
	if !ok {
		c.Status(http.StatusBadRequest)
		return
	}

	zStr := c.Param("z")
	z, err := strconv.Atoi(zStr)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	worldCfg := h.settings.Worlds[world]
	if z < worldCfg.MinZoom || z > worldCfg.MaxNativeZoom {
		c.Status(http.StatusBadRequest)
		return
	}

	x, y := c.Param("x"), c.Param("y")
	if !isInt(x) || !isInt(y) {
		c.Status(http.StatusBadRequest)
		return
	}

	baseDir := filepath.Join("data", "tiles")
	path := filepath.Clean(filepath.Join(
		baseDir,
		world.String(),
		zStr,
		x,
		y+".png",
	))

	log := logging.FromContext(c.Request.Context(), h.logger)

	if !strings.HasPrefix(path, baseDir+string(os.PathSeparator)) {
		log.Warn("invalid tile path", zap.String("path", path))
		c.Status(http.StatusBadRequest)
		return
	}

	if val, ok := h.cache.Get(path); ok {
		tile := val.(*CachedTile)
		h.handleETag(c, tile.ETag, func() {
			h.stats.Increment("cache.hit")
			h.respondTile(c, tile.ETag, tile.Data)
		})
		return
	}

	h.stats.Increment("cache.miss")

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			c.Status(http.StatusNotFound)
			return
		}
		log.Error("failed to open tile", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		log.Error("failed to stat open file", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	etag := fmt.Sprintf(`W/"%x-%x"`, info.Size(), info.ModTime().Unix())

	if c.GetHeader("If-None-Match") == etag {
		h.stats.Increment("etag.hit")
		c.Status(http.StatusNotModified)
		return
	}
	h.stats.Increment("etag.miss")

	data, err := io.ReadAll(f)
	if err != nil {
		log.Error("failed to read tile data", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	tile := &CachedTile{
		Data: data,
		ETag: etag,
	}

	h.storeInCache(path, tile)
	h.respondTile(c, etag, data)
}

func (h *TileHandler) SettingsHandle(c *gin.Context) {
	etag := fmt.Sprintf(`W/"settings-%s"`, h.settings.Version)

	h.handleETag(c, etag, func() {
		resp := dto.MapSettings(h.settings)

		c.Header("Cache-Control", "public, max-age=3600")
		c.JSON(http.StatusOK, resp)
	})
}

func (h *TileHandler) handleETag(
	c *gin.Context,
	etag string,
	onMiss func(),
) {
	if c.GetHeader("If-None-Match") == etag {
		h.stats.Increment("etag.hit")
		c.Status(http.StatusNotModified)
		return
	}

	h.stats.Increment("etag.miss")
	c.Header("ETag", etag)
	onMiss()
}

func (h *TileHandler) respondTile(
	c *gin.Context,
	etag string,
	data []byte,
) {
	c.Header("ETag", etag)
	h.setStaticHeaders(c)
	c.Data(http.StatusOK, "image/png", data)
}

func (h *TileHandler) storeInCache(path string, tile *CachedTile) {
	cost := int64(len(tile.Data) + len(tile.ETag) + 64)

	h.cache.Set(path, tile, cost)
	h.stats.Increment("cache.set")
	h.stats.Gauge("cache.cost_bytes", float64(cost))
	h.stats.Gauge("cache.items", float64(h.cache.Metrics.KeysAdded()))
}

func (h *TileHandler) setStaticHeaders(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=86400")
}

func isInt(v string) bool {
	_, err := strconv.Atoi(v)
	return err == nil
}
