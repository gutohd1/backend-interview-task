build_proto:
	protoc --go_out=./app/explore_service_protos --go_opt=paths=source_relative --go-grpc_out=./explore_service_protos --go-grpc_opt=paths=source_relative ./app/explore-service.proto
