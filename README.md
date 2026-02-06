# Greeting API

A production-ready Go REST API with structured logging for Filebeat/Kibana integration.

## Features

- Clean Architecture (Domain-Driven Design)
- Structured JSON logging with `log/slog`
- HTTP request/response logging middleware
- Mock generation with `go generate`
- Comprehensive test coverage (100%)
- Filebeat-ready log format
- golangci-lint v2 integration

## Architecture

```
filebeatTest/
├── domain/          # Domain layer: Business logic and entities
├── usecase/         # Application layer: Use cases
├── handler/         # Infrastructure layer: HTTP handlers
├── model/           # Presentation layer: DTOs
├── middleware/      # HTTP middleware (logging, etc.)
├── pkg/logger/      # Structured logging package
├── mock/            # Auto-generated mocks
└── main.go          # Application entry point
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Make (optional)

### Installation

```bash
# Clone the repository
git clone <repository-url>
cd filebeatTest

# Install dependencies
go mod download

# Install development tools
make install-tools
```

### Running the Application

```bash
# Run directly
go run main.go

# Or build and run
make build
./greeting-api
```

The server will start on `http://localhost:8080`.

### API Endpoints

#### GET /greet

Greets a user by name with a time-appropriate greeting.

**Query Parameters:**
- `name` (optional): Name to greet. Defaults to "Guest".

**Time-based Greetings:**
- **5:00-11:59**: Good morning
- **12:00-17:59**: Good afternoon
- **18:00-21:59**: Good evening
- **22:00-4:59**: Good night

**Example:**

```bash
# With name (morning)
curl "http://localhost:8080/greet?name=Alice"
# Response: {"message":"Good morning, Alice! Welcome to our API."}

# Without name (afternoon)
curl "http://localhost:8080/greet"
# Response: {"message":"Good afternoon, Guest! Welcome to our API."}

# Evening
curl "http://localhost:8080/greet?name=Bob"
# Response: {"message":"Good evening, Bob! Welcome to our API."}
```

## Structured Logging

The application outputs structured JSON logs compatible with Filebeat and Elasticsearch/Kibana.

### Log Format

```json
{
  "@timestamp": "2024-01-15T10:30:00.123Z",
  "level": "INFO",
  "msg": "HTTP request",
  "method": "GET",
  "path": "/greet",
  "query": "name=Alice",
  "remote_addr": "127.0.0.1:54321",
  "status_code": 200,
  "duration_ms": 1.234
}
```

### Log Levels

- `INFO`: Normal operations, HTTP requests
- `DEBUG`: Detailed information for debugging
- `WARN`: Warning messages
- `ERROR`: Error conditions

### Viewing Logs

Logs are written to stdout in JSON format:

```bash
# Run application and view logs
./greeting-api

# Save logs to file
./greeting-api > /var/log/greeting-api/app.log 2>&1
```

## Kubernetes Deployment (Minikube)

You can deploy a log collection environment with Filebeat sidecar and Elasticsearch/Kibana to Minikube.

### Prerequisites

- Docker Desktop (WSL2 integration enabled)
- Minikube
- kubectl
- make

### 1. Start Minikube

```bash
minikube start --memory=4096 --cpus=2
```

### 2. Build Docker Image

Build within the Minikube environment:

```bash
eval $(minikube docker-env)
make docker-build
```

### 3. Deploy to Kubernetes

```bash
make k8s-deploy
```

Wait for all Pods to start:

```bash
make k8s-status
```

Once all Pods show `Running` status, you're ready:

```
NAME                             READY   STATUS    RESTARTS   AGE
elasticsearch-xxx                1/1     Running   0          60s
greeting-api-xxx                 2/2     Running   0          60s
kibana-xxx                       1/1     Running   0          60s
```

### 4. Verify API

Access the API from another terminal:

```bash
minikube service greeting-api -n greeting --url
```

Test the API with the displayed URL:

```bash
curl "http://127.0.0.1:<port>/greet?name=test"
```

### 5. View Logs in Kibana

Open Kibana:

```bash
minikube service kibana -n greeting
```

Once the browser opens, follow these steps to view logs:

1. Click the **hamburger menu** in the top left
2. Select **Stack Management** then **Index Patterns**
3. Click **Create index pattern**
4. Enter `greeting-api-*` as the Index pattern
5. Click **Next step**
6. Select `@timestamp` for the Time field
7. Click **Create index pattern**
8. Click the **hamburger menu** then **Discover**
9. Select `greeting-api-*` from the dropdown in the top left

You can now view the greeting-api logs.

### Kubernetes Commands

