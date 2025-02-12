# Stage 1: Build the Go app
FROM golang:1.23 AS build

WORKDIR /app

RUN echo "Обновление пакетов..." && \
    apt-get update && apt-get install -y \
    build-essential \
    libssl-dev \
    pkg-config

COPY go.mod go.sum ./
RUN go mod download
COPY . /app

# Build binaries for each service
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go mod verify
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -v -o /build/counters ./cmd/.

# Stage 2: Run the app
FROM alpine:latest AS backend
WORKDIR /app
RUN apk add --no-cache libc6-compat
COPY --from=build /build/* ./
COPY counters.yaml ./
EXPOSE 8080
