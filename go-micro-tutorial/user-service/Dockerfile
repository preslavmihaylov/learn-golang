FROM alpine:latest

RUN apk --no-cache add ca-certificates bash

RUN mkdir /app
WORKDIR /app

COPY user-service .
COPY wait-for-it.sh .

CMD ["./user-service"]