| Command | Description |
|---------|-------------|
| `make k8s-deploy` | Deploy to Kubernetes |
| `make k8s-status` | Check Pod/Service status |
| `make k8s-logs` | Show greeting-api logs |
| `make k8s-logs-filebeat` | Show Filebeat logs |
| `make k8s-delete` | Delete Kubernetes resources |

### Architecture

```
┌──────────────────────────────────────────────────────────────┐
│ Kubernetes (Minikube)                                        │
│                                                              │
│  ┌─────────────────────────────────────┐                     │
│  │ greeting-api Pod                    │                     │
│  │  ├── greeting-api (main container)  │                     │
│  │  └── filebeat (sidecar) ────────────┼──┐                  │
│  └─────────────────────────────────────┘  │                  │
│                                           ▼                  │
│  ┌─────────────────┐              ┌─────────────────┐        │
│  │     Kibana      │─────────────►│  Elasticsearch  │        │
│  │ (visualization) │              │    (storage)    │        │
│  └─────────────────┘              └─────────────────┘        │
└──────────────────────────────────────────────────────────────┘
```

### Troubleshooting

**If Elasticsearch fails to start:**
```bash
minikube stop
minikube delete
minikube start --memory=6144 --cpus=2
```

**If Kibana shows 500 errors:**
```bash
kubectl rollout restart deployment kibana -n greeting
```

**If Docker image is not found:**
```bash
eval $(minikube docker-env)
make docker-build
```

---

## Filebeat Integration

### Setup

1. Install Filebeat:
```bash
# Ubuntu/Debian
curl -L -O https://artifacts.elastic.co/downloads/beats/filebeat/filebeat-8.11.0-amd64.deb
sudo dpkg -i filebeat-8.11.0-amd64.deb

# macOS
brew tap elastic/tap
brew install elastic/tap/filebeat-full
```

2. Configure Filebeat:
```bash
# Copy the provided configuration
sudo cp filebeat.yml /etc/filebeat/filebeat.yml

# Update paths in filebeat.yml to match your setup
sudo nano /etc/filebeat/filebeat.yml
```

3. Start Filebeat:
```bash
sudo filebeat setup
sudo service filebeat start
```

### Filebeat Configuration

The provided [filebeat.yml](k8s/filebeat.yml) includes:

- JSON log parsing
- Elasticsearch output configuration
- Kibana dashboard setup
- Log enrichment with metadata

**Key settings to customize:**

```yaml
filebeat.inputs:
  - type: log
    paths:
      - /path/to/your/logs/*.log  # Update this path

output.elasticsearch:
  hosts: ["your-elasticsearch-host:9200"]  # Update this

setup.kibana:
  host: "your-kibana-host:5601"  # Update this
```

### Viewing Logs in Kibana

1. Open Kibana: `http://localhost:5601`
2. Go to **Discover** or **Logs**
3. Create index pattern: `greeting-api-*`
4. Filter by fields:
   - `app: greeting-api`
   - `level: ERROR`
   - `method: GET`
   - `status_code: >= 400`

## Development

### Running Tests

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Generate HTML coverage report
make test-coverage-html
```

### Generating Mocks

```bash
# Generate mocks for testing
make generate
```

### Linting

```bash
# Run all linters
make lint

# Run golangci-lint only
make golangci-lint

# Auto-fix issues
make golangci-lint-fix
```

### Code Quality

The project uses golangci-lint v2 with 30+ linters including:

- `errcheck`: Unchecked errors
- `gosec`: Security issues
- `govet`: Suspicious constructs
- `staticcheck`: Static analysis
- `revive`: Code style
- And many more...

## Makefile Commands

```bash
make help                # Show available commands
make test                # Run tests
make test-verbose        # Run tests with verbose output
make test-coverage       # Run tests with coverage
make test-coverage-html  # Generate HTML coverage report
make generate            # Generate mocks
make build               # Build the application
make run                 # Build and run
make dev                 # Run without building
make clean               # Clean build artifacts
make deps                # Download dependencies
make install-tools       # Install development tools
make fmt                 # Format code
make vet                 # Run go vet
make golangci-lint       # Run golangci-lint
make golangci-lint-fix   # Run with auto-fix
make lint                # Run all linters
make all                 # Run all steps
```

## Project Structure Details

### Domain Layer (`domain/`)
Contains core business logic and entities. No external dependencies.

### Use Case Layer (`usecase/`)
Application-specific business rules. Orchestrates domain objects.

### Handler Layer (`handler/`)
HTTP request handlers. Converts HTTP requests to use case calls.

### Model Layer (`model/`)
Data Transfer Objects (DTOs) for API requests and responses.

### Middleware (`middleware/`)
HTTP middleware functions (logging, authentication, etc.).

### Logger Package (`pkg/logger/`)
Structured logging configuration for Filebeat integration.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and linting
5. Submit a pull request

## License

MIT License
