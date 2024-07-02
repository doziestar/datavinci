#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

# Get the absolute path to the project root
PROJECT_ROOT=$(git rev-parse --show-toplevel)

# Directory containing the proto files
PROTO_DIR="$PROJECT_ROOT/internal/auth/api/proto"

# Directory to output the generated Go code
GO_OUT_DIR="$PROJECT_ROOT/internal/auth/pb"

# Create the output directory if it doesn't exist
mkdir -p $GO_OUT_DIR

# Generate Go code from proto files
protoc \
    --proto_path=$PROTO_DIR \
    --go_out=$GO_OUT_DIR --go_opt=paths=source_relative \
    --go-grpc_out=$GO_OUT_DIR --go-grpc_opt=paths=source_relative \
    $PROTO_DIR/*.proto

echo "Proto files generated successfully in $GO_OUT_DIR"

# Optionally, run go mod tidy to ensure all dependencies are properly managed
cd $PROJECT_ROOT && go mod tidy

echo "Proto generation complete!"