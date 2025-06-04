package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type AnthropicServiceInterface interface {
	ProcessGoal(ctx context.Context, goal string) (string, error)
}

type AnthropicService struct {
	apiKey string
	httpClient *http.Client
}

func NewAnthropicService() *AnthropicService {
	apiKey := os.Getenv("CLAUDE_API_KEY")
	if apiKey == "" {
		// For development, provide a mock response if no API key
		return &AnthropicService{apiKey: "mock", httpClient: nil}
	}

	// Initialize HTTP client for real API calls
	httpClient := &http.Client{}
	return &AnthropicService{apiKey: apiKey, httpClient: httpClient}
}

func (s *AnthropicService) ProcessGoal(ctx context.Context, goal string) (string, error) {
	// Use mock response if no real API key or if key is placeholder
	if s.apiKey == "" || s.apiKey == "mock" || s.apiKey == "your-claude-api-key" || s.httpClient == nil {
		return s.generateMockResponse(goal), nil
	}

	// Make real API call to Anthropic
	return s.callAnthropicAPI(ctx, goal)
}

// AnthropicRequest represents the request payload for Anthropic API
type AnthropicRequest struct {
	Model     string `json:"model"`
	MaxTokens int    `json:"max_tokens"`
	Messages  []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}

// AnthropicResponse represents the response from Anthropic API
type AnthropicResponse struct {
	Content []struct {
		Text string `json:"text"`
		Type string `json:"type"`
	} `json:"content"`
	ID           string `json:"id"`
	Model        string `json:"model"`
	Role         string `json:"role"`
	StopReason   string `json:"stop_reason"`
	StopSequence string `json:"stop_sequence"`
	Type         string `json:"type"`
	Usage        struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

func (s *AnthropicService) callAnthropicAPI(ctx context.Context, goal string) (string, error) {
	prompt := fmt.Sprintf(`You are a helpful AI assistant that provides guidance and motivation for personal goals. 

A user has shared this goal: "%s"

Please provide a supportive, actionable response that:
1. Acknowledges their goal positively
2. Offers 2-3 specific, practical steps they can take to work toward this goal
3. Includes encouragement and motivation
4. Keeps the response concise (under 200 words)

Be warm, encouraging, and focus on actionable advice.`, goal)

	// Prepare request payload
	requestPayload := AnthropicRequest{
		Model:     "claude-3-5-sonnet-20241022",
		MaxTokens: 250,
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	// Marshal request to JSON
	jsonData, err := json.Marshal(requestPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", s.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	// Make the request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var apiResponse AnthropicResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	// Extract text from response
	if len(apiResponse.Content) > 0 && apiResponse.Content[0].Text != "" {
		return apiResponse.Content[0].Text, nil
	}

	return "", fmt.Errorf("unexpected response format from Claude")
}

func (s *AnthropicService) generateMockResponse(goal string) string {
	goalLower := strings.ToLower(goal)
	
	if strings.Contains(goalLower, "learn") && strings.Contains(goalLower, "guitar") {
		return `🎸 That's a fantastic goal! Learning guitar is incredibly rewarding and a great way to express creativity.

Here are some concrete steps to get started:

1. **Get the basics right**: Start with a beginner-friendly acoustic guitar and learn proper posture and hand positioning. Focus on basic chords like G, C, D, and Em.

2. **Practice regularly**: Even 15-20 minutes daily is better than long, infrequent sessions. Use apps like Yousician or Guitar Tabs, or find a local teacher for structured lessons.

3. **Set small milestones**: Learn one new chord each week, then work on switching between them smoothly. Pick a simple song you love and make it your first goal.

Remember, everyone starts somewhere, and your fingers will feel awkward at first - that's completely normal! Stay patient with yourself and celebrate small victories. You've got this! 🎵`
	}
	
	if strings.Contains(goalLower, "learn") && (strings.Contains(goalLower, "code") || strings.Contains(goalLower, "program")) {
		return `💻 Excellent choice! Learning to code opens up incredible opportunities and is a skill that keeps growing in value.

Here's your roadmap to success:

1. **Choose your first language**: Start with Python for its beginner-friendly syntax, or JavaScript if you're interested in web development. Both have great learning resources and job prospects.

2. **Build projects, not just tutorials**: After learning basics, create small projects like a calculator, to-do list, or personal website. Building things reinforces learning better than passive consumption.

3. **Join the community**: Use platforms like GitHub to share your work, Stack Overflow for questions, and find local coding meetups or online communities for support and networking.

Start with free resources like freeCodeCamp, Codecademy, or YouTube tutorials. Remember, coding is all about problem-solving - embrace the challenges and celebrate every bug you fix! 🚀`
	}
	
	// Generic response for other goals
	return fmt.Sprintf(`That's a wonderful goal! "%s" shows great ambition and self-awareness.

Here's how you can make meaningful progress:

1. **Break it down**: Divide your goal into smaller, specific milestones that you can achieve weekly or monthly. This makes the journey less overwhelming and more manageable.

2. **Create accountability**: Share your goal with friends, family, or join online communities related to your interest. Having others know about your commitment increases your likelihood of success.

3. **Track your progress**: Keep a simple journal or use an app to record your daily actions toward this goal. Seeing your progress builds momentum and motivation.

Remember, every expert was once a beginner. Stay consistent, be patient with yourself, and celebrate small wins along the way. You have everything it takes to achieve this goal! 🌟`, goal)
}