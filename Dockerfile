# Dockerfile

# Step 1: Use the official Go image to build the application
FROM golang:1.21-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your source code
COPY . .

# Build the Go application.
# CGO_ENABLED=0 is important for creating a static binary.
RUN CGO_ENABLED=0 GOOS=linux go build -o /ascend-api .

# Step 2: Create the final, lightweight image
FROM alpine:latest

# Set the working directory
WORKDIR /

# Copy the built binary from the 'builder' stage
COPY --from=builder /ascend-api /ascend-api

# This port will be exposed by the container
EXPOSE 8080

# Command to run the application
ENTRYPOINT ["/ascend-api"]