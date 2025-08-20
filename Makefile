PROTO_DIR=proto
OUT_DIR=gen/go
PROTOC_GEN_GO=$(shell which protoc-gen-go)
PROTOC_GEN_GO_GRPC=$(shell which protoc-gen-go-grpc)

PROTO_FILES := $(shell find $(PROTO_DIR) -name '*.proto')

all: proto

proto:
	@echo "📦 Compiling protobuf files..."
	@mkdir -p $(OUT_DIR)
	protoc \
	  --proto_path=$(PROTO_DIR) \
	  --go_out=paths=source_relative:$(OUT_DIR) \
	  --go-grpc_out=paths=source_relative:$(OUT_DIR) \
	  $(PROTO_FILES)

clean:
	@echo "🧹 Cleaning generated files..."
	@rm -rf $(OUT_DIR)

grpc:
	go run cmd/grpc/main.go

grpc-hot:
	@echo "🚀 Starting gRPC server with hot reload..."
	air --build.cmd "go build -o bin/grpc cmd/grpc/main.go" --build.bin "./bin/grpc"

lint:
	golangci-lint run ./...

test:
	@echo "Running all tests..."
	go test ./internal/tests/...

test-verbose:
	@echo "Running all tests with verbose output..."
	go test -v ./internal/tests/...

test-coverage:
	@echo "Running all tests with coverage report..."
	go test -v -cover -coverprofile=coverage.out ./internal/tests/...

test-coverage-html:
	@echo "Running all tests and generating HTML coverage report..."
	go test -v -cover -coverprofile=coverage.out ./internal/tests/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-clean:
	@echo "Cleaning test cache and running tests..."
	go clean -testcache && go test -v ./internal/tests/...

security:
	gosec -fmt sarif -out results.sarif -exclude-dir=gen ./... || true

.PHONY: all proto clean grpc grpc-hot lint test test-verbose test-coverage test-coverage-html test-clean security
