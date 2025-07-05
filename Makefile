build_test:
	go build -ldflags "-X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=ignore" -o bin/test/goCsInspect cmd/test/main.go

run_test:
	GOLANG_PROTOBUF_REGISTRATION_CONFLICT=ignore go run cmd/test/main.go
