# ---- Base Builder with Dependencies Cache ----
FROM golang:1.24-alpine AS deps
WORKDIR /app
ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org,direct
ENV GOCACHE=/root/.cache/go-build

# Install build dependencies
RUN apk add --no-cache git protobuf protobuf-dev ca-certificates && update-ca-certificates

# Install protobuf tools (these are expensive operations - cache them)
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
    go install github.com/asim/go-micro/cmd/protoc-gen-micro/v4@latest

# ---- Builder Stage ----
FROM deps AS builder
WORKDIR /app

# Copy entire repo for context
COPY . .

# Generate protobuf files
RUN protoc --proto_path=./service/producer/proto \
    --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    --micro_out=. --micro_opt=paths=source_relative service/producer/proto/producer.proto

# Change to producer directory and build without workspace (avoid module conflicts)
WORKDIR /app/service/producer
ENV GOWORK=off
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o producer-service .

# ---- Runtime Stage ----
FROM alpine:latest AS runtime
RUN apk add --no-cache poppler-utils ca-certificates
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /app
COPY --from=builder /app/service/producer/producer-service .
USER appuser
EXPOSE 8080
CMD ["./producer-service"]
