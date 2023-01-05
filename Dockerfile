FROM ubuntu:20.04

WORKDIR /app

# 更换软件源
COPY box/sources.list /etc/apt/sources.list

# 安装工具
RUN apt update && \
  apt install -y vim curl wget net-tools git traceroute

# 安装golang
RUN cd /usr/local && \
  wget https://go.dev/dl/go1.19.linux-amd64.tar.gz && \
  tar -xvf go1.19.linux-amd64.tar.gz && \
  rm go1.19.linux-amd64.tar.gz && \
  ln -s $(pwd)/go/bin/go /usr/local/bin

# 拷贝app、静态文件、配置文件等
COPY eserver /app/eserver
COPY html /app/html
COPY config /app/config

CMD ["/app/eserver"]