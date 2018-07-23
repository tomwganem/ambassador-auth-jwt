# Ambassador Auth JWT-RSA Service

This is a fork of [kminehart/ambassador-auth-jwt](https://github.com/kminehart/ambassador-auth-jwt), which is able to verify HMAC based tokens, but not RSA ones. This module is very single purpose.

## Using the service

This service is not responsible for creating and assigning JWTs.

It decode JWT / `Bearer` tokens (provided by the `Authorization` header) and verify it against a JWKSet, provided by the `JWT_ISSUER` env variable.

It will return a 200 if it can verify the token, 401 if not.

## Configuration

Provide the following environment variables:

| name | description | default value |
|------|-------------|---------------|
| `JWT_ISSUER` | public endpoint with JWKSet to verify tokens against | |
| `JWT_OUTBOUND_HEADER` | The name of the header to put the decoded payload in | `X-JWT-PAYLOAD` |
| `CHECK_EXP` | Set to true if you want this service to look at the RFC3339 timestamp in the `exp` value of the payload to determine if the token is expired | `false` |

## Run on Kubernetes

Take a look at the [`kubernetes.yaml`](kubernetes.yaml) file provided and modify accordingly, and then run

```
kubectl apply -f kubernetes.yaml
```
**note:** it is very crucial that you set a much stronger secret before deploying this to production.

## Logging

All logs are in json format to be consumed in a ELK stack.

## Ambassador

This is intended to be used with [Ambassador](getambassador.io).

To use this with Ambassador, make sure to read the [Authentication tutorial](https://www.getambassador.io/user-guide/auth-tutorial), and then check out the example [auth.yaml](auth.yaml).

```yaml
---
apiVersion: ambassador/v0
kind: Mapping
name: my_service_mapping
prefix: /my_service/
service: my-service-http:3000
---
apiVersion: ambassador/v0
kind:  Module
name:  authentication
config:
  auth_service: "example-jwt-auth:3000"
  path_prefix: "/auth"
  allowed_headers:
  - "x-jwt-payload"
```

One thing to keep in mind:  if I did not include a `path_prefix` then Envoy would put me in a redirect loop.
