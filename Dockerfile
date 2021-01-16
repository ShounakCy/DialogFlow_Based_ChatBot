FROM golang:1.12.7 as build-env

MAINTAINER parath.singh@havells.com

RUN mkdir /app

WORKDIR /app

COPY go.mod /app


RUN GO111MODULE=on go mod download

COPY . /app

RUN GO111MODULE=on CGO_ENABLED=0 go build -o /bin/nlp


FROM alpine:3.8
#Add ca certificates
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN update-ca-certificates

COPY --from=build-env /bin/nlp /nlp
RUN apk add --no-cache tini bash

ENTRYPOINT ["/sbin/tini", "--"]
