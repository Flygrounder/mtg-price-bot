FROM golang:1.20.4-alpine3.18

ARG VERSION
ENV VERSION=$VERSION

COPY . /go/src/go-mtg-vk
WORKDIR /go/src/go-mtg-vk
RUN go build ./cmd/$VERSION
RUN mkdir logs
ENV GIN_MODE=release
ENTRYPOINT ./$VERSION
