# Stage 1: Build the Go app
FROM golang:1.23 AS build

WORKDIR /app

RUN echo "Обновление пакетов..." && \
    apt-get update && apt-get install -y \
    build-essential \
    libssl-dev \
    pkg-config

COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ /app

# Build binaries for each service
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /build/posts-ws ./cmd/posts-ws/.
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /build/dialogs ./cmd/dialogs/.
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /build/socnet ./cmd/socnet/.

# Stage 2: Run the app
FROM alpine:latest AS backend
WORKDIR /app
RUN apk add --no-cache libc6-compat
COPY --from=build /build/* ./
COPY environments/hw06/backend/socnet.yaml ./
EXPOSE 8080
