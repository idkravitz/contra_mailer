FROM golang:latest

MAINTAINER Dmitry Kravtsov <idkravitz@gmail.com>

RUN useradd -r -m contra

USER contra
WORKDIR /home/contra
ENV GOPATH /home/contra

RUN mkdir -p bin pkg templates src/github.com/kravitz/contra_mailer
ADD . src/github.com/kravitz/contra_mailer/
ADD templates/ templates/

RUN go install github.com/kravitz/contra_mailer

ENTRYPOINT ["./bin/contra_mailer"]
