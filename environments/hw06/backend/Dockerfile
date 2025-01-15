# Stage 1: Build the Go app
FROM golang:1.23 AS build
WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ /app
# RUN go build -o socnet ./cmd/socnet/.
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o socnet ./cmd/socnet/.

# Stage 2: Run the app
FROM alpine:latest AS backend
WORKDIR /app
RUN apk add --no-cache libc6-compat
COPY --from=build /app/socnet ./socnet
COPY environments/hw06/backend/socnet.yaml ./
CMD ["./socnet"]
EXPOSE 8080
