package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
	"user-activity-tracking-api/models"
	"user-activity-tracking-api/service/database/repositories"
	"user-activity-tracking-api/utils"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func CreateActivityEvent(eventRepo *repositories.EventsRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var event models.Event

		if err := ctx.ShouldBindJSON(&event); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			log.Errorf("invalid request body: %v", err)
			return
		}

		if msg := utils.GenerateErrorMessage(event, validate); msg != "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": msg})
			return
		}

		c, cancel := context.WithTimeout(ctx.Request.Context(), 60*time.Second)
		defer cancel()

		if err := eventRepo.CreateEvent(c, &event); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not crete an event"})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"createdAt": time.Now().UTC()})
	}
}

func GetActivityEventByUserIdDateRange(eventRepo *repositories.EventsRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req models.GetActivityEventByUserIdDateRangeRequest

		if err := ctx.ShouldBindQuery(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if msg := utils.GenerateErrorMessage(req, validate); msg != "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": msg})
			return
		}
		var events []models.Event

		c, cancel := context.WithTimeout(ctx.Request.Context(), 60*time.Second)
		defer cancel()

		events, err := eventRepo.GetEventsByUserIdAndDateRange(c, req.UserID, req.Limit, req.Offset, req.FromEventActionTimestamp, req.ToEventActionTimestamp)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something was wrong in db"})
			return
		}

		ctx.JSON(http.StatusOK, events)
	}
}
