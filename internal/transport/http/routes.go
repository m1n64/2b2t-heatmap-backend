package http

import (
	"github.com/gin-gonic/gin"
	handlers2 "tbtt-heatmaps-service/internal/transport/http/handlers"
	"tbtt-heatmaps-service/pkg/di"
)

type handlers struct {
	tileHandler    *handlers2.TileHandler
	highwayHandler *handlers2.HighwayHandler
}

func InitRoutes(core *gin.RouterGroup, dependencies *di.Dependencies) {
	handler := initHandlers(dependencies)

	core.GET("/:world/:z/:x/:y", handler.tileHandler.TileHandle)
	core.GET("/settings", handler.tileHandler.SettingsHandle)
	core.GET("/:world/highways/geo-json", handler.highwayHandler.Handle)

	dependencies.Logger.Info("HTTP routes initialized")
}

func initHandlers(dependencies *di.Dependencies) *handlers {
	return &handlers{
		tileHandler:    handlers2.NewTileHandler(dependencies.Logger, dependencies.Cache, dependencies.Settings, dependencies.StatsD),
		highwayHandler: handlers2.NewHighwayHandler(dependencies.Logger, dependencies.Cache, dependencies.HighwayService, dependencies.StatsD),
	}
}
