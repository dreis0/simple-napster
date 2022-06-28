PORT?=10098
DOCKER_PATH?=deps/docker-compose.yaml

generate-protos:
	 @protoc protos/*.proto --go_out=protos/messages --go-grpc_out=protos/services
	 @echo Protobuf files generated

run-server:
	@go run server/*.go ${PORT}

create-deps:
	@docker-compose -f ${DOCKER_PATH} up -d
	#TODO: create database bootstrap
