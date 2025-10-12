package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
	"user-activity-tracking-api/config"
)

func SetupCorsMiddleware(cfg *config.CorsConfig) gin.HandlerFunc {
	corsConfig := cors.DefaultConfig()

	if cfg == nil {
		return cors.New(corsConfig)
	}

	origins := strings.Split(cfg.AllowedOrigins, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}
	corsConfig.AllowOrigins = origins

	if cfg.AllowMethods != "" {
		corsConfig.AllowMethods = strings.Split(cfg.AllowMethods, ",")
	}

	if cfg.AllowHeaders != "" {
		corsConfig.AllowHeaders = strings.Split(cfg.AllowHeaders, ",")
	} else {
		corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	}

	if cfg.MaxAgeHoursCache > 0 {
		corsConfig.MaxAge = time.Duration(cfg.MaxAgeHoursCache) * time.Hour
	}

	corsConfig.AllowCredentials = true

	return cors.New(corsConfig)
}
