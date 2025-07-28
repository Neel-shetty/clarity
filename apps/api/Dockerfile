# Build stage
FROM golang:1.24.5-alpine3.22 AS builder
WORKDIR /app

# Cache modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# Final stage
FROM alpine:latest
WORKDIR /root/

# Copy built binary
COPY --from=builder /app/server .

EXPOSE 8080
CMD ["./server"]
