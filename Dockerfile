ARG VERSION=1.15.8
FROM golang:${VERSION} as builder
LABEL maintainer="Tom Ganem <tganem@us.ibm.com>"

RUN go get -u golang.org/x/lint/golint

ENV GO111MODULE on
WORKDIR /go/src/github.com/tomwganem/ambassador-auth-jwt
ADD . ./
RUN go get ./...
RUN make build_static

FROM gcr.io/distroless/static
COPY --from=builder /go/src/github.com/tomwganem/ambassador-auth-jwt/bin/* /
USER 1000
CMD ["/ambassador-auth-jwt"]
