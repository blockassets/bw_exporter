DATE=$(shell date -u '+%Y-%m-%d %H:%M:%S')
VERSION=${TRAVIS_BUILD_ID} ${TRAVIS_COMMIT} ${DATE}
COMPILE_FLAGS=-ldflags="-X 'main.version=${VERSION}'"

build:
	go build ${COMPILE_FLAGS}

arm:
	@GOOS=linux GOARCH=arm go build ${COMPILE_FLAGS}

test:
	@go test ${PWD}/cgminer

dep:
	@dep ensure

clean:
	@rm -f bw_exporter

all: clean test build
