FROM golang:1.20-alpine AS builder

WORKDIR /build
COPY . .

# Download dependencies
RUN go mod init pike13sync
RUN go mod tidy
RUN go get golang.org/x/oauth2/google
RUN go get google.golang.org/api/calendar/v3
RUN go get google.golang.org/api/option

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o pike13sync .

# Create a minimal image
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

# Set environment variable to indicate Docker environment
ENV DOCKER_ENV=true

# Create app directory structure
WORKDIR /app
RUN mkdir -p /app/config /app/credentials /app/logs

# Copy the binary from builder stage
COPY --from=builder /build/pike13sync /app/

# Set permissions
RUN chmod +x /app/pike13sync

# Use non-root user for security
RUN adduser -D appuser
RUN chown -R appuser:appuser /app
USER appuser

# Set the entry point
ENTRYPOINT ["/app/pike13sync"]
