.PHONY: build run test proto

build:
    go build -o bin/visualization-server cmd/server/main.go

run: build
    ./bin/visualization-server

test:
    go test ./...

proto:
    protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        api/proto/visualization.proto