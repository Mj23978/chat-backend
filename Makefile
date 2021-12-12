pb:
	protoc -I=./protos  protos/*.proto --go_out=./pkg/proto --go-grpc_out=./pkg/proto

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o server ./cmd/server

run-server:
	go run cmd/server/server.go -c ./conf/server.toml

devspace:
	devspace dev -n demo

migration:
	go run tests/cql/run_migrations.go