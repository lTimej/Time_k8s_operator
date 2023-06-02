SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: build
build:
	go build -o bin/manager ./cmd/server/main.go && ./bin/manager

.PHONY: run
run:
	go run ./cmd/server/main.go

.PHONY: docker-build
docker-build:
	docker build -t ${IMG} .

.PHONY: proto
proto: 
	protoc --proto_path=./pb --go_out=./pb  --go-grpc_out=./pb ./pb/proto/*.proto

.PHONY: mod
mod:
	go mod tidy