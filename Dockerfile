FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum

RUN go mod download

COPY . .

RUN go mod tidy

RUN go build -o go-backend-nba

ENTRYPOINT [ "app/go-backend-nba" ]