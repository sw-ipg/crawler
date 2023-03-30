FROM golang:1.19-alpine

EXPOSE 80

WORKDIR /go/src/crawler
COPY . .

RUN go install -mod vendor

ENTRYPOINT crawler