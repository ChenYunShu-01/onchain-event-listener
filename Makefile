GO              := GO111MODULE=on CGO_ENABLED=1 go
GOBUILD         := $(GO) build $(BUILD_FLAG) -tags codes

adapter:
	$(GOBUILD) -o bin/adapter cmd/adapter/adapter.go 
init_testing_db:
	cd misc && docker-compose up -d
	sleep 30
clean_testing_db:
	cd misc && docker-compose down && rm -rf db_data
local_adapter: adapter clean_testing_db init_testing_db
	DSN='root:password@tcp(127.0.0.1:3306)/reddifi?charset=utf8mb4&parseTime=True&loc=Local' ./bin/adapter