# Ambassador Auth JWT-RSA Service

This is a fork of [kminehart/ambassador-auth-jwt](https://github.com/kminehart/ambassador-auth-jwt), which is able to verify HMAC based tokens, but not RSA ones. This module is only meant to verify RSA JWTs.

## Using the service

This service is not responsible for creating and assigning JWTs.

It decodes JWT / `Bearer` tokens (provided by the `Authorization` header) and verify the token against a JWKSet, provided by the `JWT_ISSUER` env variable.

It will return a 200 if it can verify the token, 401 if not.

## Configuration

Provide the following environment variables:

| name | description | default value |
|------|-------------|---------------|
| `JWT_ISSUER` | public endpoint with JWKSet (A set of public key) to verify tokens against | |
| `JWT_OUTBOUND_HEADER` | The name of the header to put the decoded payload in | `X-JWT-PAYLOAD` |
| `CHECK_EXP` | check if the token is expired or not | `true` |
| `ALLOW_BASIC_AUTH_PASSTHROUGH` | allow basic auth requests, without a token, to pass through  | `false` |
| `ALLOW_BASIC_AUTH_HEADER` | specify the header that has the basic auth credentials  | `Authorization` |
| `ALLOW_BASIC_AUTH_PATH_REGEX` | specify a regex to test the path of the request determine if a basic auth request should be allowed | `^/.*` |

## Run on Kubernetes

A helm chart is included as a git submodule in the helm directory. You can check out the chart at https://github.com/tomwganem/ambassador-auth-jwt-helm

## Logging

All logs are in json format to be consumed in a ELK stack.
