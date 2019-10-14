FROM golang:latest

COPY . ./go-mtg-vk
WORKDIR go-mtg-vk
ENTRYPOINT ./run.sh

