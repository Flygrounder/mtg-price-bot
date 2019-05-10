FROM golang:1.24.2-alpine3.21 AS builder
WORKDIR /src/app
COPY main.go go.mod ./
RUN go build -o mtg-price-bot main.go

FROM alpine:3.21
WORKDIR /root/
COPY --from=builder /src/app/mtg-price-bot ./mtg-price-bot
ENTRYPOINT ["./mtg-price-bot"]
