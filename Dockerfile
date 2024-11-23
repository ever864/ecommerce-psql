FROM golang:1.23.1-alpine AS builder

WORKDIR /app
COPY . .

RUN go build -o main cmd/main.go

# Run stage

FROM alpine:3.13

WORKDIR /app

COPY --from=builder /app/main .
COPY .env .

EXPOSE 8080

CMD ["/app/main"]
