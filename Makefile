.PHONY: build vet fmt clean

BUILD_DIR ?= build
MODULE ?= $(shell go list -m)

VERSION := 0.1.0
COMMIT 	:= $(shell git rev-parse --short HEAD)

LDFLAGS :=	-X '$(MODULE)/utils.Commit=$(COMMIT)'	\
			-X '$(MODULE)/utils.Version=$(VERSION)' \
			-s -w

fmt:
	go mod tidy && go fmt ./...

vet: fmt
	go vet ./...

build: vet
	@echo Build directory: $(BUILD_DIR)/
	@echo Linker flags: $(LDFLAGS)
	go build -o $(BUILD_DIR)/ -ldflags="$(LDFLAGS)"

clean:
	go clean -i ./... && rm -r $(BUILD_DIR)
