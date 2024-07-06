# DataVinci

![alt text](68119.png)

DataVinci is a comprehensive data management and visualization tool designed for the developer community. It enables users to visualize data from various sources, generate insights, analyze data with AI models, and receive real-time updates on anomalies.

[![codecov](https://codecov.io/github/doziestar/datavinci/graph/badge.svg?token=PN33XKZTE5)](https://codecov.io/github/doziestar/datavinci)

## Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Getting Started](#getting-started)
- [Development](#development)
- [Deployment](#deployment)
- [Contributing](#contributing)
- [License](#license)

## Features

- Multi-source data integration (PostgreSQL, MongoDB, Cassandra, Elasticsearch, various logs)
- Interactive data visualization with customizable dashboards
- AI-powered data analysis and anomaly detection
- Real-time data processing and alerts
- Cloud resource management and visualization (e.g., Amazon S3)
- Report generation and scheduling
- Collaboration features with version control

## Architecture

DataVinci follows a microservices architecture for scalability and maintainability. Here's a high-level overview of the system:

```mermaid

graph TB
    A[Web UI] --> B[API Gateway]
    B --> C[Authentication Service]
    B --> D[Data Source Service]
    B --> E[Visualization Service]
    B --> F[Report Service]
    B --> G[AI Analysis Service]
    B --> H[Real-time Processing Service]
    D --> I[Data Connectors]
    I --> J[(Various Data Sources)]
    E & F & G & H --> K[Data Processing Engine]
    K --> L[(Data Lake/Warehouse)]
    M[Background Jobs] --> K

```

## Getting Started

### Prerequisites

- Go 1.16+
- Node.js 14+
- Docker and Docker Compose
- Kubernetes cluster (for production deployment)

### Installation

1. Clone the repository:

```bash
git clone https://github.com/doziestar/datavinci.git
cd datavinci
```

2. Install the required dependencies:

```bash
go mod download
cd web && yarn install && cd ..
```

3. Set up the environment variables:

```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Start the development server:

```bash
docker-compose up -d
go run cmd/datavinci/main.go

cd web/src-tauri && cargo tauri dev
```

5. Access the web UI at `http://localhost:3000`.

## Development

DataVinci uses a monorepo structure with Go workspaces (go.work) for backend services and Next.js with Tauri for the frontend.

### Folder Structure

```bash
datavinci/
├── cmd/
│   └── datavinci/
│       └── main.go
├── internal/
│   ├── auth/
│   ├── datasource/
│   ├── visualization/
│   ├── report/
│   ├── ai/
│   └── realtime/
├── pkg/
│   ├── common/
│   └── models/
├── web/
│   ├── components/
│   ├── pages/
│   └── public/
├── deployments/
│   ├── docker/
│   └── k8s/
├── scripts/
├── tests/
├── go.work
├── go.mod
├── go.sum
├── package.json
├── docker-compose.yml
├── Dockerfile
└── README.md
```

### Service Communication

The backend services communicate with each other using gRPC. The API Gateway acts as a reverse proxy for the frontend and forwards requests to the appropriate service.

```mermaid
graph TB
    Client[Client] --> APIGateway[API Gateway]
    subgraph "Service Mesh"
        APIGateway --> Auth[Authentication Service]
        APIGateway --> DataSource[Data Source Service]
        APIGateway --> Visualization[Visualization Service]
        APIGateway --> Report[Report Service]
        APIGateway --> AI[AI Analysis Service]
        APIGateway --> RealTime[Real-time Processing Service]
    end
    Auth -.->|gRPC| DataSource
    DataSource -.->|gRPC| Visualization
    Visualization -.->|gRPC| Report
    DataSource -.->|gRPC| AI
    DataSource -.->|gRPC| RealTime

    MessageBroker[Message Broker] --> DataSource
    MessageBroker --> Visualization
    MessageBroker --> Report
    MessageBroker --> AI
    MessageBroker --> RealTime

    EventStore[(Event Store)] --> MessageBroker

    DataSource --> DB[(Data Sources)]
    RealTime --> DB
```

### Testing

Run the tests with:

```bash
go test ./...
cd web && yarn test && cd ..

// or

go test -v -race -coverprofile=pkg/coverage.txt -covermode=atomic ./internal/auth/...
```

To ensure that the code meets our standards, run the pre-commit hooks:

```bash
pre-commit run --all-files
```

### Linting

Lint the Go code with:

```bash
golangci-lint run
```

Lint the JavaScript code with:

```bash
cd web && yarn lint && cd ..
```

## Deployment

DataVinci can be deployed on any cloud provider or on-premises infrastructure. For production deployments, we recommend using Kubernetes with Helm charts.

### Docker

Build the Docker image with:

```bash
docker build -t datavinci:latest .
```

### Kubernetes

Deploy the application on a Kubernetes cluster with:

```bash
kubectl apply -f deployments/k8s
```

### Helm

Install the Helm chart with:

```bash
helm install datavinci deployments/helm
```

## Contributing

Contributions are welcome! Please read the [contributing guidelines](CONTRIBUTING.md) before submitting a pull request.

## Pre-commit Hooks

We use pre-commit hooks to ensure code quality and consistency. These hooks run automatically before each commit, checking your changes against our coding standards and running various linters.

### Setup

1. Install pre-commit:

   ```bash
   pip install pre-commit
   ```

2. Install the git hook scripts:
   ```bash
   pre-commit install
   ```

### Running pre-commit

The hooks will run automatically on `git commit`. If you want to run the hooks manually (for example, to test them or run them on all files), you can use:

```bash
pre-commit run --all-files
```

### Our pre-commit hooks

We use the following hooks:

- **For Go:**

  - `go-fmt`: Formats Go code
  - `go-vet`: Reports suspicious constructs
  - `go-imports`: Updates import lines
  - `go-cyclo`: Checks function complexity
  - `golangci-lint`: Runs multiple Go linters
  - `go-critic`: Provides extensive code analysis
  - `go-unit-tests`: Runs Go unit tests
  - `go-build`: Checks if the code builds
  - `go-mod-tidy`: Runs `go mod tidy`

- ** ensure that you have the following tools installed:**

  - `golangci-lint`
  - `go-critic`
  - `go-cyclo`
  - `go-unit-tests`
  - `go-build`
  - `go-mod-tidy`

    ```bash
    go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    go install github.com/go-critic/go-critic/cmd/gocritic@latest
    go install github.com/hexdigest/gounit/cmd/gounit@latest
    go install github.com/securego/gosec/v2/cmd/gosec@latest
    ```

- **For TypeScript/JavaScript:**

  - `prettier`: Formats code
  - `eslint`: Lints JavaScript and TypeScript code

- **General:**
  - `trailing-whitespace`: Trims trailing whitespace
  - `end-of-file-fixer`: Ensures files end with a newline
  - `check-yaml`: Checks yaml files for parseable syntax
  - `check-added-large-files`: Prevents giant files from being committed

### Skipping hooks

If you need to bypass the pre-commit hooks (not recommended), you can use:

```bash
git commit -m "Your commit message" --no-verify
```

However, please use this sparingly and ensure your code still meets our standards.

### Updating hooks

To update the pre-commit hooks to the latest versions, run:

```bash
pre-commit autoupdate
```

Then commit the changes to `.pre-commit-config.yaml`.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
