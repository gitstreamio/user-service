#!/bin/sh

go generate ./...

CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"'  -o user-service .

docker build -t user-service .
