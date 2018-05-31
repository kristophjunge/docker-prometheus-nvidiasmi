FROM nvidia/cuda

MAINTAINER Kristoph Junge <kristoph.junge@gmail.com>

RUN apt-get update && \
    apt-get -y install golang --no-install-recommends && \
    rm -r /var/lib/apt/lists/*

WORKDIR /go

COPY . .

RUN go build -v -o bin/app src/app.go

EXPOSE 9202

CMD ["./bin/app"]
