# based on http://carlosbecker.com/posts/small-go-apps-containers/
FROM alpine:3.2

ENV GOROOT=/usr/lib/go \
    GOPATH=/gopath \
    GOBIN=/gopath/bin \
    PATH=$PATH:$GOROOT/bin:$GOPATH/bin \
    HASHTOCK_SERVE_ADDR=:80 \
    HASHTOCK_DB_NAME=hashtock

WORKDIR /gopath/src/github.com/hashtock/hashtock-go
ADD . /gopath/src/github.com/hashtock/hashtock-go

RUN apk add -U git go && \
    go get github.com/tools/godep && \
    $GOBIN/godep go build -o /usr/bin/hashtock-go && \
    apk del git go && \
    rm -rf /gopath && \
    rm -rf /var/cache/apk/*

EXPOSE 80

CMD "hashtock-go"
