FROM golang:1.9.2-alpine3.7

RUN mkdir -p /go/src/gitlab.com/medakk/zyxdb
ADD . /go/src/gitlab.com/medakk/zyxdb
RUN go install gitlab.com/medakk/zyxdb

ENTRYPOINT ["/go/bin/zyxdb", "-host", "0.0.0.0"]
EXPOSE 5002
