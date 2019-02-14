# Change Log

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

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
