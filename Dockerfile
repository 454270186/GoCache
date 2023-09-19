FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./
COPY gocache server main.go ./

RUN go build -o /docker-gocache