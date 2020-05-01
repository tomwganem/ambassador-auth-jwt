ARG VERSION=1.14.2-alpine3.11
FROM golang:${VERSION} as builder
LABEL maintainer="Tom Ganem <tganem@us.ibm.com>"

RUN apk add --no-cache --virtual .build-deps \
        curl \
        git \
        openssh \
        make \
        g++ \
        ca-certificates \
    && curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

RUN go get -u golang.org/x/lint/golint

ENV GO111MODULE on
WORKDIR /go/src/github.com/tomwganem/ambassador-auth-jwt
ADD . ./
RUN go get ./...
RUN make build_static

FROM alpine:3.11
COPY --from=builder /go/bin/* /usr/local/bin/
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN adduser -D -u 1000 app
USER app
CMD ["ambassador-auth-jwt"]
