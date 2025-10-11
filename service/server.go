package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"user-activity-tracking-api/config"
	"user-activity-tracking-api/handlers"
	"user-activity-tracking-api/middleware"
)

func Run(httpCfg *config.HttpConfig) {
	router := gin.New()
	router.Use(middleware.Logger())

	router.Use(gin.Recovery())

	api := router.Group("/api")
	{
		api.POST("event", handlers.CreateActivityEvent)
		api.GET("events", handlers.GetActivityEventByUserIdDateRange)
	}

	if err := router.Run(fmt.Sprintf(":%d", httpCfg.Port)); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
