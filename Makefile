BINARY_NAME=stratagem-zero
MAIN_PATH=./cmd/stratagem-zero

.PHONY: all build run test lint clean install

all: build

build:
	go build -o $(BINARY_NAME) $(MAIN_PATH)

run: build
	./$(BINARY_NAME)

test:
	go test -v ./...

lint:
	go vet ./...

clean:
	rm -f $(BINARY_NAME)

install:
	go install $(MAIN_PATH)
