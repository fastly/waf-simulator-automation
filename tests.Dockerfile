FROM golang:1.21.6

ENV GO111MODULE=on

ADD . /app
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o waftests tests/main.go

ENTRYPOINT ["/app/waftests"]