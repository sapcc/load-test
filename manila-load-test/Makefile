all: build

.Phony: init
init:
	go get -u github.com/tsenart/vegeta

.Phony: build
build: 
	go build

.Phony: plot
plot:
	vegeta plot results.bin > plot.html
