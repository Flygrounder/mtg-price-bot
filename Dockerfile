FROM golang:1.24.2-alpine3.21 AS builder
WORKDIR /src/app
COPY go.mod .
COPY src/ ./src/
RUN go build -o mtg-price-bot ./src/main.go

FROM alpine:3.21
WORKDIR /root/
COPY --from=builder /src/app/mtg-price-bot ./mtg-price-bot
ENTRYPOINT ["./mtg-price-bot"]
