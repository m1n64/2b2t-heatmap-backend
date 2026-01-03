package di

import (
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"golang-service-template/pkg/utils"
	"gorm.io/gorm"
	"os"
)

type Dependencies struct {
	Logger    *zap.Logger
	Redis     *redis.Client
	DB        *gorm.DB
	Validator *validator.Validate
}

func InitDependencies() *Dependencies {
	// Infrastructure
	logger := utils.InitLogs()
	utils.LoadEnv()
	redisConn := utils.CreateRedisConn(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	dbConn := utils.InitDBConnection(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	utils.InitMigrations(dbConn)

	validate := utils.InitValidator()

	return &Dependencies{
		Logger:    logger,
		Redis:     redisConn,
		DB:        dbConn,
		Validator: validate,
	}
}
