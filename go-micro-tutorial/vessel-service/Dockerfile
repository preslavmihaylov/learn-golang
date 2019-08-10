# FROM golang:alpine as builder

# RUN apk --no-cache add git

# RUN mkdir -p /go/src/github.com/preslavmihaylov/learn-golang/go-micro-tutorial
# RUN mkdir -p /go/pkg
# RUN mkdir -p /go/bin

# WORKDIR /go/src/github.com/preslavmihaylov/learn-golang/go-micro-tutorial

# COPY . .

# RUN export GOPATH=/go
# RUN go get ./...
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./vessel-service/vessel-service ./vessel-service

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app
COPY vessel-service .

CMD ["./vessel-service"]
