package handlers

import (
	"net/http"
	"time"

	"github.com/bgoettsch/imgonna/backend/internal/models"
	"github.com/bgoettsch/imgonna/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type GoalsHandler struct {
	anthropicService services.AnthropicServiceInterface
}

func NewGoalsHandler(anthropicService services.AnthropicServiceInterface) *GoalsHandler {
	return &GoalsHandler{
		anthropicService: anthropicService,
	}
}

func (h *GoalsHandler) CreateGoal(c *gin.Context) {
	var req models.GoalRequest

	// Bind and validate the request
	if err := c.ShouldBindJSON(&req); err != nil {
		response := models.GoalResponse{
			Success:   false,
			Error:     "Invalid request: " + err.Error(),
			Timestamp: time.Now(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Process the goal with Anthropic API
	aiResponse, err := h.anthropicService.ProcessGoal(c.Request.Context(), req.Goal)
	if err != nil {
		response := models.GoalResponse{
			Success:   false,
			Error:     "Failed to process goal: " + err.Error(),
			Timestamp: time.Now(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// Return successful response
	response := models.GoalResponse{
		Success:   true,
		Response:  aiResponse,
		Timestamp: time.Now(),
	}

	c.JSON(http.StatusOK, response)
}