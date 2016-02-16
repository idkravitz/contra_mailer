FROM golang:latest

MAINTAINER Dmitry Kravtsov <idkravitz@gmail.com>

RUN useradd -r -m contra

USER contra
WORKDIR /home/contra
ENV GOPATH /home/contra

RUN mkdir -p bin pkg templates src/github.com/kravitz/contra_mailer
RUN go get -u	golang.org/x/net/context && \
	  go get -u golang.org/x/oauth2 && \
	  go get -u golang.org/x/oauth2/google && \
	  go get -u google.golang.org/api/gmail/v1

ADD . src/github.com/kravitz/contra_mailer/

# ADD config.json config.json
# ADD credentials.json credentials.json
# ADD client_secret.json client_secret.json

ADD templates/ templates/

RUN go install github.com/kravitz/contra_mailer

ENTRYPOINT ["./bin/contra_mailer"]
