BINARY_NAME=gofind
.PHONY: default build install clean 

OS := $(shell uname -s)

ifeq ($(OS),Linux)
    BUILD_CMD=GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${BINARY_NAME} main.go
endif
ifeq ($(OS),Darwin)
    BUILD_CMD=GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o ${BINARY_NAME} main.go
endif

default: build

build:
	$(BUILD_CMD)

install: build
	sudo cp ${BINARY_NAME} /usr/local/bin/${BINARY_NAME}

clean:
	go clean