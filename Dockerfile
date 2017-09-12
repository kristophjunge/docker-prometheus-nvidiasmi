FROM golang:alpine

MAINTAINER Kritoph Junge <kristoph.junge@gmail.com>

WORKDIR /go/src/app

COPY . .

RUN go build -v -o bin/app app.go

CMD ["./bin/app"]