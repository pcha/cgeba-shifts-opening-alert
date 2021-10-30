FROM golang:1.16-alpine as builder

RUN echo "Build...."

ENV GO111MODULE=on

WORKDIR /go/src/app/
RUN apk add git
RUN apk add build-base musl-dev

COPY ./ ./
RUN go mod download -x

RUN go get -d -v ./cmd/bot
RUN go build ./cmd/bot

ENTRYPOINT ./bot