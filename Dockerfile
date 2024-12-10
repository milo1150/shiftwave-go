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


# # Build stage
# FROM golang:1.23 as builder

# WORKDIR /app
# COPY go.mod go.sum ./
# RUN go mod download && go mod verify
# COPY . .

# # Build the binary
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# # Runtime stage
# FROM debian:bullseye-slim
# WORKDIR /app

# # Copy the binary from the builder stage
# COPY --from=builder /app/main .

# # Ensure the binary is executable
# RUN chmod +x ./main

# # Expose port if necessary
# EXPOSE 8080

# # Run the binary
# CMD ["./main"]
