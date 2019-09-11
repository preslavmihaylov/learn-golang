#!/bin/bash

make -C consignment-cli build
make -C consignment-service build
make -C vessel-service build
make -C user-cli build
make -C user-service build
make -C email-service build

docker-compose down
docker-compose build

make -C consignment-cli clean
make -C consignment-service clean
make -C vessel-service clean
make -C user-cli clean
make -C user-service clean
make -C email-service clean
