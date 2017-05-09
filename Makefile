BUILD_TARGET ?= cmd/cfseeker/*.go
OUTPUT_DIR ?= bin
LDFLAGS := -X "github.com/thomasmmitchell/cfseeker/config.Version=$(VERSION)"
OUTPUT_NAME ?= cfseeker
BUILD := go build -v -ldflags='$(LDFLAGS)' -o $(OUTPUT_DIR)/$(OUTPUT_NAME) $(BUILD_TARGET)

.PHONY: build darwin linux all clean
.DEFAULT: build
build:
	GOOS=$(GOOS) GOARCH=amd64 $(BUILD)

darwin:
	GOOS=darwin OUTPUT_NAME=cfseeker-darwin VERSION=$(VERSION) $(MAKE)

linux:
	GOOS=linux OUTPUT_NAME=cfseeker-linux VERSION=$(VERSION) $(MAKE)

all: darwin linux

clean:
	rm -f bin/*

embed:
	go run utils/embed.go assets/web