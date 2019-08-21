BINARY=$(shell basename $(CURDIR))

build:
	CGO_ENABLED=0 go build -o $(BINARY) .

clean:
	rm $(BINARY)
