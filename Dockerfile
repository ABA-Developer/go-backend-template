FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/go-backend-nba ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/go-backend-nba ./

COPY --from=builder /app/.env .env

ENTRYPOINT [ "./go-backend-nba" ]