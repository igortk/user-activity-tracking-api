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
	"user-activity-tracking-api/internal/configs"
	"user-activity-tracking-api/internal/rest/handlers"
	middleware2 "user-activity-tracking-api/internal/rest/middleware"
	"user-activity-tracking-api/internal/service/database"
	"user-activity-tracking-api/internal/service/database/repositories"
)

type Server struct {
	httpCfg *configs.HttpConfig

	eventRepo repositories.EventsRepository
}

func NewServer(cfg *configs.Config, dbCl *database.Client) *Server {
	return &Server{
		httpCfg:   &cfg.HttpConfig,
		eventRepo: repositories.NewEventsRepository(dbCl.GetDb()),
	}
}

func (s *Server) Run(wg *sync.WaitGroup, stopCh <-chan struct{}) {
	defer wg.Done()

	router := s.initRouter(s.httpCfg)
	s.initApis(router)
	s.serve(router, s.httpCfg, stopCh)
}

func (s *Server) initRouter(cfg *configs.HttpConfig) *gin.Engine {
	router := gin.New()

	router.Use(middleware2.MaxBodySize(1 << 20))
	router.Use(middleware2.SetupCorsMiddleware(&cfg.CorsConfig))
	router.Use(middleware2.TrackMetrics())
	router.Use(middleware2.Logger())
	router.Use(gin.Recovery())

	return router
}

func (s *Server) initApis(router *gin.Engine) {
	//api for grafana
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	api := router.Group("/api")
	{
		api.POST("event", handlers.CreateActivityEvent(s.eventRepo))
		api.GET("events", handlers.GetActivityEventByUserIdDateRange(s.eventRepo))
	}
}

func (s *Server) serve(router *gin.Engine, httpCfg *configs.HttpConfig, stopCh <-chan struct{}) {
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", httpCfg.Port),
		Handler:           router,
		ReadTimeout:       time.Duration(httpCfg.ReadTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(httpCfg.ReadHeaderTimeout) * time.Second,
		WriteTimeout:      time.Duration(httpCfg.WriteTimeout) * time.Second,
		IdleTimeout:       time.Duration(httpCfg.IdleTimeout) * time.Second,
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
