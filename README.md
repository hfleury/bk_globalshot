# bk_globalshot


## Local Dev
### Set Namespace
Set your current context to the namespace so you don’t have to specify it every time:

```shell
    kubectl config set-context --current --namespace=globalshot
```

## Development Commands
You can view all available `just` commands at any time by running:

```sh
just --list
```

This will display a list of all recipes (commands) defined in the `Justfile`, along with their descriptions (if provided as comments above each command).  
Adding comments above each recipe in the `Justfile` helps make the output of `just --list` more informative.

This project uses [just](https://github.com/casey/just) for task automation. Below are the available commands:

| Command                | Description                                                                                   |
|------------------------|-----------------------------------------------------------------------------------------------|
| `just run`             | Run the application locally (`go run cmd/main.go`).                                           |
| `just build`           | Build the Docker image using `docker-compose build --no-cache`.                              |
| `just test`            | Run all Go tests with verbose output.                                                        |
| `just mock`            | Generate Go mocks using `go generate ./...`.                                                 |
| `just fmt`             | Format all Go code using `go fmt ./...`.                                                     |
| `just lint`            | Lint the codebase using `golangci-lint`.                                                     |
| `just clean`           | Remove the built binary (`./build/app`).                                                     |
| `just migrate-up`      | Run database migrations up using the `migrate` tool.                                         |
| `just migrate-down`    | Roll back database migrations using the `migrate` tool.                                      |
| `just dev`             | Run `set-env`, generate mocks, migrate up, and run the app (for local development).          |
| `just test-coverage`   | Run tests with coverage and generate an HTML report.                                         |
| `just up`              | Start the app and database using Docker Compose.                                             |
| `just down`            | Stop the app and database containers.                                                        |
| `just logs-docker`     | View Docker Compose logs.                                                                    |
| `just shell`           | Open a shell in the running app container.                                                   |
| `just restart-docker`  | Rebuild and restart the Docker containers.                                                   |
| `just set-env`         | Set Docker environment variables for Minikube.                                               |
| `just build-image`     | Build the Docker image for Kubernetes deployment.                                            |
| `just apply-namespace` | Apply the Kubernetes namespace manifest.                                                     |
| `just apply-config`    | Apply the Kubernetes ConfigMap for the app.                                                  |
| `just apply-secrets`   | Apply the Kubernetes secrets for the app and database.                                       |
| `just apply-volume`    | Apply the PersistentVolumeClaim for the database.                                            |
| `just apply-db-deployment` | Deploy the PostgreSQL database to Kubernetes.                                            |
| `just apply-app-deployment` | Deploy the application to Kubernetes.                                                   |
| `just apply-services`  | Apply the Kubernetes Service manifests for app and database.                                 |
| `just start`           | Build image, apply all manifests, and deploy everything to Kubernetes.                       |
| `just open-service`    | Open the app service in your browser via Minikube.                                           |
| `just logs`            | View logs from the running app pod in Kubernetes.                                            |
| `just delete`          | Delete the entire Kubernetes namespace and all resources.                                    |
| `just restart`         | Delete and redeploy everything in Kubernetes.                                                |

> **Tip:** You can run any command with `just <command>`. For example:  
> `just build-image`
