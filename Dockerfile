# Stage 1: Build stage
FROM golang:1.23 AS build

ARG BUILD_TARGET

# Set the working directory
WORKDIR /app

# Copy and download dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o ${BUILD_TARGET} ./cmd/${BUILD_TARGET}

# Stage 2: Final stage
FROM alpine:edge

ARG BUILD_TARGET

ENV APP_NAME=${BUILD_TARGET}

# Set the working directory
WORKDIR /app

COPY data /app/data

# Copy the binary from the build stage
COPY --from=build /app/${BUILD_TARGET} .

# Set the entrypoint command
ENTRYPOINT ["/bin/sh", "-c", "/app/$APP_NAME"]
