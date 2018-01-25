PWD := $(shell basename `pwd`)

build:
	$(shell go build)
test:
	go test ${PWD}/cgminer
arm: clean
	$(shell GOOS=linux GOARCH=arm go build)
dep:
	$(shell dep ensure)
clean:
	@rm -f bw_exporter
