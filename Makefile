PORT?=10098
DOCKER_PATH?=deps/docker-compose.yaml

.PHONY: protos
protos:
	 @protoc protos/*.proto --go_out=. --go-grpc_out=.
	 @echo Protobuf files generated

run-server:
	@go run server/*.go ${PORT}

run-peer:
	@go run peer/*.go ${PORT}

create-deps:
	@docker-compose -f ${DOCKER_PATH} up -d
	@sleep 5
	@go run deps/*.go
