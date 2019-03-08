# Change Log

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

<a name="0.3.0-rc.1"></a>
# [0.3.0-rc.1](https://github.com/tomwganem/ambassador-auth-jwt/compare/v0.2.0...v0.3.0-rc.1) (2019-03-08)


### Features

* add helm chart as sub repo ([b2ddcc2](https://github.com/tomwganem/ambassador-auth-jwt/commit/b2ddcc2))



<a name="0.3.0-rc.0"></a>
# [0.3.0-rc.0](https://github.com/tomwganem/ambassador-auth-jwt/compare/v0.2.0...v0.3.0-rc.0) (2019-03-08)


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
