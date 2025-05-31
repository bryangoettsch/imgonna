.PHONY: up down local-backend stop-backend test test-backend test-frontend migrate migrate-down clean help

# Default target
help:
	@echo "Available targets:"
	@echo "  up              - Start all services (postgres, backend, frontend) in Docker"
	@echo "  down            - Stop all Docker services"
	@echo "  local-backend   - Stop Docker backend and run locally for debugging"
	@echo "  stop-backend    - Stop only the backend container"
	@echo "  test            - Run all tests (backend + frontend)"
	@echo "  test-backend    - Run backend tests"
	@echo "  test-frontend   - Run frontend tests"
	@echo "  migrate         - Run database migrations"
	@echo "  migrate-down    - Rollback database migrations"
	@echo "  clean           - Clean up Docker volumes and containers"

# Start all services in Docker
up:
	@echo "🚀 Starting all services..."
	docker compose -f docker-compose.dev.yml up

# Stop all services
down:
	@echo "🛑 Stopping all services..."
	docker compose -f docker-compose.dev.yml down

# Stop backend container and run locally
local-backend: stop-backend
	@echo "🔧 Starting local backend development..."
	@echo "📦 Keeping PostgreSQL running in Docker..."
	docker compose -f docker-compose.dev.yml up postgres -d
	@echo "🏃 Run the following commands in separate terminals:"
	@echo "   cd backend && go run main.go"
	@echo "   OR for hot reload: cd backend && air"
	@echo "   OR debug in VS Code with F5"

# Stop only the backend container
stop-backend:
	@echo "⏹️  Stopping backend container..."
	docker compose -f docker-compose.dev.yml stop backend-dev

# Run all tests
test: test-backend test-frontend

# Run backend tests
test-backend:
	@echo "🧪 Running backend tests..."
	cd backend && go test ./...

# Run frontend tests
test-frontend:
	@echo "🧪 Running frontend tests..."
	cd frontend && npm test -- --run

# Run database migrations
migrate:
	@echo "📊 Running database migrations..."
	cd backend && go run cmd/migrate/main.go up

# Rollback database migrations
migrate-down:
	@echo "📊 Rolling back database migrations..."
	cd backend && go run cmd/migrate/main.go down

# Clean up Docker resources
clean:
	@echo "🧹 Cleaning up Docker resources..."
	docker compose -f docker-compose.dev.yml down -v
	docker system prune -f