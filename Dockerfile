# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.3/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.21
WORKDIR /app

ARG APP_ENV_FILE
COPY ${APP_ENV_FILE:-app.env} ./app.env

COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./migration

EXPOSE 8080
LABEL authors="vitaliiiavurek"
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]