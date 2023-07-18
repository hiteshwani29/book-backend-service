# Start with a base image containing the Go runtime
FROM golang:1.19-alpine AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the source code to the working directory
COPY . .

# Build the Go application
RUN go build -o book-backend-svc

# Start with a fresh, minimal image
FROM alpine:latest

# Copy the binary from the build stage to the final image
COPY --from=build /app/book-backend-svc /usr/local/bin/book-backend-svc

# Set the entrypoint command
ENTRYPOINT ["/usr/local/bin/book-backend-svc"]
