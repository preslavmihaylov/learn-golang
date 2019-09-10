build:
	protoc -I. --go_out=. --micro_out=. proto/consignment/consignment.proto
	# CGO_ENABLED=0 go build -a -installsuffix cgo -o consignment-service .
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o consignment-service .

clean:
	rm $(shell basename $(CURDIR))
