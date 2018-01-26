DATE=$(shell date -u '+%Y-%m-%d %H:%M:%S')
COMMIT=$(shell git log --format=%h -1)
VERSION=main.version=${TRAVIS_BUILD_NUMBER} ${COMMIT} ${DATE}
COMPILE_FLAGS=-ldflags="-X '${VERSION}'"

build:
	@go build ${COMPILE_FLAGS}

arm:
	GOOS=linux GOARCH=arm GOARM=7 go build ${COMPILE_FLAGS}

test:
	@go test ./cgminer

dep:
	@dep ensure

clean:
	@rm -f bw_exporter

all: clean test build
