build:
	protoc -I. --go_out=plugins=micro:. proto/vessel/vessel.proto
	# CGO_ENABLED=0 go build -a -installsuffix cgo -o vessel-service .
	CGO_ENABLED=0 go build -o vessel-service .

clean:
	rm $(shell basename $(CURDIR))
