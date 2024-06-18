-include .makerc
SHELL = /bin/bash
PACKAGE ?= family-budget-api

#Build server placeholders (BUILD_NUMBER and BRANCH_NAME are overwritten by Jenkins)
BUILD_NUMBER ?= 1
BRANCH_NAME ?= local
BRANCH_NAME_CLEAN = $(shell echo '$(BRANCH_NAME)' | tr "[:upper:]" "[:lower:]" | tr "/" "-")
VERSION ?= $(shell sh build/version.sh $(BRANCH_NAME_CLEAN) $(BUILD_NUMBER) 2>/dev/null)

#Static
GIT_COMMIT ?= $(shell git rev-parse HEAD 2> /dev/null)
GIT_AUTHORS ?= $(shell git log --format='%aN' | sort -u | awk -vORS=, '{ print }' | sed 's/,$$//')

#Output purposes
OUTPUT_DIR = $(CURDIR)/output
BIN_OUTPUT_DIR = $(OUTPUT_DIR)/bin
TEST_OUTPUT_DIR = $(OUTPUT_DIR)/test
DIRS=$(BIN_OUTPUT_DIR) $(TEST_OUTPUT_DIR)
$(shell mkdir -p $(DIRS))

#Docker buildkit setup
DOCKER_BUILDKIT=1
BUILDKIT_INLINE_CACHE=1
BUILDKIT_PROGRESS=plain

#Build flags
# MAIN_PACKAGE: "main" package in Go, leave empty if there is no main package (for a lib)
MAIN_PACKAGE = ./cmd/app
LDFLAGS ?= "-X 'main.version=$(VERSION)' -X 'main.gitCommit=$(GIT_COMMIT)' -X 'main.application=$(PACKAGE)'"
# EXTRA_BUILD_FLAGS and EXTRA_LINT_FLAGS can be used to make a difference between Dockerfile builds and local builds
BUILD_FLAGS ?= $(EXTRA_BUILD_FLAGS) $(MAIN_PACKAGE)
LINT_FLAGS ?= -c ./.golangci.yaml --out-format checkstyle $(EXTRA_LINT_FLAGS)
TEST_FLAGS ?= "-tags=unit"
PACKAGE_EXTENSION ?= $(shell if [ "$(GOOS)" = windows ]; then echo .exe; fi)
GOPROXY=https://dev-athens.be-mobile.biz
GONOSUMDB=bitbucket.org/be-mobile
GO111MODULE=on
CGO_ENABLED=0

#Docker
DOCKER_BUILD_ARGS ?=--build-arg ARG_GIT_COMMIT=$(GIT_COMMIT) --build-arg ARG_VERSION=$(VERSION) --build-arg ARG_AUTHORS="$(GIT_AUTHORS)"
DOCKER_FILE_PATH ?= ./build/docker/Dockerfile

.NOTPARALLEL: ; # wait for this target to finish
.EXPORT_ALL_VARIABLES: ; # send all vars to shell
.PHONY: version build docs test scripts api cmd configs examples

deps: ## Add dependencies for your project

version: ## Return version
	@echo $(VERSION)

help: ## Show Help
	@echo "Usage:"
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort |\
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

lint: ## Lint
	golangci-lint run $(LINT_FLAGS) | tee $(TEST_OUTPUT_DIR)/lint-report.xml
	test -s $(TEST_OUTPUT_DIR)/lint-report.xml  # Check that lint output is not empty

build: ## Build the app
	go build -o $(BIN_OUTPUT_DIR)/$(PACKAGE)$(PACKAGE_EXTENSION) --ldflags=$(LDFLAGS) $(BUILD_FLAGS)

run: ## Run the app
	$(BIN_OUTPUT_DIR)/$(PACKAGE)$(PACKAGE_EXTENSION)

dev:
	go run cmd/app/main.go

test: ## Run tests
	go test $(TEST_FLAGS) ./... $(BUILD_FLAGS)

benchmark: ## Run benchmark tests
	go test -bench=. $(TEST_FLAGS) ./... $(BUILD_FLAGS)

test-report: ## Launch tests and output go-junit-report
	go test -v $(TEST_FLAGS) ./... $(BUILD_FLAGS) > $(TEST_OUTPUT_DIR)/tests.output; cat $(TEST_OUTPUT_DIR)/tests.output
	cat $(TEST_OUTPUT_DIR)/tests.output | go-junit-report > $(TEST_OUTPUT_DIR)/go_test_report.xml

test-coverage: ## Run coverage tool for go
	go test $(TEST_FLAGS) ./... -coverpkg ./... -coverprofile $(TEST_OUTPUT_DIR)/cover.out $(BUILD_FLAGS)
	gocov convert $(TEST_OUTPUT_DIR)/cover.out | gocov-xml -source=$(CURDIR) > $(TEST_OUTPUT_DIR)/coverage.xml

docker-build-report: ## Build: lint and test report
	docker build --target report --output type=local,dest=output $(DOCKER_BUILD_ARGS) -t $(PACKAGE):$(VERSION) -f $(DOCKER_FILE_PATH) .

docker-build: docker-build-report ## Build: docker
	docker build --target final $(DOCKER_BUILD_ARGS) -t $(PACKAGE):$(VERSION) -f $(DOCKER_FILE_PATH) .
