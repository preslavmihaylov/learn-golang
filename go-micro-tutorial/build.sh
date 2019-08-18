#!/bin/bash

make -C consignment-cli
make -C consignment-service
make -C vessel-service
docker-compose build
