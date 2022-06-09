# docker build . -t ghcr.io/reddio-com/red-adapter:latest
FROM golang:1.16-buster as builder

RUN mkdir /build
COPY . /build
RUN cd /build && make adapter

FROM ubuntu:22.04

RUN apt-get update && apt-get install -y ca-certificates && apt-get clean

RUN mkdir /config
COPY ./config/adapter.toml /config/adapter.toml
COPY --from=builder /build/bin/adapter /adapter

CMD ["/adapter"]
