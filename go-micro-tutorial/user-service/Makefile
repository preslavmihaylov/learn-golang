BINARY=$(shell basename $(CURDIR))

build:
	protoc -I. --go_out=plugins=micro:. proto/user/user.proto
	# CGO_ENABLED=0 go build -a -installsuffix cgo -o user-service .
	CGO_ENABLED=0 go build -o $(BINARY) .

clean:
	rm $(BINARY)
