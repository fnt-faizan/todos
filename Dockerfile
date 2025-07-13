FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o todos-app .

FROM alpine:3.18

# Install stdbuf for unbuffered output
RUN apk add --no-cache coreutils

# Set working directory
WORKDIR /app

# Copy the built application
COPY --from=builder /app/todos-app .

# Expose the application port
EXPOSE 8080

# Command to run the application with unbuffered output
CMD ["stdbuf", "-oL", "./todos-app"]
