BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

generate:
	go generate ./api

init: docker-down-clear \
	clear \
	docker-pull docker-build docker-up \
	wait-postgres init-postgres migrations-migrate

up: docker-up
down: docker-down
restart: down up

docker-pull:
	docker compose pull

docker-build:
	docker compose build --pull

docker-up:
	docker-compose up -d

docker-down:
	docker compose down --remove-orphans

docker-down-clear:
	docker compose down -v --remove-orphans

wait-postgres:
	sleep 10

clear:
	rm -rf var/*

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) -config ./configs/config.toml

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

init-postgres:
	docker exec -it calendar_postgres bash -c "\
		psql -Utest -dtest -c 'create database calendar;'; \
		psql -Utest -dtest -c 'create user calendar_user with encrypted password '\''calendar_pass'\'';'; \
		psql -Utest -dtest -c 'grant all privileges on database calendar to calendar_user;'; \
    "

# go get -u github.com/pressly/goose/cmd/goose
# brew install goose
# goose -dir migrations create init sql

migrations-migrate:
	goose -dir migrations postgres "user=calendar_user password=calendar_pass dbname=calendar host=localhost sslmode=disable" up