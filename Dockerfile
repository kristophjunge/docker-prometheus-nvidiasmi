FROM nvidia/cuda:8.0-runtime-ubuntu16.04

MAINTAINER Kristoph Junge <kristoph.junge@gmail.com>

RUN apt-get update && \
    apt-get -y install golang --no-install-recommends && \
    rm -r /var/lib/apt/lists/*

WORKDIR /go

COPY . .

RUN go build -v -o bin/app src/app.go

CMD ["./bin/app"]
