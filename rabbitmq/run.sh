#!/bin/bash
docker build -t my-rabbitmq .
docker run -d -p 5672:5672 -p 15672:15672 --hostname my-rabbit --name rabbit -e RABBITMQ_DEFAULT_USER=user -e RABBITMQ_DEFAULT_PASS=user my-rabbitmq /usr/local/bin/init.sh

