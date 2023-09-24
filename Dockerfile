FROM golang:1.21

WORKDIR /app

COPY . .

RUN go build -o docker-gocache main.go

ENV PORT=8001

ENV isAPI=false

CMD [ "./docker-gocache" ]