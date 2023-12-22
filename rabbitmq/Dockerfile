FROM rabbitmq:3-management

# 复制自定义的初始化脚本到容器中
COPY init.sh /usr/local/bin/

# 设置可执行权限
RUN chmod +x /usr/local/bin/init.sh