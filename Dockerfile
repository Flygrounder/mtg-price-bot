FROM golang:1.20.4-alpine3.18
COPY . /go/src/go-mtg-vk
WORKDIR /go/src/go-mtg-vk
RUN go build ./cmd/bot

FROM alpine:3.18
RUN mkdir /app
WORKDIR /app
COPY --from=0 /go/src/go-mtg-vk/bot .
ENV GIN_MODE=release
ENTRYPOINT ./bot