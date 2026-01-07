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
	"strings"
	"tbtt-heatmaps-service/internal/config/settings"
	"tbtt-heatmaps-service/internal/enums"
	"tbtt-heatmaps-service/internal/transport/http/dto"
	"tbtt-heatmaps-service/internal/transport/http/httpx"
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

	z, ok := httpx.ParamInt(c, "z")
	if !ok {
		c.Status(http.StatusBadRequest)
		return
	}
	zStr := c.Param("z")

	worldCfg := h.settings.Worlds[world]
	if z < worldCfg.MinZoom || z > worldCfg.MaxNativeZoom {
		c.Status(http.StatusBadRequest)
		return
	}

	x, y := c.Param("x"), c.Param("y")

	log := logging.FromContext(c.Request.Context(), h.logger)

	baseDir := filepath.Join("data", "tiles")
	path, ok := httpx.SafeJoin(
		baseDir,
		world.String(),
		zStr,
		x,
		y+".png",
	)
	if !ok {
		log.Warn("invalid tile path")
		c.Status(http.StatusBadRequest)
		return
	}

	if !strings.HasPrefix(path, baseDir+string(os.PathSeparator)) {
		log.Warn("invalid tile path", zap.String("path", path))
		c.Status(http.StatusBadRequest)
		return
	}

	if val, ok := h.cache.Get(path); ok {
		tile := val.(*CachedTile)

		h.stats.Increment("cache.hit")

		if httpx.HandleETag(c, tile.ETag, httpx.ETagHandler{
			OnHit: func() {
				h.stats.Increment("etag.hit")
			},
			OnMiss: func() {
				h.stats.Increment("etag.miss")
				h.respondTile(c, tile.ETag, tile.Data)
			},
		}) {
			return
		}
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

	if httpx.HandleETag(c, etag, httpx.ETagHandler{
		OnHit: func() {
			h.stats.Increment("etag.hit")
		},
		OnMiss: func() {
			h.stats.Increment("etag.miss")
		},
	}) {
		return
	}

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

	if httpx.HandleETag(c, etag, httpx.ETagHandler{
		OnHit: func() {
			h.stats.Increment("etag.hit")
		},
		OnMiss: func() {
			h.stats.Increment("etag.miss")
			httpx.SetStaticCache(c, 3600)
			c.JSON(http.StatusOK, dto.MapSettings(h.settings))
		},
	}) {
		return
	}
}

func (h *TileHandler) respondTile(
	c *gin.Context,
	etag string,
	data []byte,
) {
	c.Header("ETag", etag)
	httpx.SetStaticCache(c, 86400)
	c.Data(http.StatusOK, "image/png", data)
}

func (h *TileHandler) storeInCache(path string, tile *CachedTile) {
	cost := int64(len(tile.Data) + len(tile.ETag) + 64)

	h.cache.Set(path, tile, cost)
	h.stats.Increment("cache.set")
	h.stats.Gauge("cache.cost_bytes", float64(cost))
	h.stats.Gauge("cache.items", float64(h.cache.Metrics.KeysAdded()))
}
