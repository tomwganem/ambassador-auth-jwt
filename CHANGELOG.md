# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

### [0.4.2-rc.1](https://github.com/tomwganem/ambassador-auth-jwt/compare/v0.4.1...v0.4.2-rc.1) (2021-03-05)


### Features

* **OPSENG-985:** supporting multiple jwt issuers ([6a38a67](https://github.com/tomwganem/ambassador-auth-jwt/commit/6a38a679688c2168e73534c04633cb7e2ee65202))
* distroless ([af7f97f](https://github.com/tomwganem/ambassador-auth-jwt/commit/af7f97fea8259194f28cb60f8984cde4ea9e2737))

### [0.4.2-rc.0](https://github.com/tomwganem/ambassador-auth-jwt/compare/v0.4.1...v0.4.2-rc.0) (2020-05-13)


### Features

* **OPSENG-985:** supporting multiple jwt issuers ([6a38a67](https://github.com/tomwganem/ambassador-auth-jwt/commit/6a38a679688c2168e73534c04633cb7e2ee65202))

### [0.4.1](https://github.com/tomwganem/ambassador-auth-jwt/compare/v0.4.0...v0.4.1) (2020-05-01)

<a name="0.4.0"></a>
# [0.4.0](https://github.com/tomwganem/ambassador-auth-jwt/compare/v0.3.0...v0.4.0) (2019-07-23)


### Bug Fixes

* allow basic auth credentials in the Authorization header to pass ([222c287](https://github.com/tomwganem/ambassador-auth-jwt/commit/222c287))
* allow setting log level from environment variable ([c234b3f](https://github.com/tomwganem/ambassador-auth-jwt/commit/c234b3f))
* delete "sub" key from claims ([e3867ba](https://github.com/tomwganem/ambassador-auth-jwt/commit/e3867ba))
* fix basic auth regex check ([7f91aed](https://github.com/tomwganem/ambassador-auth-jwt/commit/7f91aed))
* go mod tidy ([e28e7cf](https://github.com/tomwganem/ambassador-auth-jwt/commit/e28e7cf))
* provide backward compatibility in the error msg structure returned ([0831fa7](https://github.com/tomwganem/ambassador-auth-jwt/commit/0831fa7))
* set CHECK_EXP to true by default ([328763f](https://github.com/tomwganem/ambassador-auth-jwt/commit/328763f))
* update submodule ([7d6015c](https://github.com/tomwganem/ambassador-auth-jwt/commit/7d6015c))


### Features

* allow Basic Auth Requests to be passed through ([f6c5f36](https://github.com/tomwganem/ambassador-auth-jwt/commit/f6c5f36))
* allows passing basic auth creds in multiple headers ([93c6755](https://github.com/tomwganem/ambassador-auth-jwt/commit/93c6755))
* provide more structured error message ([bcab0f1](https://github.com/tomwganem/ambassador-auth-jwt/commit/bcab0f1))
* specify header to allow basic auth pass through ([c409aff](https://github.com/tomwganem/ambassador-auth-jwt/commit/c409aff))
* use golang version 1.12.6 ([8c9c9cb](https://github.com/tomwganem/ambassador-auth-jwt/commit/8c9c9cb))
* use regex to allow basic auth requests ([87bde66](https://github.com/tomwganem/ambassador-auth-jwt/commit/87bde66))



<a name="0.3.0"></a>
# [0.3.0](https://github.com/tomwganem/ambassador-auth-jwt/compare/v0.2.0...v0.3.0) (2019-03-08)


### Features

* add helm chart as sub repo ([b2ddcc2](https://github.com/tomwganem/ambassador-auth-jwt/commit/b2ddcc2))



<a name="0.2.0"></a>
# [0.2.0](https://github.com/tomwganem/ambassador-auth-jwt/compare/v0.1.0...v0.2.0) (2019-03-07)


### Bug Fixes

* fix issue where we are unable to auth tokens with 'expires_at' ([2de4949](https://github.com/tomwganem/ambassador-auth-jwt/commit/2de4949))


### Features

* add sentry ([967e21e](https://github.com/tomwganem/ambassador-auth-jwt/commit/967e21e))



<a name="0.1.0"></a>
# 0.1.0 (2019-02-14)


### Bug Fixes

* fix timecode parse ([3efb577](https://github.com/tomwganem/ambassador-auth-jwt/commit/3efb577))
* improve documentation for httpserver package ([9b43d99](https://github.com/tomwganem/ambassador-auth-jwt/commit/9b43d99))
* log in nanoseconds ([2000b7b](https://github.com/tomwganem/ambassador-auth-jwt/commit/2000b7b))
* use token's KeyID to find public key to verify with ([94e6c3c](https://github.com/tomwganem/ambassador-auth-jwt/commit/94e6c3c))


### Features

* assume that "exp" field will be in unix epoch time instead of rfc3339 ([40498e8](https://github.com/tomwganem/ambassador-auth-jwt/commit/40498e8))
* decode token for url params ([477899b](https://github.com/tomwganem/ambassador-auth-jwt/commit/477899b))
* enable cors ([2444a8f](https://github.com/tomwganem/ambassador-auth-jwt/commit/2444a8f))
* move DecodeHTTPHandler to httpserver package ([c4d0ec8](https://github.com/tomwganem/ambassador-auth-jwt/commit/c4d0ec8))
* verify RSA jwtokens by getting a JWKSet from a url ([57f1441](https://github.com/tomwganem/ambassador-auth-jwt/commit/57f1441))
