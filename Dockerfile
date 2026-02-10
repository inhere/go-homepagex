# Build stage
FROM golang:1.21-alpine AS backend-builder

WORKDIR /app
COPY backend/ ./backend/
RUN cd backend && go mod tidy && go build -o home-dashboard

# Frontend is static, no build needed

# Final stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy backend binary
COPY --from=backend-builder /app/backend/home-dashboard ./

# Copy frontend static files
COPY frontend/build/ ./frontend/build/

# Copy config and pages
COPY config/ ./config/
COPY pages/ ./pages/

# Expose port
EXPOSE 8080

# Run
CMD ["./home-dashboard", "./config/config.yaml"]
