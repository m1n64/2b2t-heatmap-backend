package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	httpTransport "tbtt-heatmaps-service/internal/transport/http"
	"tbtt-heatmaps-service/internal/transport/http/middlewares"
	"tbtt-heatmaps-service/pkg/di"
	"time"
)

func main() {
	dependencies := di.InitDependencies()

	r := gin.New()

	defaultGroup := r.Group("")

	var handler http.Handler = r
	if os.Getenv("APP_DEBUG") == "true" {
		fmt.Println("CRITICAL: Global CORS active")
		dependencies.Logger.Debug("debug enabled, adding CORS middleware")
		handler = middlewares.ApplyGlobalCORS(r)
	}

	defaultGroup.Use(gin.Recovery(),
		middlewares.TraceparentMiddleware(),
		middlewares.AccessLogMiddleware(dependencies.Logger, dependencies.StatsD),
	)
	httpTransport.InitRoutes(defaultGroup, dependencies)

	srv := &http.Server{Addr: ":8000", Handler: handler}

	dependencies.SystemCollector.Run(10 * time.Second)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			dependencies.Logger.Fatal(fmt.Sprintf("listen: %s", err.Error()), zap.Error(err))
		}
	}()

	dependencies.Logger.Info(fmt.Sprintf("Service started on 8000 port (%s on dev external)", os.Getenv("SERVICE_PORT")))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)

	dependencies.Logger.Info("Service stopped")
}
