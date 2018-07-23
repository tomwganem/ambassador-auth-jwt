FROM golang:1.10.2-alpine3.7 as builder
MAINTAINER Tom Ganem <tganem@us.ibm.com>

RUN apk add --no-cache --virtual .build-deps \
        curl \
        git \
        openssh \
        make \
    && curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go/src/github.com/tomwganem/ambassador-auth-jwt
ADD . ./
RUN dep ensure -v
RUN go build -v

FROM alpine:3.7
COPY --from=builder /go/bin/* /usr/local/bin/
CMD ["auth-jwt-http"]