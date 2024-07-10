BINARY_NAME=gofind
.PHONY: default build install clean 

default: build

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${BINARY_NAME} main.go

install: build
	sudo cp ${BINARY_NAME} /usr/local/bin/${BINARY_NAME}

clean:
	go clean