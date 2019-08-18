#!/bin/bash

make -C consignment-cli build
make -C consignment-service build
make -C vessel-service build

docker-compose build

make -C consignment-cli clean
make -C consignment-service clean
make -C vessel-service clean
