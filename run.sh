#!/bin/bash
if [[ $MODE = "test" ]]
then
	go test ./...
elif [[ $MODE = "prod" ]]
then
	go build .
	export GIN_MODE="release"
	./go-mtg-vk
fi
