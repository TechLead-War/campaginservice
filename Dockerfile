FROM golang:1.23 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o main .
# Install migrate from source (correct arch)
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.16.2

FROM debian:bookworm-slim
WORKDIR /app

RUN apt-get update && apt-get install -y postgresql-client && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/main .
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY migrations /app/migrations

COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["./main"]