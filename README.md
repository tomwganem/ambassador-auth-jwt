# Ambassador Auth JWT Service

## Using the service

This service is not responsible for creating and assigning JWTs.

It will take the provided secret and decode the JWT (provided by the `Authorization` header) and assign the decoded key to `x-decoded-jwt`.

## Configuration

Provide the following environment variables:

| name | description | default value |
|------|-------------|---------------|
| `JWT_SECRET` | The secret used to encode / decode the JWT | |

## Run on Kubernetes
