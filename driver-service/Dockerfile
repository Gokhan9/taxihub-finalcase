
# Use the official Golang image to create a build artifact.
# This is known as a multi-stage build.
FROM golang:1.19-alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
# COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
# RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o driver-service ./cmd/driver-service

# ---

# Use a smaller, secure image to run the app
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/driver-service .

# Expose port 8081 to the outside world
EXPOSE 8081

# Command to run the executable
CMD ["./driver-service"]
