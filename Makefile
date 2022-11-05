BIN := "./bin"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" .


run-integration-test-pg:
	@docker-compose -f tests/docker-compose.yml -f tests/postgres/docker-compose.pg.local.yml \
		--project-directory tests/postgres up \
		--build --abort-on-container-exit

	@docker-compose -f tests/docker-compose.yml -f tests/postgres/docker-compose.pg.local.yml \
		--project-directory tests/postgres down \
		 --rmi local \
		--volumes \
		--timeout 60;