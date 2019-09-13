FROM debian:latest

RUN mkdir -p /app
WORKDIR /app

ADD consignment.json .
ADD consignment-cli .

CMD ./consignment-cli
