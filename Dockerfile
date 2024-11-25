FROM golang:1.23.1-alpine AS builder

WORKDIR /app
COPY . .

RUN go build -o main cmd/main.go
RUN apk add curl
# Run stage
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar xvz
RUN mv migrate migrate.linux-amd64

RUN ls -lh /app/

FROM alpine:3.13

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate.linux-amd64 ./migrate

COPY .env .
COPY start.sh .
COPY wait-for.sh .
COPY cmd/migrate/migrations ./migrations

RUN chmod +x ./start.sh ./wait-for.sh ./migrate

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]
