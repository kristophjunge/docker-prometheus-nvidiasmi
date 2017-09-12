FROM nvidia/cuda

MAINTAINER Kritoph Junge <kristoph.junge@gmail.com>

RUN apt-get update && \
    apt-get -y install golang --no-install-recommends && \
    rm -r /var/lib/apt/lists/*

WORKDIR /go/src/app

COPY . .

RUN go build -v -o bin/app app.go

CMD ["./bin/app"]