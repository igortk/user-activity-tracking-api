package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"user-activity-tracking-api/models"
	"user-activity-tracking-api/service/database"
	"user-activity-tracking-api/service/database/repositories"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func CreateActivityEvent(ctx *gin.Context) {
	var event models.Event

	if err := ctx.ShouldBindJSON(&event); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate.Struct(event); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": validationErrors.Error()})
		return
	}

	if err := database.Session.Create(&event).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not crete an event"})
		return
	}

	ctx.Status(http.StatusCreated)
}

func GetActivityEventByUserIdDateRange(ctx *gin.Context) {
	var req models.GetActivityEventByUserIdDateRangeRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": validationErrors.Error()})
		return
	}
	var events []models.Event

	events, err := repositories.GetEventsByUserIdAndDateRange(req.UserID, req.FromEventActionTimestamp, req.ToEventActionTimestamp)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something was wrong in db"})
		return
	}

	ctx.JSON(http.StatusOK, events)
}
