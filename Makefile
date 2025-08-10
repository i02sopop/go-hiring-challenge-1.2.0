BIN=server

.DEFAULT_GOAL := build

.PHONY: build
build:
	@go build -o $(PWD)/bin/${BIN} $(PWD)/cmd/${BIN}/...

.PHONY: tidy
tidy:
	@go mod tidy && go mod vendor

.PHONY: seed
seed:
	@go run cmd/seed/main.go

.PHONY: run
run: build
	@$(PWD)/bin/${BIN}

.PHONY: test
test:
	@go test -v -count=1 -race ./... -coverprofile=coverage.out -covermode=atomic

.PHONY: test-update
test-update:
	@go test -v -count=1 -race ./... -coverprofile=coverage.out -covermode=atomic -tags=update

.PHONY: docker-up
docker-up:
	docker compose up -d

.PHONY: docker-down
docker-down:
	docker compose down

.PHONY: lint
lint:
	@golangci-lint run -c $(PWD)/.golangci.yaml

.PHONY: lint-fix
lint-fix:
	@golangci-lint run -c $(PWD)/.golangci.yaml --fix

.PHONY: clean
clean:
	@rm -fr ${BIN}* bin cover.out

.PHONY: dependencies
dependencies:
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
