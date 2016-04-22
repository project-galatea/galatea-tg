FROM golang:1.6

ARG BOT_TOKEN
ENV BOT_TOKEN ${BOT_TOKEN}


ADD . /go/src/github.com/project-galatea/galatea-tg

WORKDIR /go/src/github.com/project-galatea/galatea-tg


RUN apt-get update

RUN apt-get -y install protobuf-compiler

RUN go get -u github.com/golang/protobuf/proto && go get -u github.com/golang/protobuf/protoc-gen-go

RUN BOT_TOKEN=$BOT_TOKEN make install

ENTRYPOINT ["/go/bin/galatea-tg"]

EXPOSE 8443
