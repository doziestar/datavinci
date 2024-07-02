# DataVinci Authentication and Authorization System

## Table of Contents

1. [Overview](#overview)
2. [System Architecture](#system-architecture)
3. [Authentication Flow](#authentication-flow)
4. [Authorization Flow](#authorization-flow)
5. [Service Integration](#service-integration)
6. [User Login Process](#user-login-process)
7. [Token Management](#token-management)
8. [Service-to-Service Communication](#service-to-service-communication)
9. [Security Considerations](#security-considerations)
10. [API Reference](#api-reference)
11. [Configuration](#configuration)
12. [Troubleshooting](#troubleshooting)

## Overview

The DataVinci Authentication and Authorization System is designed to provide secure, scalable, and efficient access control for the DataVinci microservices ecosystem. It uses JWT (JSON Web Tokens) for authentication and a role-based access control (RBAC) system for authorization.

Key features:

- Centralized Authentication Service
- gRPC-based communication
- JWT for stateless authentication
- Role-Based Access Control (RBAC)
- Mutual TLS for service-to-service authentication
- gRPC interceptors for seamless integration

## System Architecture

```mermaid
graph TB
    A[Client] -->|HTTPS| B(API Gateway)
    B -->|gRPC| C{Auth Service}
    B -->|gRPC| D{Data Source Service}
    B -->|gRPC| E{Visualization Service}
    C -->|Validate Token| F[(User DB)]
    C -->|Store Token| G[(Token Cache)]
    H{Authorization Service} -->|Check Permissions| I[(RBAC DB)]
    J{Auth Interceptor} -->|Intercept Requests| D
    J -->|Intercept Requests| E
    J -->|Validate Token| C
    J -->|Check Permissions| H
    K{mTLS} -->|Service-to-Service Auth| D
    K -->|Service-to-Service Auth| E
    K -->|Service-to-Service Auth| C
    K -->|Service-to-Service Auth| H
```

```mermaid
graph TB
    A[Client] -->|HTTPS| B(API Gateway)
    B -->|gRPC| C{Auth Service}
    B -->|gRPC| D{Data Source Service}
    B -->|gRPC| E{Visualization Service}

    C -->|Validate Token| F[(User DB)]
    C -->|Store Token| G[(Token Cache)]

    H{Authorization Service} -->|Check Permissions| I[(RBAC DB)]

    J{Auth Interceptor} -->|Intercept Requests| D
    J -->|Intercept Requests| E
    J -->|Validate Token| C
    J -->|Check Permissions| H

    K{mTLS} -->|Service-to-Service Auth| D
    K -->|Service-to-Service Auth| E
    K -->|Service-to-Service Auth| C
    K -->|Service-to-Service Auth| H

    subgraph "External"
    A
    end

    subgraph "API Gateway"
    B
    end

    subgraph "Microservices"
    C
    D
    E
    H
    end

    subgraph "Security"
    J
    K
    end

    subgraph "Data Stores"
    F
    G
    I
    end
```

```mermaid
graph TB
    A[Client] -->|Request| B(API Gateway)
    B --> C{Auth Middleware}
    C -->|Authenticate| D[AuthManager]
    D -->|Verify Credentials| E[CredentialStore]
    C -->|Authorize| F[PermissionManager]
    F -->|Check Permissions| G[ACL Store]
    C -->|Log Event| H[AuditLogger]
    C -->|If Authorized| I[Data Source Connector]
    I --> J[Data Source]

    subgraph "Authentication"
    D --> K[BasicAuthProvider]
    D --> L[APIKeyAuthProvider]
    D --> M[OAuthProvider]
    end

    subgraph "Authorization"
    F --> N[RoleManager]
    F --> O[PermissionChecker]
    end

    H --> P[(Audit Logs)]
    E --> Q[(Encrypted Credentials)]
    G --> R[(ACL Database)]
```

## Folder Structure

```plaintext
datavinci/
├── api/
│   └── proto/
│       ├── auth.proto
│       ├── authorization.proto
│       └── common.proto
├── cmd/
│   ├── authservice/
│   │   └── main.go
│   ├── authzservice/
│   │   └── main.go
│   └── gateway/
│       └── main.go
├── internal/
│   ├── auth/
│   │   ├── handler/
│   │   │   └── auth_handler.go
│   │   ├── repository/
│   │   │   └── user_repository.go
│   │   ├── service/
│   │   │   └── auth_service.go
│   │   └── token/
│   │       └── jwt.go
│   ├── authz/
│   │   ├── handler/
│   │   │   └── authz_handler.go
│   │   ├── repository/
│   │   │   └── rbac_repository.go
│   │   └── service/
│   │       └── authz_service.go
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   ├── postgres.go
│   │   └── redis.go
│   ├── interceptor/
│   │   ├── auth_interceptor.go
│   │   └── logging_interceptor.go
│   └── model/
│       ├── user.go
│       └── role.go
├── pkg/
│   ├── cryptography/
│   │   └── hash.go
│   ├── logger/
│   │   └── logger.go
│   └── validator/
│       └── validator.go
├── scripts/
│   ├── generate_proto.sh
│   └── generate_certs.sh
├── deployments/
│   ├── docker/
│   │   ├── authservice.Dockerfile
│   │   ├── authzservice.Dockerfile
│   │   └── gateway.Dockerfile
│   └── k8s/
│       ├── authservice.yaml
│       ├── authzservice.yaml
│       └── gateway.yaml
├── tests/
│   ├── integration/
│   │   ├── auth_test.go
│   │   └── authz_test.go
│   └── unit/
│       ├── auth_test.go
│       └── authz_test.go
├── go.mod
├── go.sum
├── Makefile
├── docker-compose.yml
└── README.md
```

## Authentication Flow

1. Client sends login request to API Gateway
2. API Gateway forwards request to Auth Service
3. Auth Service validates credentials and generates JWT
4. JWT is returned to client via API Gateway
5. Client includes JWT in subsequent requests
6. gRPC interceptors validate JWT for each request

## Authorization Flow

1. Client makes request with JWT to a service
2. gRPC interceptor extracts JWT and user information
3. Interceptor calls Authorization Service to check permissions
4. If authorized, the request proceeds; otherwise, it's denied

## Service Integration

Each service in the DataVinci ecosystem integrates with the auth system as follows:

1. Implement gRPC interceptors:

   ```go
   func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
       // Extract JWT from metadata
       token, err := extractToken(ctx)
       if err != nil {
           return nil, status.Errorf(codes.Unauthenticated, "invalid token")
       }

       // Validate JWT
       claims, err := validateToken(token)
       if err != nil {
           return nil, status.Errorf(codes.Unauthenticated, "invalid token")
       }

       // Check permissions
       if !hasPermission(claims, info.FullMethod) {
           return nil, status.Errorf(codes.PermissionDenied, "unauthorized")
       }

       return handler(ctx, req)
   }
   ```

2. Register the interceptor with the gRPC server:

   ```go
   server := grpc.NewServer(
       grpc.UnaryInterceptor(AuthInterceptor),
   )
   ```

3. Use the Authorization Service for fine-grained permission checks within your service methods.

## User Login Process

1. User submits credentials to the client application
2. Client sends login request to API Gateway
3. API Gateway forwards request to Auth Service
4. Auth Service validates credentials against User DB
5. If valid, Auth Service generates a JWT with user claims
6. JWT is returned to the client via API Gateway
7. Client stores JWT securely (e.g., in memory for web apps, secure storage for mobile apps)
8. Client includes JWT in the Authorization header for subsequent requests

Example login request:

```go
func Login(username, password string) (string, error) {
    conn, err := grpc.Dial(authServiceAddress, grpc.WithInsecure())
    if err != nil {
        return "", err
    }
    defer conn.Close()

    client := pb.NewAuthServiceClient(conn)
    resp, err := client.Login(context.Background(), &pb.LoginRequest{
        Username: username,
        Password: password,
    })
    if err != nil {
        return "", err
    }

    return resp.Token, nil
}
```

## Token Management

- JWTs are stateless and contain encoded user information and permissions
- Tokens have an expiration time (e.g., 1 hour)
- We will need to implement token refresh mechanism for long-lived sessions
- Then store active tokens in a distributed cache (e.g., Redis) for quick validation and revocation

Token refresh example:

```go
func RefreshToken(refreshToken string) (string, error) {
    conn, err := grpc.Dial(authServiceAddress, grpc.WithInsecure())
    if err != nil {
        return "", err
    }
    defer conn.Close()

    client := pb.NewAuthServiceClient(conn)
    resp, err := client.RefreshToken(context.Background(), &pb.RefreshTokenRequest{
        RefreshToken: refreshToken,
    })
    if err != nil {
        return "", err
    }

    return resp.NewToken, nil
}
```

## Service-to-Service Communication

Services communicate with each other using mutual TLS (mTLS) for authentication:

1. Generate TLS certificates for each service
2. Configure gRPC servers and clients to use mTLS:

   Server:

   ```go
   cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
   if err != nil {
       log.Fatalf("Failed to load key pair: %v", err)
   }

   server := grpc.NewServer(
       grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
       grpc.UnaryInterceptor(AuthInterceptor),
   )
   ```

   Client:

   ```go
   cert, err := tls.LoadX509KeyPair("client.crt", "client.key")
   if err != nil {
       log.Fatalf("Failed to load client cert: %v", err)
   }

   config := &tls.Config{Certificates: []tls.Certificate{cert}}
   conn, err := grpc.Dial(address, grpc.WithTransportCredentials(credentials.NewTLS(config)))
   ```

## Security Considerations

- We will use strong, frequently rotated encryption keys for JWT signing
- Then we will Implement rate limiting to prevent brute-force attacks
- All external communicationsUse will use HTTPS with TLS
- Proper error handling to avoid information leakage

## API Reference

Auth Service:

- `Login(username, password) returns (JWT token)`
- `ValidateToken(token) returns (claims)`
- `RefreshToken(refresh_token) returns (new_token)`

Authorization Service:

- `CheckPermission(user_id, resource, action) returns (bool)`
- `GetUserRoles(user_id) returns (roles)`
- `AssignRole(user_id, role) returns (success)`

## Configuration

Example configuration file (config.yaml):

```yaml
auth_service:
  address: "localhost:50051"
  jwt_secret: "your-secret-key"
  token_expiration: 3600 # seconds

authorization_service:
  address: "localhost:50052"

database:
  user_db:
    driver: "postgres"
    connection_string: "postgres://user:password@localhost/userdb"
  rbac_db:
    driver: "postgres"
    connection_string: "postgres://user:password@localhost/rbacdb"

cache:
  address: "localhost:6379"
  password: ""
  db: 0
```

## Troubleshooting

Common issues and solutions:

1. "Invalid token" error

   - Check token expiration
   - Verify JWT signing key
   - Ensure clock sync between services

2. "Permission denied" error

   - Check user roles and permissions
   - Verify resource access rules
   - Check for typos in role or permission names

3. Service-to-service communication fails
   - Verify mTLS certificate validity
   - Check network connectivity and firewall rules
   - Ensure correct service addresses in configuration

For more complex issues, enable debug logging and check service logs for detailed error messages.
