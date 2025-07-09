# ---- Base Builder ----
FROM golang:1.24-alpine AS builder
ENV GO111MODULE=on
ENV GOPROXY="https://proxy.golang.org,direct"
ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org,direct
ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org,direct

# Set the Current Working Directory inside the container
WORKDIR /app

RUN apk add --no-cache git ca-certificates && update-ca-certificates
RUN apk add --no-cache git ca-certificates && update-ca-certificates

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api-service ./service/api/main.go

# ---- Final Stage ----
FROM alpine:latest

# Install poppler-utils for pdftotext
RUN apk add --no-cache poppler-utils

# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/api-service .

# Expose port 3000 to the outside world
EXPOSE 3000

# Switch to our non-root user
USER appuser

# Command to run the executable
CMD ["./api-service"]