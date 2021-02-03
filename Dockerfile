FROM golang:1.15.4-alpine3.12

COPY . /go/src/go-mtg-vk
WORKDIR /go/src/go-mtg-vk/cmd/go-mtg-vk
RUN go install
WORKDIR /go/bin
RUN mkdir logs
ENV GIN_MODE=release
ENTRYPOINT ./go-mtg-vk
