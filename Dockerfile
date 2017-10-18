FROM scratch
MAINTAINER Kevin Minehart <kmineh0151@gmail.com>
ADD ./bin/auth-jwt-http /auth-jwt-http
ENTRYPOINT ["/auth-jwt-http"]
