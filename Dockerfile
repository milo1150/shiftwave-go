# Use official Go image
FROM golang:1.23 as builder

# Set the working directory
WORKDIR /app

# Install Air for live reloading
RUN go install github.com/air-verse/air@latest

# Copy Go modules and dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the rest of the application code
COPY . .

# Expose the application port
EXPOSE 8080

# For running in dev mode:
CMD ["air", "-c", ".air.toml"]