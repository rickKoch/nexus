FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git for go mod download
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /segment-service ./internal/segments

FROM alpine:3.19

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /segment-service .

EXPOSE 8080

CMD ["./segment-service"]
