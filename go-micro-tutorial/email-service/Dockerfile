FROM alpine:latest

RUN apk --no-cache add ca-certificates bash

RUN mkdir /app
WORKDIR /app

COPY email-service .

CMD ./email-service
