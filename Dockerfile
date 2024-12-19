# Use official Go image
FROM golang:1.23 as builder

# Set the working directory
WORKDIR /app

# Install Air for live reloading
RUN go install github.com/air-verse/air@latest

# Install Staticcheck
RUN go install honnef.co/go/tools/cmd/staticcheck@latest

# Copy Go modules and dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod tidy && go mod verify

# Copy the rest of the application code
COPY . .

# Run Staticcheck to check the code
RUN staticcheck ./...

# Expose the application port
EXPOSE 8080

# For running in dev mode:
CMD ["air", "-c", ".air.toml"]