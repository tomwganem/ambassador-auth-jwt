ARG VERSION=1.19.2
FROM golang:${VERSION} as builder
LABEL maintainer="Tom Ganem <tganem@us.ibm.com>"

RUN go install golang.org/x/lint/golint@latest

ENV GO111MODULE on
WORKDIR /go/src/github.com/tomwganem/ambassador-auth-jwt
ADD . ./
RUN go get ./...
RUN make build_static

FROM scratch
COPY --from=builder /go/src/github.com/tomwganem/ambassador-auth-jwt/bin/* /
USER 10101
CMD ["/ambassador-auth-jwt"]
