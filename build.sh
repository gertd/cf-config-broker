#!/bin/bash

CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo  .

docker build -t gertd/cf-config-broker .

docker push gertd/cf-config-broker

