FROM golang:1.23 AS devbackend
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
# COPY . .
RUN go install github.com/air-verse/air@latest
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
CMD ["air"]