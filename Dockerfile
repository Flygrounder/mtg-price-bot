ARG VERSION

FROM golang:1.20.4-alpine3.18
ARG VERSION
ENV VERSION=$VERSION
COPY . /go/src/go-mtg-vk
WORKDIR /go/src/go-mtg-vk
RUN go build ./cmd/$VERSION

FROM alpine:3.18
RUN mkdir /app
WORKDIR /app
ARG VERSION
ENV VERSION=$VERSION
COPY --from=0 /go/src/go-mtg-vk/$VERSION .
ENV GIN_MODE=release
ENTRYPOINT ./$VERSION