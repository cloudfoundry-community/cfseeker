BUILD_TARGET ?= cmd/cfseeker/*.go
OUTPUT_DIR ?= bin

build:
	go build -v -o $(OUTPUT_DIR)/cfseeker $(BUILD_TARGET)

darwin:
	GOOS=darwin GOARCH=amd64 go build -v -o $(OUTPUT_DIR)/cfseeker-darwin $(BUILD_TARGET)

linux:
	GOOS=linux GOARCH=amd64 go build -v -o $(OUTPUT_DIR)/cfseeker-linux $(BUILD_TARGET)

all: darwin linux

clean:
	rm -f bin/*