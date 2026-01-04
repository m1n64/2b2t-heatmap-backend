package dto

type WorldSettingsResponse struct {
	MinZoom       int    `json:"min_zoom"`
	MaxZoom       int    `json:"max_zoom"`
	MaxNativeZoom int    `json:"max_native_zoom"`
	TileSize      int    `json:"tile_size"`
	Attribution   string `json:"attribution,omitempty"`
}

type SettingsResponse struct {
	AvailableWorlds []string                         `json:"available_worlds"`
	Worlds          map[string]WorldSettingsResponse `json:"worlds"`
}
