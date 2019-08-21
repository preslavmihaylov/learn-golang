FROM alpine:latest

RUN apk --no-cache add ca-certificates bash

RUN mkdir /app
WORKDIR /app

COPY user-cli .
COPY wait-for-it.sh .

CMD ./wait-for-it.sh user-service:50051 -t 5 -- ./user-cli --name="pm" --email="pm@pm.com" --password="test" --company="BBC"
