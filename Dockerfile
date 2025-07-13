FROM golang:1.24-alpine

# Install stdbuf for unbuffered output
RUN apk add --no-cache coreutils

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o todos-app .

# Expose the application port
EXPOSE 8080

# Command to run the application with unbuffered output
CMD ["stdbuf", "-oL", "./todos-app"]
