.PHONY: build vet fmt clean

BIN_DIR ?= build

fmt:
	go mod tidy && go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build -o $(BIN_DIR)/

clean:
	rm -r build