BUILD_TARGET ?= cmd/cfseeker/*.go
OUTPUT_DIR ?= bin
OUTPUT_NAME ?= cfseeker
VERSION ?= development build: $(shell date)
LDFLAGS := -X "github.com/cloudfoundry-community/cfseeker/config.Version=$(VERSION)"
BUILD := go build -v -ldflags='$(LDFLAGS)' -o $(OUTPUT_DIR)/$(OUTPUT_NAME) $(BUILD_TARGET)

.PHONY: build darwin linux all clean
.DEFAULT: build
build:
	@echo $(VERSION)
	GOOS=$(GOOS) GOARCH=amd64 $(BUILD)

darwin:
	GOOS=darwin OUTPUT_NAME=cfseeker-darwin VERSION="$(VERSION)" $(MAKE)

linux:
	GOOS=linux OUTPUT_NAME=cfseeker-linux VERSION="$(VERSION)" $(MAKE)

all: darwin linux

clean:
	rm -f bin/*

embed:
	go run utils/embed.go assets/web