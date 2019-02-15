FROM golang:1.11 as build-env

COPY ./ /go/src/github.com/kristaxox/unilog
WORKDIR /go/src/github.com/kristaxox/unilog

RUN go get -v ./...
RUN go build -o /opt/collector.bin ./cmd

######

FROM alpine

COPY --from=build-env /opt/collector.bin /opt/collector.bin

ENTRYPOINT /opt/collector.bin