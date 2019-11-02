FROM golang:latest

COPY . /go/src/go-mtg-vk
WORKDIR /go/src/go-mtg-vk
RUN mkdir logs
ENTRYPOINT ./run.sh