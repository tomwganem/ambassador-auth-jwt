build:
	go build -o ./bin/auth-jwt-http ./cmd/auth-jwt-http

build_static:
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/auth-jwt-http -a -tags netgo -ldflags '-w' ./cmd/auth-jwt-http

docker: build_static
	docker build . -t "kminehart/ambassador-auth-jwt:latest"
	docker build . -t "kminehart/ambassador-auth-jwt:v1.2.1"
test:
	go test ./pkg/...
