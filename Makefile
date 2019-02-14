build:
	go build -o ./bin/ambassador-auth-jwt

build_static:
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/ambassador-auth-jwt -a -tags netgo -ldflags '-w'

docker: build_static
	docker build . -t "tomwganem/ambassador-auth-jwt:latest"

test:
	go test ./pkg/...
