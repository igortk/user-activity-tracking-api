package rest

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"time"
	"user-activity-tracking-api/config"
	"user-activity-tracking-api/rest/handlers"
	"user-activity-tracking-api/rest/middleware"
)

func Run(wg *sync.WaitGroup, httpCfg *config.HttpConfig, stopCh <-chan struct{}) {
	defer wg.Done()

	router := initRouter(httpCfg)
	initApis(router)
	serve(router, httpCfg, stopCh)
}

func initRouter(cfg *config.HttpConfig) *gin.Engine {
	router := gin.New()

	router.Use(middleware.SetupCorsMiddleware(&cfg.CorsConfig))
	router.Use(middleware.TrackMetrics())
	router.Use(middleware.Logger())
	router.Use(gin.Recovery())

	return router
}

func initApis(router *gin.Engine) {
	//api for grafana
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	api := router.Group("/api")
	{
		api.POST("event", handlers.CreateActivityEvent)
		api.GET("events", handlers.GetActivityEventByUserIdDateRange)
	}
}

func serve(router *gin.Engine, httpCfg *config.HttpConfig, stopCh <-chan struct{}) {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", httpCfg.Port),
		Handler: router,
	}

	go func() {
		log.Infof("Starting Rest server on port %d", httpCfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Rest server failed to start: %v", err)
		}
	}()

	<-stopCh

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Info("Rest server forced shutdown: %v", err)
	}
	log.Info("Rest server stopped")
}
