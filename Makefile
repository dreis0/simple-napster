
generate-protos:
	 @protoc protos/*.proto --go_out=protos/messages --go-grpc_out=protos/services
	 @echo Protobuf files generated
