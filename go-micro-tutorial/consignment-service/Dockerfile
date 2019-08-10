# We use the official golang image, which contains all the 
# correct build tools and libraries. Notice `as builder`,
# this gives this container a name that we can reference later on. 
# FROM golang:alpine as builder

# RUN apk --no-cache add git

# RUN mkdir -p /go/src/github.com/preslavmihaylov/learn-golang/go-micro-tutorial
# RUN mkdir -p /go/pkg
# RUN mkdir -p /go/bin

# # Set our workdir to our current service in the gopath
# WORKDIR /go/src/github.com/preslavmihaylov/learn-golang/go-micro-tutorial

# # Copy the current code into our workdir
# COPY . .

# RUN export GOPATH=/go
# RUN go get ./...

# # Build the binary, with a few flags
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
#     -o ./consignment-service/consignment-service ./consignment-service

# Here we're using a second FROM statement, which is strange,
# but this tells Docker to start a new build process with this
# image.
FROM alpine:latest

# Security related package, good to have.
RUN apk --no-cache add ca-certificates

# Same as before, create a directory for our app.
RUN mkdir /app
WORKDIR /app

# Here, instead of copying the binary from our host machine,
# we pull the binary from the container named `builder`, within
# this build context. This reaches into our previous image, finds
# the binary we built, and pulls it into this container. Amazing!
COPY consignment-service .

# Run the binary as per usual! This time with a binary build in a
# separate container, with all of the correct dependencies and
# run time libraries.
CMD ["./consignment-service"]
