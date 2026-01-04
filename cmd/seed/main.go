package main

import (
	"tbtt-heatmaps-service/pkg/di"
)

var dependencies *di.Dependencies

func init() {
	dependencies = di.InitDependencies()
}

func main() {
	dependencies.Logger.Info("Successfully seeded")
}
