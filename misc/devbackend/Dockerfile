FROM golang:1.23 AS devbackend
WORKDIR /app
COPY go.mod go.sum entrypoint.sh ./
RUN go mod download
RUN apt-get update && apt-get install -y gettext-base
RUN go install github.com/air-verse/air@latest
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest
CMD ["/app/entrypoint.sh"]
# air сконфигурирован таким образом, что запустит app через dlv, поэтому здесь будет доступен не только порт приложения, но и порт отладчика
EXPOSE 2345
EXPOSE 9090
