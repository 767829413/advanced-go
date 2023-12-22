#!/bin/bash

docker run -d -p 5672:5672 -p 15672:15672 --hostname my-rabbit --name rabbit -e RABBITMQ_DEFAULT_USER=user -e RABBITMQ_DEFAULT_PASS=user rabbitmq:3-management

curl -u user:user -X PUT http://172.93.221.226:15672/api/vhosts/vhost1
curl -u user:user -X PUT http://172.93.221.226:15672/api/vhosts/vhost2
curl -u user:user -X PUT http://172.93.221.226:15672/api/vhosts/vhost3
curl -u user:user -X PUT http://172.93.221.226:15672/api/vhosts/vhost4