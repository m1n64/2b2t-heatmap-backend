package dto

type WorldSettingsResponse struct {
	Bounds        [][]int `json:"bounds"`
	MinZoom       int     `json:"min_zoom"`
	MaxZoom       int     `json:"max_zoom"`
	MinNativeZoom int     `json:"min_native_zoom"`
	MaxNativeZoom int     `json:"max_native_zoom"`
	ZoomDelta     int     `json:"zoom_delta"`
	ZoomSnap      int     `json:"zoom_snap"`
	TileSize      int     `json:"tile_size"`
	NoWrap        bool    `json:"no_wrap"`
	Attribution   string  `json:"attribution,omitempty"`
}

type SettingsResponse struct {
	AvailableWorlds []string                         `json:"available_worlds"`
	Worlds          map[string]WorldSettingsResponse `json:"worlds"`
}
