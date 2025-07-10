build_test:
	go build -ldflags "-X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=ignore" -o bin/test/goCsInspect cmd/test/main.go

run:
	GOLANG_PROTOBUF_REGISTRATION_CONFLICT=ignore go run cmd/test/main.go

test:
	GOLANG_PROTOBUF_REGISTRATION_CONFLICT=ignore go test ./... -v
