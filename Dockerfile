FROM golang:1.12-alpine as builder

RUN apk add --no-cache make git;
WORKDIR /deadlock/agent/
RUN cd /deadlock/agent/ && mkdir -p build
COPY hlserver hlserver
COPY handler handler
COPY zaplog zaplog
COPY example example
COPY Makefile Makefile
COPY main.go main.go
COPY Dockerfile Dockerfile
COPY go.mod go.mod
COPY go.sum go.sum
RUN make go

FROM alpine
LABEL author="tinywell"

RUN apk add --no-cache \
    bash;
RUN mkdir -p /etc/agent/
COPY --from=builder /deadlock/agent/build/bin/deadlock /usr/local/bin/
COPY --from=builder /deadlock/agent/build/bin/benchmark /usr/local/bin/

WORKDIR /etc/agent/
EXPOSE 8000
EXPOSE 6060

CMD ["deadlock"]

