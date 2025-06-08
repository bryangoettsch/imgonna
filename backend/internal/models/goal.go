package models

import "time"

type GoalRequest struct {
	Goal string `json:"goal" binding:"required,min=1,max=500"`
}

type GoalResponse struct {
	Success   bool      `json:"success"`
	Response  string    `json:"response,omitempty"`
	Error     string    `json:"error,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}