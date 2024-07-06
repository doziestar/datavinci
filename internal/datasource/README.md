# Data Source Module

## Overview

The Data Source module is the foundation of DataVinci, responsible for managing connections to various data sources and providing a unified interface for data retrieval. This module enables the system to interact with different types of databases, file systems, and APIs.

## Key Features

- Connection management for multiple data source types (SQL, NoSQL, File systems, APIs)
- Unified query interface
- Connection pooling and optimization
- Error handling and retry mechanisms

## Structure

```
internal/datasource/
├── api/
│   └── handlers.go
├── internal/
│   ├── service/
│   │   └── service.go
│   └── repository/
│       └── repository.go
├── models/
│   └── models.go
├── connectors/
│   ├── sql.go
│   ├── nosql.go
│   ├── file.go
│   └── api.go
├── config/
│   └── config.go
└── main.go
```

## Technical Design

```mermaid
graph TD
    A[Client] -->|Request| B(API Handler)
    B --> C{Data Source Type}
    C -->|SQL| D[SQL Connector]
    C -->|NoSQL| E[NoSQL Connector]
    C -->|File| F[File Connector]
    C -->|API| G[API Connector]
    D --> H[(SQL Database)]
    E --> I[(NoSQL Database)]
    F --> J[File Storage]
    G --> K[External API]
    L[Connection Pool] --> D
    L --> E
    L --> F
    L --> G
```

```mermaid
sequenceDiagram
    participant OtherService as Other Service
    participant gRPC as grpc/
    participant Service as internal/service
    participant Repository as internal/repository
    participant Connectors as connectors/
    participant ExternalDB as External Databases
    participant Events as events/

    OtherService->>gRPC: gRPC Request
    gRPC->>Service: Process Request
    Service->>Repository: Fetch/Store Data
    Repository->>Connectors: Use Appropriate Connector
    Connectors->>ExternalDB: Query/Update Data
    ExternalDB-->>Connectors: Data Response
    Connectors-->>Repository: Formatted Data
    Repository-->>Service: Processed Data
    Service->>Events: Publish Event (if needed)
    Service-->>gRPC: Prepare gRPC Response
    gRPC-->>OtherService: gRPC Response

    Note over Events: Asynchronous event publishing
    Events->>MessageBroker: Publish Event
    MessageBroker->>OtherService: Consume Event (if subscribed)
```

## Testing

Run tests using:

```
go test ./internal/datasource/...
```
