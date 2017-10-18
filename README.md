# Ambassador Auth JWT Service

## Using the service

This service is not responsible for creating and assigning JWTs.

It will take the provided secret and decode the JWT / `Bearer` token (provided by the `Authorization` header) and assign the decoded key to `x-decoded-jwt`. (or whatever you configure it to)

## Configuration

Provide the following environment variables:

| name | description | default value |
|------|-------------|---------------|
| `JWT_SECRET` | The secret used to encode / decode the JWT | |
| `JWT_COOKIE_NAME` | The name of the cookie to check for the token | `jwt` |
| `JWT_OUTBOUND_HEADER` | The name of the header to put the decoded payload in | `x-decoded-jwt` |
| `CHECK_EXP` | Set to true if you want this service to look at the RFC3339 timestamp in the `exp` value of the payload to determine if the token is expired | `false` |

## Run on Kubernetes

Take a look at the [`kubernetes.yaml`](kubernetes.yaml) file provided and modify accordingly, and then run 

```
kubectl apply -f kubernetes.yaml
```

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
