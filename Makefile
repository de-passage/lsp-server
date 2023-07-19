BINARY_NAME=lsp
BUILD_DIR=build

FULL_BINARY_NAME=$(BUILD_DIR)/$(BINARY_NAME)

.PHONY: all
all: build

.PHONY: build
build: $(FULL_BINARY_NAME)

$(FULL_BINARY_NAME): $(wildcard *.go)
	go build -o $(FULL_BINARY_NAME)

.PHONY: test
test:
	go test -v ./...

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

.PHONY: run
run: build
	$(FULL_BINARY_NAME)
