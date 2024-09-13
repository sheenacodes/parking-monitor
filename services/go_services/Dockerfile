# First stage: build the Go applications
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the main backend application
RUN CGO_ENABLED=0 GOOS=linux go build -o svc_backend cmd/svc_backend/main.go

# Build the generator tool
RUN CGO_ENABLED=0 GOOS=linux go build -o svc_generator cmd/svc_generator/main.go

# Second stage: create the runtime container
FROM alpine:latest


WORKDIR /app

# Copy the built binaries from the builder stage
COPY --from=builder /app/svc_backend .
COPY --from=builder /app/svc_generator .

