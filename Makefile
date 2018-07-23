build:
	go build -o ./bin/auth-jwt-http

build_static:
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/auth-jwt-http -a -tags netgo -ldflags '-w'

docker: build_static
	docker build . -t "tomwganem/ambassador-auth-jwt:latest"

test:
	go test ./pkg/...
