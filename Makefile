PROJECT=red-adapter
GO              := GO111MODULE=on CGO_ENABLED=1 go
GOBUILD         := $(GO) build $(BUILD_FLAG) -tags codes
PACKAGE_LIST  := go list ./...
PACKAGE_DIRECTORIES := $(PACKAGE_LIST) | sed 's|github.com/reddio-com/$(PROJECT)/||'

all: check adapter

adapter:
	$(GOBUILD) -o bin/adapter cmd/adapter/adapter.go 
init_testing_db:
	cd misc && docker-compose up -d
	sleep 30
clean_testing_db:
	cd misc && docker-compose down && rm -rf db_data
local_adapter: adapter clean_testing_db init_testing_db
	DSN='root:password@tcp(127.0.0.1:3306)/reddifi?charset=utf8mb4&parseTime=True&loc=Local' ./bin/adapter
tools-dir:
	mkdir -p tools/bin
golangci-lint:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b ./tools/bin v1.44.2

check: golangci-lint
	GO111MODULE=on CGO_ENABLED=0 tools/bin/golangci-lint run -v $$($(PACKAGE_DIRECTORIES)) --config .golangci.yml