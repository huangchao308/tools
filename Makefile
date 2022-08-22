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
	rm -rf bin/
	rm -rf out/

.PHONY: test
test:
	$(GOTEST) -v -count=1 ./...
