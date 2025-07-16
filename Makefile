buildName = goCsInspect
goCsInsRestDir = goCsInspectAPI/
dataFetcherDir = protoFetcher/


build_fetcher:
	go build -ldflags "-X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=ignore" -o bin/$(dataFetcherDir)/$(buildName) cmd/$(dataFetcherDir)/main.go

fetch_data: build_fetcher
	bin/$(dataFetcherDir)/$(buildName) --skip

gen_sql:
	cd pkg/storage/sqlite/sql && sqlc generate

build_all: build_fetcher build_api

build_api: gen_sql
	go build -ldflags "-X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=ignore" -o bin/$(goCsInsRestDir)/$(buildName) cmd/$(goCsInsRestDir)/main.go

run_api: build_api
	bin/$(goCsInsRestDir)/$(buildName)

test: gen_sql fetch_data
	GOLANG_PROTOBUF_REGISTRATION_CONFLICT=ignore go test ./... -v
