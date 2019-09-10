BINARY=$(shell basename $(CURDIR))

build:
	protoc -I. --go_out=. --micro_out=. proto/user/user.proto
	# CGO_ENABLED=0 go build -a -installsuffix cgo -o user-service .
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BINARY) .

clean:
	rm $(BINARY)
