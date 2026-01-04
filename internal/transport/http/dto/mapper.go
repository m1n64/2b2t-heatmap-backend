package dto

import "tbtt-heatmaps-service/internal/config/settings"

func MapSettings(cfg *settings.Settings) SettingsResponse {
	worlds := make(map[string]WorldSettingsResponse)

	for w, s := range cfg.Worlds {
		worlds[w.String()] = WorldSettingsResponse{
			MinZoom:       s.MinZoom,
			MaxZoom:       s.MaxZoom,
			MaxNativeZoom: s.MaxNativeZoom,
			TileSize:      s.TileSize,
			Attribution:   s.Attribution,
		}
	}

	available := make([]string, 0, len(cfg.AvailableWorlds))
	for _, w := range cfg.AvailableWorlds {
		available = append(available, w.String())
	}

	return SettingsResponse{
		AvailableWorlds: available,
		Worlds:          worlds,
	}
}
