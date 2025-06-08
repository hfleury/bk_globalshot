# Default target
default: run

# Load environment variables from .env if available
set-env:
    export $(cat .env | grep -v '^#' | xargs)

# Run the application
run:
    go run cmd/main.go

# Build the binary
build:
    docker-compose build --no-cache

# Run tests
test:
    go test ./... -v

# Generate mocks
mock:
    go generate ./...

# Format code
fmt:
    go fmt ./...

# Lint code
lint:
    golangci-lint run ./...

# Clean build files
clean:
    rm -f ./build/app

# Migrate database up
migrate-up:
    migrate -database $(DB_DSN) -path ./migrations up

# Migrate database down
migrate-down:
    migrate -database $(DB_DSN) -path ./migrations down

# Run all migrations and start app
dev: set-env mock migrate-up run

# Run unit tests with coverage
test-coverage:
    go test ./... -coverprofile=coverage.out
    go tool cover -html=coverage.out -o coverage.html
    echo "Coverage report generated at coverage.html"

# Start the app and DB
up:
    docker-compose up -d

down:
    docker-compose down

# View logs
logs:
    docker-compose logs -f

# Shell into running app container
shell:
    docker exec -it bk_globalshot sh

# Rebuild and restart app
restart:
    just build
    just down
    just up