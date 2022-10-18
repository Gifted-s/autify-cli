 # syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /app
ADD . /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
RUN go env -w GO111MODULE=on

RUN go build -o fetch