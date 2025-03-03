.PHONY: build

BINARY_NAME=goCsInspect
BUILD_DIR=./build
PROJECT_ROOT_DIR=$(PWD)
SQLC_DIR=$(PROJECT_ROOT_DIR)/cmd/storage/db/sqlc/sql

build:
	cd $(SQLC_DIR) && sqlc generate
	cd $(PROJECT_ROOT_DIR) && go clean -i
	cd $(PROJECT_ROOT_DIR) && go build -o $(BUILD_DIR)/$(BINARY_NAME) -ldflags "-X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=ignore" ./main.go
	> log.log


run:
	$(MAKE) build
	./build/goCsInspect

