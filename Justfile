# Default target
default: run

# Load environment variables from .env if available
set-env-var:
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
    echo "Running DB migrations locally..."
    migrate -database "postgres://globalshotuser:globalshotsecret@localhost:5432/globalshotdb?sslmode=disable" -path ./migrations up

migrate-down:
    migrate -database "postgres://globalshotuser:globalshotsecret@globalshotdb:5432/globalshotdb?sslmode=disable" -path ./migrations down

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
logs-docker:
    docker-compose logs -f

# Shell into running app container
shell:
    docker exec -it bk_globalshot sh

# Rebuild and restart app
restart-docker:
    just build
    just down
    just up

set-env:
    eval $(minikube docker-env)

build-image:
    docker build --no-cache --force-rm -t bk_globalshot:latest -f infra/Dockerfile .

apply-namespace:
    kubectl apply -f infra/namespace.yaml

apply-config:
    kubectl apply -f infra/config/app.yaml

apply-secrets:
    kubectl apply -f infra/secrets/app.yaml
    kubectl apply -f infra/secrets/psql.yaml

apply-volume:
    kubectl apply -f infra/volume/psql_data.yaml

apply-db-deployment:
    kubectl apply -f infra/deployment/psql.yaml

apply-app-deployment:
    kubectl apply -f infra/deployment/app.yaml

apply-services:
    kubectl apply -f infra/service/psql.yaml
    kubectl apply -f infra/service/app.yaml

start: set-env build-image apply-namespace apply-config apply-secrets apply-volume apply-db-deployment apply-app-deployment apply-services
    echo "✅ Application deployed to namespace 'globalshot'"
    echo "Use 'kubectl -n globalshot get pods' to check status"
    echo "To access the service: just open-service"

open-service:
    minikube service -n globalshot globalshot-service

logs:
    kubectl -n globalshot logs -l app=globalshot-app --follow

delete:
    kubectl delete ns globalshot
    echo "Namespace 'globalshot' deleted"

restart: delete start
