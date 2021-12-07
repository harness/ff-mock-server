ifndef GOPATH
	GOPATH := $(shell go env GOPATH)
endif
ifndef GOBIN
	GOBIN := $(shell go env GOPATH)/bin
endif
ifndef DOCKER_BUILD_OPTS
	DOCKER_BUILD_OPTS := --build
endif

tools = $(addprefix $(GOBIN)/, golangci-lint gosec goimports)
deps = $(addprefix $(GOBIN)/, oapi-codegen)

.DEFAULT_GOAL := all

dep: $(deps) ## Install the deps required to generate code and build mock-server
	@echo "Installing dependences"
	@go mod download

tools: $(tools) ## Install tools required for the build
	@echo "Installed tools"

generate: dep ## Generate code
	@echo "Generating Code"
	@go generate ./pkg/api/generate.go

test: dep generate  ## Run the go tests
	@echo "Running tests"
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

all: dep tools generate check build test ## Build and run the test for mock-server
	@echo "Run `make start`  to start the services"

image: ## Build the server docker image
	@echo "Building Mock Server Image"
	@docker build --build-arg FF_VERSION=latest \
			--build-arg FF_COMMIT=${FF_COMMIT} \
			 -t ff-mock-server:latest \
			 -f ./Dockerfile .

run: dep ## Run the feature flag binary from source
	@go run -race -ldflags="-X github.com/wings-software/ff-mock-server/pkg/version.Version=1.0.0" ./cmd/server

build: dep generate ## Build the service binary
	@echo "Building Mock Server"
	CGO_ENABLED=0 go build -ldflags="-X github.com/wings-software/ff-mock-server/pkg/version.Version=${FF_VERSION}" -o ./cmd/server/server ./cmd/server


#########################################
# Checks
# These lint, format and check the code for potential vulnerabilities
#########################################

check: generate lint format sec

format: tools # Format go code and error if any changes are made
	@echo "Formating ..."
	@goimports -w .
	@echo "Formatting complete"

lint: tools generate # lint the golang code
	@echo "Linting $(1)"
	@golint ./...
	@go vet ./...

sec: tools # Run the security checks
	@echo "Checking for security problems ..."
	@gosec -quiet -confidence high -severity medium ./...
	@echo "No problems found"

###########################################
# Install Tools and deps
#
# These targets specify the full path to where the tool is installed
# If the tool already exists it wont be re-installed.
###########################################

# Install golangci-lint
$(GOBIN)/golangci-lint:
	@echo "🔘 Installing golangci-lint... (`date '+%H:%M:%S'`)"
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin

# Install goimports to format code
$(GOBIN)/goimports:
	@echo "🔘 Installing goimports ... (`date '+%H:%M:%S'`)"
	@go install golang.org/x/tools/cmd/goimports

# Install gosec for security scans
$(GOBIN)/gosec:
	@echo "🔘 Installing gosec ... (`date '+%H:%M:%S'`)"
	@curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b $(GOPATH)/bin

# Install oapi-codegen to generate ff server code from the apis
$(GOBIN)/oapi-codegen:
	@go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.8.3