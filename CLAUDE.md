# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

imgonna is a full-stack goal-setting application that uses Claude AI to help users process and achieve their goals. The stack consists of:
- **Backend**: Go with Gin framework, GORM, and Auth0 authentication
- **Frontend**: React with TypeScript, Vite, Tailwind CSS, and Zustand
- **Database**: PostgreSQL 15
- **AI**: Anthropic Claude API integration

## Essential Commands

### Development
```bash
# Start all services (recommended for development)
make up

# Run backend locally with DB in Docker
make local-backend

# Stop all services
make down

# Run database migrations
make migrate

# Rollback migrations
make migrate-down
```

### Testing
```bash
# Run all tests
make test

# Backend tests only
make test-backend
# Or: cd backend && go test ./...

# Frontend tests only  
make test-frontend
# Or: cd frontend && npm test
```

### Frontend Development
```bash
cd frontend
npm run dev      # Start dev server (port 5173)
npm run build    # Production build
npm run lint     # Run ESLint
npm run typecheck # Run TypeScript type checking
```

## Architecture

### API Flow
1. Frontend makes requests to backend API at `/api/v1/*`
2. Backend validates Auth0 JWT tokens for authentication
3. Goals endpoint (`POST /api/v1/goals`) processes user input through Claude AI
4. Responses are stored in PostgreSQL and returned to frontend

### Key Backend Patterns
- **Handlers**: Located in `backend/internal/handlers/` - HTTP endpoint logic
- **Models**: Located in `backend/internal/models/` - GORM database models
- **Services**: Located in `backend/internal/services/` - Business logic (e.g., Anthropic integration)
- **Database**: Connection logic in `backend/internal/database/`

### Frontend State Management
- Uses Zustand stores in `frontend/src/store/`:
  - `useAuthStore`: Auth0 authentication state
  - `useGoalStore`: Goal management
  - `useAppStore`: General app state

### Testing Approach
- Backend: Testify framework with mocks for database and external services
- Frontend: Vitest with React Testing Library, MSW for API mocking
- Test files follow `*_test.go` and `*.test.ts(x)` naming conventions

## Environment Configuration

Required environment variables (see `.env.local` for template):
- Database: `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`
- Auth0: `AUTH0_DOMAIN`, `AUTH0_CLIENT_ID`, `AUTH0_AUDIENCE`
- Claude: `CLAUDE_API_KEY`
- Frontend: `REACT_APP_API_URL`, `REACT_APP_AUTH0_*` variables

## Current Implementation

The application currently supports:
- User goal input through a React component (`GoalInput.tsx`)
- Backend processing of goals using Claude AI
- Basic goal storage and retrieval
- Auth0 authentication setup (integration in progress)