FROM golang:alpine

MAINTAINER Kritoph Junge <kristoph.junge@gmail.com>

LABEL com.nvidia.volumes.needed="nvidia_driver"

ENV PATH /usr/local/nvidia/bin:/usr/local/cuda/bin:${PATH}
ENV LD_LIBRARY_PATH /usr/local/nvidia/lib:/usr/local/nvidia/lib64

WORKDIR /go/src/app

COPY . .

RUN go build -v -o bin/app app.go

CMD ["./bin/app"]