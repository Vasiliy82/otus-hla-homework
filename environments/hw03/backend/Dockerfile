# Stage 1: Build the Go app
FROM golang:1.23 AS build
WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ /app
# RUN go build -o main ./app/.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./app/.

# Stage 2: Run the app
FROM alpine:latest AS backend
WORKDIR /app
COPY --from=build /app/main ./main
CMD ["./main"]
EXPOSE 8080
