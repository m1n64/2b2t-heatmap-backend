package settings

import (
	"github.com/goccy/go-json"
	"os"
	"tbtt-heatmaps-service/internal/enums"
)

func LoadFromFile(path string) (*Settings, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var tmp struct {
		Version         string                   `json:"version"`
		AvailableWorlds []string                 `json:"availableWorlds"`
		Worlds          map[string]WorldSettings `json:"worlds"`
	}

	if err := json.Unmarshal(raw, &tmp); err != nil {
		return nil, err
	}

	cfg := &Settings{
		Version:         tmp.Version,
		AvailableWorlds: make([]enums.World, 0),
		Worlds:          make(map[enums.World]WorldSettings),
	}

	for _, w := range tmp.AvailableWorlds {
		world, ok := enums.ParseWorld(w)
		if ok {
			cfg.AvailableWorlds = append(cfg.AvailableWorlds, world)
		}
	}

	for k, v := range tmp.Worlds {
		world, ok := enums.ParseWorld(k)
		if ok {
			cfg.Worlds[world] = v
		}
	}

	return cfg, nil
}
