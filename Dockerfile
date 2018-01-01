FROM golang:1.9.2-alpine3.7

RUN apk update && apk add git

# TODO: Don't put deps here
RUN go get github.com/gorilla/handlers
RUN go get github.com/gorilla/mux
RUN go get gopkg.in/yaml.v2

RUN mkdir -p /go/src/gitlab.com/medakk/zyxdb
ADD . /go/src/gitlab.com/medakk/zyxdb
RUN mkdir -p /etc/zyxdb/
RUN mv /go/src/gitlab.com/medakk/zyxdb/zyxdb.yml /etc/zyxdb/
RUN go install gitlab.com/medakk/zyxdb

ENTRYPOINT ["/go/bin/zyxdb"]
