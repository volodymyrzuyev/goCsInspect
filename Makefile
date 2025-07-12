build_test:
	go build -ldflags "-X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=ignore" -o bin/test/goCsInspect cmd/test/main.go

build_dataFetcher:
	go build -ldflags "-X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=ignore" -o bin/dataFetcher cmd/getProtoData/main.go

run_dataFetcher: build_dataFetcher
	bin/dataFetcher

run:
	GOLANG_PROTOBUF_REGISTRATION_CONFLICT=ignore go run cmd/test/main.go



test:
	GOLANG_PROTOBUF_REGISTRATION_CONFLICT=ignore go test ./... -v
