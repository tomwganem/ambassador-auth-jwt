FROM golang:1.11.5-alpine3.9 as builder
LABEL maintainer="Tom Ganem <tganem@us.ibm.com>"

RUN apk add --no-cache --virtual .build-deps \
        curl \
        git \
        openssh \
        make \
    && curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go/src/github.com/tomwganem/ambassador-auth-jwt
ADD . ./
RUN dep ensure -v
RUN go build -v -o /go/bin/ambassador-auth-jwt

FROM alpine:3.9
COPY --from=builder /go/bin/* /usr/local/bin/
CMD ["ambassador-auth-jwt"]
