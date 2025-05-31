# imgonna

A modern full-stack application built with Go, PostgreSQL, React, and Auth0 authentication.

## Architecture

- **Backend**: Go REST API with Gin framework
- **Database**: PostgreSQL with GORM and migrations
- **Frontend**: React with TypeScript, Vite, and Tailwind CSS
- **State Management**: Zustand
- **Authentication**: Auth0 with JWT validation
- **Testing**: Testify (Go) and Vitest (React)
- **UI Components**: Storybook
- **Containerization**: Docker and Docker Compose

## Prerequisites

- Go 1.24+
- Node.js 20+
- Docker and Docker Compose
- Auth0 account (for authentication setup)
- Claude API key (for AI features)

## Quick Start

### 1. Clone and Setup

```bash
git clone <repository-url>
cd imgonna
cp env.local .env
```

### 2. Configure Environment Variables

Edit `.env` with your actual values:

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=imgonna
DB_USER=postgres
DB_PASSWORD=postgres
DB_SSL_MODE=disable

# Auth0 - Get these from your Auth0 dashboard
AUTH0_DOMAIN=your-domain.auth0.com
AUTH0_CLIENT_ID=your-client-id
AUTH0_CLIENT_SECRET=your-client-secret
AUTH0_AUDIENCE=your-api-audience

# Claude API
CLAUDE_API_KEY=your-claude-api-key

# Server
PORT=8080
ENVIRONMENT=development

# Frontend
REACT_APP_API_URL=http://localhost:8080
REACT_APP_AUTH0_DOMAIN=your-domain.auth0.com
REACT_APP_AUTH0_CLIENT_ID=your-client-id
REACT_APP_AUTH0_AUDIENCE=your-api-audience
```

### 3. Development with Docker

```bash
# Start all services (PostgreSQL, Backend, Frontend)
docker-compose -f docker-compose.dev.yml up

# Or start individual services
docker-compose -f docker-compose.dev.yml up postgres  # Database only
docker-compose -f docker-compose.dev.yml up backend-dev  # Backend only
docker-compose -f docker-compose.dev.yml up frontend-dev  # Frontend only
```

### 4. Manual Development Setup

#### Backend

```bash
cd backend

# Install dependencies
go mod download

# Run database migrations
go run cmd/migrate/main.go up

# Start the server
go run main.go
```

#### Frontend

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev

# Run tests
npm test

# Start Storybook
npm run storybook
```

## Database Migrations

```bash
cd backend

# Apply migrations
go run cmd/migrate/main.go up

# Rollback migrations
go run cmd/migrate/main.go down

# Check migration version
go run cmd/migrate/main.go version

# Force migration to specific version
go run cmd/migrate/main.go force 1
```

## Testing

### Backend Tests

```bash
cd backend
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./internal/models -v
```

### Frontend Tests

```bash
cd frontend

# Run all tests
npm test

# Run tests with UI
npm run test:ui

# Run tests with coverage
npm run test:coverage
```

## API Endpoints

### Health Check
- `GET /health` - Service health status

### API v1
- `GET /api/v1/` - API information
- `GET /api/v1/users` - Get user profile (authenticated)
- `PUT /api/v1/users` - Update user profile (authenticated)
- `GET /api/v1/admin/users` - List all users (admin only)

## Authentication Setup

### Auth0 Configuration

1. Create an Auth0 application
2. Configure callback URLs:
   - `http://localhost:3000/callback` (development)
   - `https://yourdomain.com/callback` (production)
3. Set up social connections (Google, GitHub, etc.)
4. Configure roles and permissions
5. Set up Auth0 Rules/Actions for custom claims

### Environment Variables

Make sure to set the following in your `.env` file:
- `AUTH0_DOMAIN`
- `AUTH0_CLIENT_ID`
- `AUTH0_CLIENT_SECRET`
- `AUTH0_AUDIENCE`

## Deployment

### Production Build

```bash
# Build backend
cd backend
go build -o main .

# Build frontend
cd frontend
npm run build
```

### Docker Production

```bash
# Build and start production containers
docker-compose up --build
```

## Project Structure

```
imgonna/
├── backend/                 # Go API server
│   ├── cmd/
│   │   └── migrate/        # Migration tool
│   ├── internal/
│   │   ├── database/       # Database connection
│   │   ├── models/         # Data models
│   │   ├── handlers/       # HTTP handlers
│   │   ├── middleware/     # HTTP middleware
│   │   └── services/       # Business logic
│   ├── main.go
│   └── Dockerfile
├── frontend/               # React application
│   ├── src/
│   │   ├── components/     # React components
│   │   ├── store/          # Zustand stores
│   │   ├── pages/          # Page components
│   │   ├── hooks/          # Custom hooks
│   │   └── test/           # Test utilities
│   ├── public/
│   ├── package.json
│   └── Dockerfile
├── database/
│   └── migrations/         # SQL migrations
├── docker/                 # Docker configurations
├── docs/                   # Documentation
└── docker-compose.yml      # Docker services
```

## Development Workflow

1. Create feature branch: `git checkout -b feature/your-feature`
2. Make changes and test locally
3. Run linting and tests
4. Commit changes
5. Create pull request

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

[Add your license here]