SERVER_PORT?=10098
SERVER_IP?=localhost
PEER_PORT?=3000
FILES_PATH?=./peer_1

DOCKER_PATH?=deps/docker-compose.yaml

.PHONY: protos
protos:
	 @protoc protos/*.proto --go_out=. --go-grpc_out=.
	 @echo Protobuf files generated

run-server:
	@go run server/*.go --port ${SERVER_PORT}

run-peer:
	@go run peer/*.go --port ${PEER_PORT} --files-path ${FILES_PATH} --server-ip ${SERVER_IP} --server-port ${SERVER_PORT}

create-deps:
	@docker-compose -f ${DOCKER_PATH} up -d
	@sleep 5
	@go run deps/*.go --env .env
