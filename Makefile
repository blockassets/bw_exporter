DATE=$(shell date -u '+%Y-%m-%d %H:%M:%S')
COMMIT=$(shell git log --format=%h -1)
UNAME=$(shell uname)
VERSION=main.version=${TRAVIS_BUILD_NUMBER} ${COMMIT} ${DATE} ${UNAME}
COMPILE_FLAGS=-ldflags="-X '${VERSION}'"

build:
	@go build ${COMPILE_FLAGS}

arm:
	@GOOS=linux GOARCH=arm go build ${COMPILE_FLAGS}

test:
	@go test ./cgminer

dep:
	@dep ensure

clean:
	@rm -f bw_exporter

all: clean test build
