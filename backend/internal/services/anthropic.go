package services

import (
	"context"
	"fmt"
	"os"
	"strings"
)

type AnthropicServiceInterface interface {
	ProcessGoal(ctx context.Context, goal string) (string, error)
}

type AnthropicService struct {
	apiKey string
}

func NewAnthropicService() *AnthropicService {
	apiKey := os.Getenv("CLAUDE_API_KEY")
	if apiKey == "" {
		// For development, provide a mock response if no API key
		return &AnthropicService{apiKey: "mock"}
	}

	return &AnthropicService{apiKey: apiKey}
}

func (s *AnthropicService) ProcessGoal(ctx context.Context, goal string) (string, error) {
	// For now, provide a mock response if no real API key or if key is "mock"
	if s.apiKey == "" || s.apiKey == "mock" || s.apiKey == "your-claude-api-key" {
		return s.generateMockResponse(goal), nil
	}

	// TODO: Implement real Anthropic API call when API key is provided
	// For now, return mock response even with real API key to avoid issues
	return s.generateMockResponse(goal), nil
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