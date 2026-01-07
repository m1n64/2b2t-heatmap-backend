package di

import (
	"github.com/alexcesaro/statsd"
	"github.com/dgraph-io/ristretto"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"os"
	"tbtt-heatmaps-service/internal/config/settings"
	"tbtt-heatmaps-service/internal/highways"
	cache2 "tbtt-heatmaps-service/pkg/cache"
	"tbtt-heatmaps-service/pkg/logging"
	"tbtt-heatmaps-service/pkg/metrics"
	"tbtt-heatmaps-service/pkg/utils"
)

type Dependencies struct {
	Logger          *zap.Logger
	Cache           *ristretto.Cache
	StatsD          *statsd.Client
	SystemCollector *metrics.SystemCollector
	Validator       *validator.Validate
	Settings        *settings.Settings
	HighwayService  *highways.HighwayService
}

func InitDependencies() *Dependencies {
	// Infrastructure
	utils.LoadEnv()

	vectorAddr := os.Getenv("VECTOR_ADDRESS")
	logger, err := logging.InitLogs("tiles-backend-service", vectorAddr)
	if err != nil {
		panic(err)
	}

	statsdAddr := os.Getenv("VMAGENT_ADDRESS")
	metricsClient, err := metrics.NewStatsDClient(statsdAddr)
	if err != nil {
		panic(err)
	}

	cache, err := cache2.NewRistrettoCache()
	if err != nil {
		logger.Fatal("failed to initialize cache", zap.Error(err))
	}

	systemCollector, err := metrics.NewSystemCollector(metricsClient, cache)
	if err != nil {
		panic(err)
	}

	validate := utils.InitValidator()
	cfg, err := settings.LoadFromFile("data/settings.json")
	if err != nil {
		logger.Fatal("failed to load settings.json", zap.Error(err))
		panic(err)
	}

	highwayGenerator := highways.NewHighwayGenerator()
	projector := highways.NewDefaultProjector()
	highwayService := highways.NewHighwayService(highwayGenerator, projector)

	logger.Info("Dependencies initialized")

	return &Dependencies{
		Logger:          logger,
		Cache:           cache,
		StatsD:          metricsClient,
		SystemCollector: systemCollector,
		Validator:       validate,
		Settings:        cfg,
		HighwayService:  highwayService,
	}
}
