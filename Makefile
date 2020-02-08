VERSION=0.1
COMMIT=$(shell git rev-parse --short HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
LDFLAGS="-X config.Version=${VERSION} -X config.Commit=${COMMIT} -X config.Branch=${BRANCH}"

GOCMD=go
GOBUILD=$(GOCMD) build -ldflags ${LDFLAGS}
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY=bid-tracker
DOCKER_IMAGE=${BINARY}:latest

SRC = cmd/api/*.go

.PHONY: all help test test-race bench vendor build run run-demo run-demo-race clean docker-build docker-run docker-stop

## ----------------------------------------------------------------------
## Help: Makefile for app: bid-tracker
## ----------------------------------------------------------------------

help:               ## Show this help (default)
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST)

test:               ## Run tests
	$(GOTEST) -timeout 15s -cover -covermode=atomic -v ./...

test-race:          ## Run tests with race detector
	$(GOTEST) -race -v ./...

bench:              ## Run benchmarks
	$(GOTEST) -benchmem -bench=. -v ./...

vendor:             ## Download and tidy go depenencies
	@go mod tidy

build:              ## Build app
	$(GOBUILD) -o $(BINARY) $(SRC)

run:                ## Run app with empty db
	$(GOCMD) run $(SRC)

run-demo:           ## Run app with demo data
	$(GOCMD) run -race $(SRC) -demo

run-demo-race:      ## Run app with demo data and -race flag
	$(GOBUILD) -o $(BINARY) -race $(SRC)
	$(GOCMD) run -race $(SRC) -demo

clean:              ## Remove compiled binary
	rm -f $(BINARY)

docker-build:       ## Build Docker image
	docker build -t $(DOCKER_IMAGE) -f build/Dockerfile .

docker-run:         ## Run bid tracker in Docker
	-@docker rm $(BINARY)
	docker run --name $(BINARY) --rm -d \
		-p 8080:9000 \
		$(DOCKER_IMAGE)

docker-stop:        ## Stop bid-tracker Docker
	-@docker stop $(BINARY)
	-@docker rm $(BINARY)
