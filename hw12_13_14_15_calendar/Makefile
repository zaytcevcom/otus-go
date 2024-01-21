BIN := "./bin/calendar"
BIN_SCHEDULER := "./bin/calendar_scheduler"
BIN_SENDER := "./bin/calendar_sender"

DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

generate:
	go generate ./api

init: docker-down-clear \
	clear \
	docker-network \
	docker-pull docker-build docker-up \
	wait-postgres migrations-migrate

up: docker-network docker-up wait-postgres migrations-migrate
down: docker-down
restart: down up

docker-pull:
	docker compose -f ./deployments/docker-compose.yml pull

docker-build:
	docker compose -f ./deployments/docker-compose.yml build --pull

docker-up:
	docker compose -f ./deployments/docker-compose.yml up -d

docker-down:
	docker compose -f ./deployments/docker-compose.yml down --remove-orphans

docker-down-clear:
	docker compose -f ./deployments/docker-compose.yml down -v --remove-orphans

docker-network:
	docker network create calendar_network || true

wait-postgres:
	sleep 10

clear:
	rm -rf var/postgres/* var/rabbitmq/*

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN_SCHEDULER) ./cmd/calendar_scheduler
	go build -v -o $(BIN_SENDER) ./cmd/calendar_sender

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version


test:
	go test -race ./internal/... #./pkg/...


remove-lint-deps:
	rm $(which golangci-lint)

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.54.2

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint

# go get -u github.com/pressly/goose/cmd/goose
# brew install goose
# goose -dir migrations create init sql

migrations-migrate:
	goose -dir migrations postgres "user=test password=test dbname=calendar host=localhost sslmode=disable" up

integration-tests:
	set -e ;\
	docker compose -f ./deployments/docker-compose.test.yml up --build -d ;\
	test_status_code=0 ;\
	sleep 10 ;\
	goose -dir migrations postgres "user=test password=test dbname=calendar host=localhost sslmode=disable" up ;\
	docker compose -f ./deployments/docker-compose.test.yml run integration_tests go test || test_status_code=$$? ;\
	docker compose -f ./deployments/docker-compose.test.yml down \
                --rmi local \
        		--volumes \
        		--remove-orphans \
        		--timeout 60; \
	echo $$test_status_code ;\
	exit $$test_status_code ;