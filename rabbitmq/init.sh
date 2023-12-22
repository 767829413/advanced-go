#!/bin/bash

# 等待RabbitMQ启动
sleep 10

# 创建vhost1
rabbitmqctl add_vhost vhost1
rabbitmqctl set_permissions -p vhost1 user ".*" ".*" ".*"

# 创建vhost2
rabbitmqctl add_vhost vhost2
rabbitmqctl set_permissions -p vhost2 user ".*" ".*" ".*"

# 其他vhost...

# 保持容器运行
rabbitmq-server