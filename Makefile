PROJECTNAME=$(shell basename "$(PWD)")
VERSION=-ldflags="-X main.Version=$(shell git describe --tags)"


.PHONY: help get build test clean
help: Makefile
	@echo
	@echo "Choose a make command to run in "$(PROJECTNAME)":"
	@echo
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'
	@echo

all:
	make clean
	make fmt
	make lint
	make build
	make test

fmt:
	go fmt ./...

lint:
	golangci-lint run ./...

pre_commit:
	pre-commit run --all-files

get:
	@echo "  >  \033[32mDownloading & Installing all the modules...\033[0m "
	go mod tidy && go mod download

dev_docker:
	@echo "  >  \033[32mStarting docker environment... \033[0m "
	docker-compose -f docker/docker-compose.yml -f docker/docker-compose.dev.yml up

dev_hub:
	@echo "  >  \033[32mHappy coding hub! 😄😄😄 \033[0m "
	go run hub/main.go

dev_indexer:
	@echo "  >  \033[32mHappy coding indexer! 😄😄😄 \033[0m "
	go run hub/indexer.go

build:
	make build_go

build_go:
	@echo "  >  \033[32mBuilding binary...\033[0m "
	go build -o dist/hub ./hub/
	go build -o dist/indexer ./indexer/

build_docker:
	@echo "  >  \033[32mBuilding docker image...\033[0m "
	docker build -t rss3/pregod-hub -f hub.Dockerfile .
	docker build -t rss3/pregod-indexer -f indexer.Dockerfile .

## Runs go test for all packages
test:
	@echo "  >  \033[32mRunning tests...\033[0m "
	go test -p 1 -coverprofile=cover.out -v `go list ./...`

clean:
	rm -rf dist/
