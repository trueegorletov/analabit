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

# Copy entire repository first to get proper module structure
COPY . .

# Generate protobuf files
RUN protoc --proto_path=./service/producer/proto \
    --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    --micro_out=. --micro_opt=paths=source_relative service/producer/proto/producer.proto

# Change to producer service directory to build with correct module context
WORKDIR /app/service/producer

# Build the Go app for producer service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o producer-service .

# ---- Final Stage ----
FROM alpine:latest

# Install poppler-utils for pdftotext
RUN apk add --no-cache poppler-utils

# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/service/producer/producer-service .

# Expose port 8080 to the outside world
EXPOSE 8080

# Switch to our non-root user
USER appuser

# Command to run the executable
CMD ["./producer-service"]
