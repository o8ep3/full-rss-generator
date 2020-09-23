FROM golang:1.15.2-alpine3.12

COPY /data/api /go/src/api/

ENV GO111MODULE=on

WORKDIR /go/src/api/

RUN apk update \
  && apk add --no-cache git \
  && apk --no-cache add postgresql-client \
  && go get github.com/oxequa/realize \
  && go mod download \
  && go build -o fullRssGeneratorAPI

EXPOSE 8080
