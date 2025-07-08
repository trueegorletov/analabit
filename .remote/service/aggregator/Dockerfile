# ---- Base Builder ----
FROM golang:1.24-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# We want to install git to be able to pull over ssh from private repos
# and protobuf-dev for code generation
RUN apk add --no-cache git protobuf protobuf-dev ca-certificates && update-ca-certificates

# Install go tools for protobuf
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN go install github.com/asim/go-micro/cmd/protoc-gen-micro/v4@latest

# Copy go mod and sum files
# Copy aggregator module definition
COPY service/aggregator/go.mod service/aggregator/go.sum ./
# Download dependencies
RUN go mod download
# Copy proto definitions
COPY service/aggregator/proto ./proto
# Generate protobuf files
RUN protoc --proto_path=./proto \
    --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    --micro_out=. --micro_opt=paths=source_relative ./proto/aggregator.proto
# Copy aggregator service code
COPY service/aggregator .

# Generate protobuf files
RUN cd /app && \
    protoc --proto_path=. \
    --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    --micro_out=. --micro_opt=paths=source_relative \
    service/aggregator/proto/aggregator.proto

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o aggregator-service ./service/aggregator/main.go

# ---- Final Stage ----
FROM alpine:latest

# Install poppler-utils for pdftotext
RUN apk add --no-cache poppler-utils

# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/aggregator-service .

# Expose port 8081 to the outside world
EXPOSE 8081

# Switch to our non-root user
USER appuser

# Command to run the executable
CMD ["./aggregator-service"] 