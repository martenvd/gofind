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
	sudo cp gfdir.sh /usr/local/bin/gfdir
	sudo chmod 755 /usr/local/bin/gfdir
	@if ! grep -q "alias gf=" ~/.zshrc; then \
        echo "alias gf='source /usr/local/bin/gfdir'" >> ~/.zshrc; \
		echo "Alias 'gf' added to ~/.zshrc"; \
		echo "Please source your .zshrc file!"; \
    fi
	@if ! grep -q "alias gf=" ~/.bashrc; then \
		echo "alias gf='source /usr/local/bin/gfdir'" >> ~/.bashrc; \
		echo "Alias 'gf' added to ~/.bashrc"; \
		echo "Please source your .bashrc file!"; \
	fi

clean:
	go clean