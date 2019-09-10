build:
	protoc -I. --go_out=. --micro_out=. proto/vessel/vessel.proto
	# CGO_ENABLED=0 go build -a -installsuffix cgo -o vessel-service .
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o vessel-service .

clean:
	rm $(shell basename $(CURDIR))
