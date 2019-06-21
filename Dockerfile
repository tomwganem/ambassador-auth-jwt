ARG VERSION=1.12.6-alpine3.10
FROM golang:${VERSION} as builder
LABEL maintainer="Tom Ganem <tganem@us.ibm.com>"

RUN apk add --no-cache --virtual .build-deps \
        curl \
        git \
        openssh \
        make \
        g++ \
    && curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

RUN go get -u golang.org/x/lint/golint

ENV GO111MODULE on
WORKDIR /go/src/github.com/tomwganem/ambassador-auth-jwt
ADD . ./
RUN go get ./...
RUN go build -v -o /go/bin/ambassador-auth-jwt

FROM alpine:3.10
COPY --from=builder /go/bin/* /usr/local/bin/
RUN adduser -D -u 1000 app
USER app
CMD ["ambassador-auth-jwt"]
