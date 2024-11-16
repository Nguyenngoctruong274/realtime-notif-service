NAME         =ads-api-go
MAIN_FILE    =main.go
BIN_DIR      =bin
LINT_DIR_LIST = $(shell ls -d */ | grep -v -E scripts\|vendor\|log\|bin\|docs\|mocks/)

.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

build:
	@echo "STEP: BUILD"
	@echo "   1. create dir: $(BIN_DIR)" \
		&& mkdir -p $(BIN_DIR)\
		&& echo "   ==> ok"
	@echo "   2. build: $(MAIN_FILE)" \
		&& go build -o $(BIN_DIR)/$(NAME) $(MAIN_FILE) \
		&& echo "   ==> ok: SERVICE=$(BIN_DIR)/$(NAME)"

## run
run:
	$(BIN_DIR)/$(NAME) api

## start
start: build run

## cron
cron: build
	$(BIN_DIR)/$(NAME) cron
	
## clean
clean:
	@echo "STEP: CLEAN"
	@echo "   1. remove dir: $(BIN_DIR)"
	@rm -rf bin \
	 	&& echo "   ==> ok"

## swagger
swagger:
	swag init -g main.go

# install swagger before running
## swagger-install
swagger-install:
	go install github.com/swaggo/swag/cmd/swag@latest
	export PATH=$$PATH:$$(go env GOPATH)/bin

## test-coverage
test-coverage:
	go test -coverprofile=coverage.out -covermode count ./internal/usecase
	go tool cover -func=coverage.out

## cobertura-report
cobertura-report: $(GOCOVER_COBERTURA)
	echo "Begin generate coverage.xml"
	gocover-cobertura < coverage.out > coverage.xml
	echo "Finish generate coverage.xml"

## $(GOCOVER_COBERTURA)
$(GOCOVER_COBERTURA):
	go get github.com/boumenot/gocover-cobertura

## check-lint
check-lint:
	@echo "STEP: CHECK INDEX"
	$(eval INDEX = 0)
	$(eval COUNT = $(words $(LINT_DIR_LIST)))
	@echo ╔═══ START LINTING
	@for dir in $(LINT_DIR_LIST); do \
		INDEX=$$(($${INDEX}+1)); \
		echo ╠═ Checking [$$INDEX/$(COUNT)] $$dir; \
		golangci-lint run ./$$dir/... --timeout=5m; \
	done
	@echo ╚═══ PASSED LINTING

## test
test:
	go test ./...

## mockgen-repo
mockgen-repo:
	@read -p "Enter service name:" service; \
	read -p "Enter repo name:" repo; \
	mockgen -source=service/$$service/repository/$$repo.go -destination=mocks/repo/$$service/$${repo}_mock.go \
		&& echo "Begin Mockgen"
	@echo "   ==> ok"

## mockgen-repo
mockgen-http:
	@read -p "Enter service name:" service; \
	read -p "Enter go file name:" file; \
	mockgen -source=pkg/xservice/$$service/$$file.go -destination=mocks/http/$$service/$${service}_client_mock.go \
		&& echo "Begin Mockgen"
	@echo "   ==> ok"
