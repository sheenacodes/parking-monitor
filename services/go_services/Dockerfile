# First stage: build the Go applications
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy Go modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project into the container
COPY . .

# Build the main application
RUN CGO_ENABLED=0 GOOS=linux go build -o svc_backend cmd/svc_backend/main.go

# Build the second tool
RUN CGO_ENABLED=0 GOOS=linux go build -o svc_generator cmd/svc_generator/main.go

# Second stage: create the runtime container
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the built binaries from the builder stage
COPY --from=builder /app/svc_backend .
COPY --from=builder /app/svc_generator .

# # Expose any necessary ports (if your apps need to expose ports)
# # EXPOSE 8080

# # Define entry points or commands
# CMD ["./app"]                # Start the main app
# # If you need to run another tool in the background or after the main app starts, you can chain the commands:
# # CMD ["sh", "-c", "./app & ./another-tool"]
