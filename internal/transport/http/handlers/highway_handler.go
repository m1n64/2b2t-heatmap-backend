package handlers

import (
	"fmt"
	"github.com/alexcesaro/statsd"
	"github.com/dgraph-io/ristretto"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"tbtt-heatmaps-service/internal/enums"
	"tbtt-heatmaps-service/internal/highways"
	"tbtt-heatmaps-service/internal/transport/http/httpx"
)

type HighwayHandler struct {
	logger         *zap.Logger
	cache          *ristretto.Cache
	highwayService *highways.HighwayService
	stats          *statsd.Client
}

func NewHighwayHandler(logger *zap.Logger, cache *ristretto.Cache, highwayService *highways.HighwayService, stats *statsd.Client) *HighwayHandler {
	return &HighwayHandler{
		logger:         logger,
		cache:          cache,
		highwayService: highwayService,
		stats:          stats,
	}
}

func (h *HighwayHandler) Handle(c *gin.Context) {
	world, ok := enums.ParseWorld(c.Param("world"))
	if !ok {
		c.Status(http.StatusBadRequest)
		return
	}

	cacheKey := fmt.Sprintf("highways:%s:v1", world.String())
	etag := fmt.Sprintf(`W/"v1-highways-geojson-%s"`, world.String())

	if httpx.HandleETag(c, etag, httpx.ETagHandler{
		OnHit: func() {
			// TODO: statsD metrics to constants/enums
			h.stats.Increment("etag.hit")
		},
		OnMiss: func() {
			h.stats.Increment("etag.miss")
		},
	}) {
		return
	}

	if val, ok := h.cache.Get(cacheKey); ok {
		h.stats.Increment("cache.hit")
		httpx.SetStaticCache(c, 86400)
		c.JSON(http.StatusOK, val)
		return
	}

	h.stats.Increment("cache.miss")

	var data highways.FeatureCollection
	switch world {
	case enums.Nether:
		data = h.highwayService.GenerateNether()
	default:
		h.logger.Error(fmt.Sprintf("highway geo-json not implemented for world %s", world.String()), zap.String("world", world.String()))
		c.Status(http.StatusNotImplemented)
		return
	}

	h.cache.Set(cacheKey, data, int64(len(data.Features))*256)

	httpx.SetStaticCache(c, 86400) // 1 week
	c.JSON(200, data)
}
