# syntax=docker/dockerfile:1

FROM golang:1.22.4 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build --no-cache -o /creategorestapi ./cmd/api

CMD [ "/creategorestapi" ]