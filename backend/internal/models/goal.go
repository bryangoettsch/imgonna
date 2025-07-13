package models

import "time"

type GoalRequest struct {
	Goal string `json:"goal" binding:"required,min=1,max=500"`
}

type MediaItem struct {
	Title       string `json:"title"`
	Link        string `json:"link,omitempty"`
	Platform    string `json:"platform,omitempty"`
	Description string `json:"description,omitempty"`
}

type MediaRecommendations struct {
	Podcasts  []MediaItem `json:"podcasts"`
	Streaming []MediaItem `json:"streaming"`
	Books     []MediaItem `json:"books"`
	Websites  []MediaItem `json:"websites"`
}

type GoalResponse struct {
	Success              bool                  `json:"success"`
	Response             string                `json:"response,omitempty"`
	MediaRecommendations *MediaRecommendations `json:"mediaRecommendations,omitempty"`
	Error                string                `json:"error,omitempty"`
	Timestamp            time.Time             `json:"timestamp"`
}