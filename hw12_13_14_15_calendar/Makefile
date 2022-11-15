BIN_CND := "./bin/calendar"
BIN_CND_SENDER := "./bin/calendar_sender"
BIN_CND_SCHEDULER := "./bin/calendar_scheduler"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN_CND) -ldflags "$(LDFLAGS)" ./cmd/calendar  
	go build -v -o $(BIN_CND_SENDER) -ldflags "$(LDFLAGS)" ./cmd/calendar_sender
	go build -v -o $(BIN_CND_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/calendar_scheduler

run-calendar: build
	$(BIN_CND) -config ./configs/config.toml

run-sender: build
	$(BIN_CND_SENDER) -config ./configs/sender_config.toml

run-scheduler: build
	$(BIN_CND_SCHEDULER) -config ./configs/scheduler_config.toml
	
up:
	docker-compose -f ./deployments/docker-compose.yaml up

down:
	docker-compose -f ./deployments/docker-compose.yaml down
	
version: build
	$(BIN_CND) version
	$(BIN_CND_SENDER) version
	$(BIN_CND_SCHEDULER) version
test:
	go test -race ./...

generate:
	protoc ./internal/server/grpc/api/calendar.proto --go_out=./internal/server/grpc/pb --go-grpc_out=./internal/server/grpc/pb  

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint
