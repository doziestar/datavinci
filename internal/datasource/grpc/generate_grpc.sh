#!/bin/bash

# Set the directory where your .proto file is located
PROTO_DIR="./internal/datasource/grpc/proto"

# Set the output directory for generated Go files
GO_OUT_DIR="./internal/datasource/grpc"

CLIENT_OUT_DIR="./internal/datasource/grpc/client"

# Ensure the output directory exists
mkdir -p $CLIENT_OUT_DIR

# Generate Go code for client
protoc --proto_path=${PROTO_DIR} \
       --go_out=${CLIENT_OUT_DIR} --go_opt=paths=source_relative \
       --go-grpc_out=${CLIENT_OUT_DIR} --go-grpc_opt=paths=source_relative \
       ${PROTO_DIR}/connector.proto

# Generate Go code
protoc --proto_path=${PROTO_DIR} \
       --go_out=${GO_OUT_DIR} --go_opt=paths=source_relative \
       --go-grpc_out=${GO_OUT_DIR} --go-grpc_opt=paths=source_relative \
       ${PROTO_DIR}/connector.proto

echo "DataSource gRPC code generation completed."