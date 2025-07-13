package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bgoettsch/imgonna/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Anthropic Service
type MockAnthropicService struct {
	mock.Mock
}

func (m *MockAnthropicService) ProcessGoal(ctx context.Context, goal string) (string, *models.MediaRecommendations, error) {
	args := m.Called(ctx, goal)
	if args.Get(1) == nil {
		return args.String(0), nil, args.Error(2)
	}
	return args.String(0), args.Get(1).(*models.MediaRecommendations), args.Error(2)
}

func TestGoalsHandler_CreateGoal_Success(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockAnthropicService)
	handler := NewGoalsHandler(mockService)

	// Mock successful response
	mockService.On("ProcessGoal", mock.Anything, "Learn to play guitar").
		Return("Great goal! Here's how you can start learning guitar...", nil, nil)

	// Create request
	goalRequest := models.GoalRequest{Goal: "Learn to play guitar"}
	jsonData, _ := json.Marshal(goalRequest)

	req, _ := http.NewRequest("POST", "/goals", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Setup Gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Execute
	handler.CreateGoal(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.GoalResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Great goal! Here's how you can start learning guitar...", response.Response)
	assert.Empty(t, response.Error)

	mockService.AssertExpectations(t)
}

func TestGoalsHandler_CreateGoal_WithMediaRecommendations(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockAnthropicService)
	handler := NewGoalsHandler(mockService)

	// Create mock media recommendations
	mockMedia := &models.MediaRecommendations{
		Podcasts: []models.MediaItem{
			{Title: "Test Podcast", Platform: "Spotify", Link: "https://spotify.com", Description: "A test podcast"},
		},
		Streaming: []models.MediaItem{
			{Title: "Test Video", Platform: "YouTube", Description: "A test video"},
		},
		Books: []models.MediaItem{
			{Title: "Test Book", Link: "https://amazon.com", Description: "A test book"},
		},
		Websites: []models.MediaItem{
			{Title: "Test Website", Link: "https://example.com", Description: "A test website"},
		},
	}

	// Mock successful response with media
	mockService.On("ProcessGoal", mock.Anything, "Learn to play guitar").
		Return("Great goal! Here's how you can start learning guitar...", mockMedia, nil)

	// Create request
	goalRequest := models.GoalRequest{Goal: "Learn to play guitar"}
	jsonData, _ := json.Marshal(goalRequest)

	req, _ := http.NewRequest("POST", "/goals", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Setup Gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Execute
	handler.CreateGoal(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.GoalResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Great goal! Here's how you can start learning guitar...", response.Response)
	assert.NotNil(t, response.MediaRecommendations)
	assert.Len(t, response.MediaRecommendations.Podcasts, 1)
	assert.Len(t, response.MediaRecommendations.Streaming, 1)
	assert.Len(t, response.MediaRecommendations.Books, 1)
	assert.Len(t, response.MediaRecommendations.Websites, 1)
	assert.Empty(t, response.Error)

	mockService.AssertExpectations(t)
}

func TestGoalsHandler_CreateGoal_InvalidRequest(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockAnthropicService)
	handler := NewGoalsHandler(mockService)

	// Create invalid request (empty goal)
	goalRequest := models.GoalRequest{Goal: ""}
	jsonData, _ := json.Marshal(goalRequest)

	req, _ := http.NewRequest("POST", "/goals", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Setup Gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Execute
	handler.CreateGoal(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.GoalResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Contains(t, response.Error, "Invalid request")
	assert.Empty(t, response.Response)
}

func TestGoalsHandler_CreateGoal_TooLongGoal(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockAnthropicService)
	handler := NewGoalsHandler(mockService)

	// Create request with goal that's too long
	longGoal := make([]byte, 501)
	for i := range longGoal {
		longGoal[i] = 'a'
	}
	goalRequest := models.GoalRequest{Goal: string(longGoal)}
	jsonData, _ := json.Marshal(goalRequest)

	req, _ := http.NewRequest("POST", "/goals", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Setup Gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Execute
	handler.CreateGoal(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.GoalResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Contains(t, response.Error, "Invalid request")
}