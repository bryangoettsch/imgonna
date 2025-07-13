package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/bgoettsch/imgonna/backend/internal/models"
)

type AnthropicServiceInterface interface {
	ProcessGoal(ctx context.Context, goal string) (string, *models.MediaRecommendations, error)
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

func (s *AnthropicService) ProcessGoal(ctx context.Context, goal string) (string, *models.MediaRecommendations, error) {
	// Use mock response if no real API key or if key is placeholder
	if s.apiKey == "" || s.apiKey == "mock" || s.apiKey == "your-claude-api-key" || s.httpClient == nil {
		slog.Debug("Using mock response", 
			"api_key_empty", s.apiKey == "",
			"api_key_is_mock", s.apiKey == "mock",
			"goal", goal)
		response, media := s.generateMockResponse(goal)
		return response, media, nil
	}

	// Make real API call to Anthropic
	slog.Debug("Making API call to Anthropic", "goal", goal)
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

// APIResponseFormat represents the expected JSON structure from Claude
type APIResponseFormat struct {
	GoalResponse         string                       `json:"goalResponse"`
	MediaRecommendations models.MediaRecommendations `json:"mediaRecommendations"`
}

func (s *AnthropicService) callAnthropicAPI(ctx context.Context, goal string) (string, *models.MediaRecommendations, error) {
	prompt := fmt.Sprintf(`You are a helpful AI assistant that provides guidance and motivation for personal goals, along with relevant media recommendations.

A user has shared this goal: "%s"

Please provide a response in the following JSON format:
{
  "goalResponse": "Your supportive response here (under 200 words)",
  "mediaRecommendations": {
    "podcasts": [
      {"title": "Podcast Name", "link": "spotify/apple podcast link", "platform": "Spotify", "description": "Brief description"},
      ... (5-10 items)
    ],
    "streaming": [
      {"title": "Show/Movie Name", "platform": "Netflix/Hulu/YouTube/etc", "description": "Brief description"},
      ... (5-10 items)
    ],
    "books": [
      {"title": "Book Title by Author", "link": "amazon link if possible", "description": "Brief description"},
      ... (5-10 items)
    ],
    "websites": [
      {"title": "Website Name", "link": "https://...", "description": "Brief description"},
      ... (5-10 items)
    ]
  }
}

For the goal response:
1. Acknowledge their goal positively
2. Offer 2-3 specific, practical steps they can take
3. Include encouragement and motivation
4. Keep it concise and actionable

For media recommendations:
- Provide 5-10 specific, real titles for each category
- For podcasts: Include Spotify or Apple Podcasts links when possible
- For streaming: Indicate the platform (Netflix, Hulu, Prime Video, YouTube, etc.)
- For books: Include Amazon links when possible
- For websites: Provide direct URLs
- Make all recommendations highly relevant to the stated goal

Return ONLY valid JSON, no additional text.`, goal)

	// Prepare request payload
	requestPayload := AnthropicRequest{
		Model:     "claude-3-5-sonnet-20241022",
		MaxTokens: 2500,
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
		return "", nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", s.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	// Make the request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		slog.Error("Anthropic API request failed", 
			"status_code", resp.StatusCode,
			"response_body", string(body),
			"goal", goal)
		return "", nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var apiResponse AnthropicResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return "", nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Extract text from response
	if len(apiResponse.Content) > 0 && apiResponse.Content[0].Text != "" {
		// Parse the JSON response from Claude
		var responseFormat APIResponseFormat
		if err := json.Unmarshal([]byte(apiResponse.Content[0].Text), &responseFormat); err != nil {
			// If JSON parsing fails, return just the text without media
			slog.Warn("Failed to parse JSON response from Claude", 
				"error", err,
				"raw_response", apiResponse.Content[0].Text)
			return apiResponse.Content[0].Text, nil, nil
		}
		slog.Debug("Successfully parsed response with media recommendations",
			"podcasts_count", len(responseFormat.MediaRecommendations.Podcasts),
			"streaming_count", len(responseFormat.MediaRecommendations.Streaming),
			"books_count", len(responseFormat.MediaRecommendations.Books),
			"websites_count", len(responseFormat.MediaRecommendations.Websites))
		return responseFormat.GoalResponse, &responseFormat.MediaRecommendations, nil
	}

	return "", nil, fmt.Errorf("unexpected response format from Claude")
}

func (s *AnthropicService) generateMockResponse(goal string) (string, *models.MediaRecommendations) {
	goalLower := strings.ToLower(goal)
	
	if strings.Contains(goalLower, "learn") && strings.Contains(goalLower, "guitar") {
		response := `🎸 That's a fantastic goal! Learning guitar is incredibly rewarding and a great way to express creativity.

Here are some concrete steps to get started:

1. **Get the basics right**: Start with a beginner-friendly acoustic guitar and learn proper posture and hand positioning. Focus on basic chords like G, C, D, and Em.

2. **Practice regularly**: Even 15-20 minutes daily is better than long, infrequent sessions. Use apps like Yousician or Guitar Tabs, or find a local teacher for structured lessons.

3. **Set small milestones**: Learn one new chord each week, then work on switching between them smoothly. Pick a simple song you love and make it your first goal.

Remember, everyone starts somewhere, and your fingers will feel awkward at first - that's completely normal! Stay patient with yourself and celebrate small victories. You've got this! 🎵`
		
		media := &models.MediaRecommendations{
			Podcasts: []models.MediaItem{
				{Title: "Guitar Lessons and Gear with Tomo Fujita", Platform: "Spotify", Link: "https://open.spotify.com/show/2zRQHGPmBJBGrJvT5vXvAI", Description: "Learn from Berklee professor and John Mayer's guitar teacher"},
				{Title: "The Guitar Hour Podcast", Platform: "Apple Podcasts", Link: "https://podcasts.apple.com/us/podcast/the-guitar-hour-podcast/id1455379351", Description: "Interviews with guitarists and practical playing tips"},
				{Title: "No Guitar Is Safe", Platform: "Spotify", Link: "https://open.spotify.com/show/0iI65jLqzxJCYVlHUe9bXm", Description: "Jude Gold's deep dives into guitar techniques and history"},
			},
			Streaming: []models.MediaItem{
				{Title: "JustinGuitar YouTube Channel", Platform: "YouTube", Description: "Free comprehensive guitar lessons from beginner to advanced"},
				{Title: "Marty Music", Platform: "YouTube", Description: "Popular song tutorials and guitar techniques"},
				{Title: "Paul Davids", Platform: "YouTube", Description: "Music theory and advanced guitar concepts explained clearly"},
				{Title: "GuitarLessons365", Platform: "YouTube", Description: "Song-by-song tutorials for popular music"},
			},
			Books: []models.MediaItem{
				{Title: "Guitar: The First 100 Chords for Guitar by Joseph Alexander", Link: "https://www.amazon.com/dp/1910403334", Description: "Essential chord reference for beginners"},
				{Title: "The Guitar Handbook by Ralph Denyer", Link: "https://www.amazon.com/dp/0679742751", Description: "Comprehensive guide covering all aspects of guitar"},
				{Title: "Hal Leonard Guitar Method, Complete Edition", Link: "https://www.amazon.com/dp/0881881392", Description: "Progressive method book with audio examples"},
			},
			Websites: []models.MediaItem{
				{Title: "Ultimate Guitar", Link: "https://www.ultimate-guitar.com", Description: "Largest collection of guitar tabs and chords"},
				{Title: "JustinGuitar", Link: "https://www.justinguitar.com", Description: "Free structured guitar courses from beginner to advanced"},
				{Title: "Guitar Tricks", Link: "https://www.guitartricks.com", Description: "Comprehensive paid platform with step-by-step lessons"},
				{Title: "Fender Play", Link: "https://www.fender.com/play", Description: "Fender's official learning platform with guided paths"},
			},
		}
		
		return response, media
	}
	
	if strings.Contains(goalLower, "learn") && (strings.Contains(goalLower, "code") || strings.Contains(goalLower, "program")) {
		response := `💻 Excellent choice! Learning to code opens up incredible opportunities and is a skill that keeps growing in value.

Here's your roadmap to success:

1. **Choose your first language**: Start with Python for its beginner-friendly syntax, or JavaScript if you're interested in web development. Both have great learning resources and job prospects.

2. **Build projects, not just tutorials**: After learning basics, create small projects like a calculator, to-do list, or personal website. Building things reinforces learning better than passive consumption.

3. **Join the community**: Use platforms like GitHub to share your work, Stack Overflow for questions, and find local coding meetups or online communities for support and networking.

Start with free resources like freeCodeCamp, Codecademy, or YouTube tutorials. Remember, coding is all about problem-solving - embrace the challenges and celebrate every bug you fix! 🚀`
		
		media := &models.MediaRecommendations{
			Podcasts: []models.MediaItem{
				{Title: "CodeNewbie", Platform: "Spotify", Link: "https://open.spotify.com/show/2T2OwucPOy3dxOCfgyxTMA", Description: "Stories and interviews about people learning to code"},
				{Title: "Syntax - Tasty Web Development Treats", Platform: "Apple Podcasts", Link: "https://podcasts.apple.com/us/podcast/syntax-tasty-web-development-treats/id1253186678", Description: "Full stack web development discussions"},
				{Title: "Talk Python To Me", Platform: "Spotify", Link: "https://open.spotify.com/show/4sYdZnbikqCABkYhHH5fZg", Description: "Python and its ecosystem explored"},
			},
			Streaming: []models.MediaItem{
				{Title: "CS50x - Harvard's Introduction to Computer Science", Platform: "YouTube", Description: "Complete Harvard CS course available free"},
				{Title: "The Coding Train", Platform: "YouTube", Description: "Creative coding tutorials with p5.js and Processing"},
				{Title: "Traversy Media", Platform: "YouTube", Description: "Web development crash courses and project tutorials"},
				{Title: "Programming with Mosh", Platform: "YouTube", Description: "Clear programming tutorials for multiple languages"},
			},
			Books: []models.MediaItem{
				{Title: "Python Crash Course by Eric Matthes", Link: "https://www.amazon.com/dp/1718502702", Description: "Hands-on, project-based introduction to Python"},
				{Title: "JavaScript: The Good Parts by Douglas Crockford", Link: "https://www.amazon.com/dp/0596517742", Description: "Essential JavaScript concepts and best practices"},
				{Title: "Clean Code by Robert C. Martin", Link: "https://www.amazon.com/dp/0132350882", Description: "Writing maintainable, professional code"},
			},
			Websites: []models.MediaItem{
				{Title: "freeCodeCamp", Link: "https://www.freecodecamp.org", Description: "Free coding bootcamp with certifications"},
				{Title: "Codecademy", Link: "https://www.codecademy.com", Description: "Interactive coding lessons in various languages"},
				{Title: "LeetCode", Link: "https://leetcode.com", Description: "Coding challenges and interview preparation"},
				{Title: "MDN Web Docs", Link: "https://developer.mozilla.org", Description: "Comprehensive web development documentation"},
			},
		}
		
		return response, media
	}
	
	// Generic response for other goals
	response := fmt.Sprintf(`That's a wonderful goal! "%s" shows great ambition and self-awareness.

Here's how you can make meaningful progress:

1. **Break it down**: Divide your goal into smaller, specific milestones that you can achieve weekly or monthly. This makes the journey less overwhelming and more manageable.

2. **Create accountability**: Share your goal with friends, family, or join online communities related to your interest. Having others know about your commitment increases your likelihood of success.

3. **Track your progress**: Keep a simple journal or use an app to record your daily actions toward this goal. Seeing your progress builds momentum and motivation.

Remember, every expert was once a beginner. Stay consistent, be patient with yourself, and celebrate small wins along the way. You have everything it takes to achieve this goal! 🌟`, goal)
	
	// Generic media recommendations
	media := &models.MediaRecommendations{
		Podcasts: []models.MediaItem{
			{Title: "The Tim Ferriss Show", Platform: "Spotify", Link: "https://open.spotify.com/show/5qSUyCrk9KR69lEiXbjwXM", Description: "Interviews with world-class performers on habits and routines"},
			{Title: "How I Built This", Platform: "Apple Podcasts", Link: "https://podcasts.apple.com/us/podcast/how-i-built-this/id1150510297", Description: "Stories of entrepreneurs and innovators"},
		},
		Streaming: []models.MediaItem{
			{Title: "TED Talks on Goal Setting", Platform: "YouTube", Description: "Inspiring talks on achievement and goal setting"},
			{Title: "MasterClass", Platform: "MasterClass", Description: "Learn from experts in various fields"},
		},
		Books: []models.MediaItem{
			{Title: "Atomic Habits by James Clear", Link: "https://www.amazon.com/dp/0735211299", Description: "Build good habits and break bad ones"},
			{Title: "The 7 Habits of Highly Effective People by Stephen Covey", Link: "https://www.amazon.com/dp/1982137274", Description: "Timeless principles for personal effectiveness"},
		},
		Websites: []models.MediaItem{
			{Title: "Coursera", Link: "https://www.coursera.org", Description: "Online courses from top universities"},
			{Title: "Khan Academy", Link: "https://www.khanacademy.org", Description: "Free educational resources for all subjects"},
		},
	}
	
	return response, media
}