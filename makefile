# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_FOLDER=bin
BINARY_NAME=spark

all: test build
build: 
				$(GOBUILD) -o $(BINARY_FOLDER)/$(BINARY_NAME) -v
test: 
				$(GOTEST) -v ./...
clean: 
				$(GOCLEAN)
				rm -rf $(BINARY_FOLDER)
run:
				$(GOBUILD) -o $(BINARY_FOLDER)/$(BINARY_NAME) -v ./...
				./$(BINARY_FOLDER)/$(BINARY_NAME)
deps:
				$(GOGET) gopkg.in/yaml.v2
