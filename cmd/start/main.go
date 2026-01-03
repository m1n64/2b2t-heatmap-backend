package main

import (
	"golang-service-template/pkg/di"
	"golang-service-template/pkg/utils"
)

var dependencies *di.Dependencies

func init() {
	dependencies = di.InitDependencies()
}

func main() {
	utils.InitMigrations(dependencies.DB)
	dependencies.Logger.Info("Service started")

	select {}
}
