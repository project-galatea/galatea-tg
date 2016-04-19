# These are the values we want to pass for Version and BuildTime
VERSION=0.0.1
BUILD_TIME=$(shell date +%s)

# Setup the -ldflags option for go build here, interpolate the variable values

LDFLAGS += -X \"main.Version=$(VERSION)\"
LDFLAGS += -X \"main.BuildTime=$(BUILD_TIME)\"
LDFLAGS += -X \"main.BotToken=$(BOT_TOKEN)\"

.PHONY: build clean

build:
	go build -ldflags "$(LDFLAGS)"

install:
	go install -ldflags "$(LDFLAGS)"

clean:
	go clean -i ./...

clean-mac: clean
	find . -name ".DS_Store" -print0 | xargs -0 rm
