FROM golang:1.21-alpine as server-build

WORKDIR  /go/src/full-stack-server
COPY . .

RUN apk upgrade --update && \
    apk --no-cache add git