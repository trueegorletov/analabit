# ---- Base Builder ----
FROM golang:1.24-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app
ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org,direct

# We want to install git to be able to pull over ssh from private repos
# and protobuf-dev for code generation
RUN apk add --no-cache git protobuf protobuf-dev ca-certificates && update-ca-certificates

# Install go tools for protobuf
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN go install github.com/asim/go-micro/cmd/protoc-gen-micro/v4@latest

# Copy workspace and module definitions
COPY go.work go.mod go.sum ./
# Copy service modules to proper locations
COPY service/aggregator/go.mod service/aggregator/go.sum service/aggregator/
COPY service/api/go.mod service/api/go.sum service/api/
COPY service/producer/go.mod service/producer/go.sum service/producer/
# Sync workspace build list to modules
RUN go work sync
# Copy entire repository
COPY . .

# Generate protobuf files
RUN protoc --proto_path=./service/aggregator/proto \
    --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    --micro_out=. --micro_opt=paths=source_relative service/aggregator/proto/aggregator.proto

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o aggregator-service ./service/aggregator

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
