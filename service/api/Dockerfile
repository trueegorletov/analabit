FROM golang:1.24-alpine

WORKDIR /app

# Install git for fetching modules
RUN apk add --no-cache git

# Copy source code
COPY . .

# Ensure modules are tidy and download deps
RUN go mod tidy

# Build the API binary
RUN go build -o api-service ./service/api/main.go

# Expose the API port (matches config default)
EXPOSE 3000

CMD ["./api-service"] 