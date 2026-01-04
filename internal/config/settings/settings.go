package settings

import "tbtt-heatmaps-service/internal/enums"

type WorldSettings struct {
	MinZoom       int    `json:"minZoom"`
	MaxZoom       int    `json:"maxZoom"`
	MaxNativeZoom int    `json:"maxNativeZoom"`
	TileSize      int    `json:"tileSize"`
	Attribution   string `json:"attribution,omitempty"`
}

type Settings struct {
	Version         string                        `json:"version"`
	AvailableWorlds []enums.World                 `json:"availableWorlds"`
	Worlds          map[enums.World]WorldSettings `json:"worlds"`
}
