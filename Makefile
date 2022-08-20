# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

.PHONY: build
build:
	$(GOBUILD) -o bin/tools -v ./main.go


.PHONY: clean
clean:
	$(GOCLEAN)

.PHONY: test
test:
	$(GOTEST) -v -count=1 ./...
