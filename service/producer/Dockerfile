# ---- Base Builder ----
FROM golang:1.21-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# We want to install git to be able to pull over ssh from private repos
# and protobuf-dev for code generation
RUN apk add --no-cache git protobuf-dev

# Install go tools for protobuf
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN go install github.com/asim/go-micro/cmd/protoc-gen-micro/v4@latest

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Generate protobuf files
RUN cd /app && \
    protoc --proto_path=. \
    --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    --micro_out=. --micro_opt=paths=source_relative \
    service/producer/proto/producer.proto

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o producer-service ./service/producer/main.go

# ---- Final Stage ----
FROM alpine:latest

# Install poppler-utils for pdftotext
RUN apk add --no-cache poppler-utils

# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/producer-service .

# Expose port 8080 to the outside world
EXPOSE 8080

# Switch to our non-root user
USER appuser

# Command to run the executable
CMD ["./producer-service"] 