BINARY=$(shell basename $(CURDIR))

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BINARY) .

clean:
	rm $(BINARY)
