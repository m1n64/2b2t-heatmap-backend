package settings

import "tbtt-heatmaps-service/internal/enums"

type WorldSettings struct {
	Bounds        [][]int `json:"bounds"`
	MinZoom       int     `json:"minZoom"`
	MaxZoom       int     `json:"maxZoom"`
	MinNativeZoom int     `json:"minNativeZoom"`
	MaxNativeZoom int     `json:"maxNativeZoom"`
	ZoomDelta     int     `json:"zoomDelta"`
	ZoomSnap      int     `json:"zoomSnap"`
	TileSize      int     `json:"tileSize"`
	NoWrap        bool    `json:"noWrap"`
	Attribution   string  `json:"attribution,omitempty"`
}

type Settings struct {
	Version         string                        `json:"version"`
	AvailableWorlds []enums.World                 `json:"availableWorlds"`
	Worlds          map[enums.World]WorldSettings `json:"worlds"`
}
