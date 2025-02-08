FROM alpine:3.18.2

LABEL maintainer="cookeem"
LABEL email="cookeem@qq.com"
LABEL version="v1.6.6"

COPY dorycli /usr/bin
RUN apk --update add ca-certificates bash-completion bash git tree htop curl zip jq && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/* && \
    adduser -h /home/dory -s /bin/bash -u 1000 -D dory && \
    mkdir -p /etc/bash_completion.d/
RUN dorycli completion bash > /etc/bash_completion.d/dorycli
COPY .bashrc /home/dory/
COPY .bashrc /root/
WORKDIR /home/dory
USER dory

# docker rmi doryengine/dorycli:v1.6.6-alpine
# docker build --platform linux/amd64 -t doryengine/dorycli:v1.6.6-alpine .
# docker push doryengine/dorycli:v1.6.6-alpine

# 创建外部目录保存.dorycli/config.yaml
# mkdir -p .dorycli && sudo chown -R 1000:1000 .dorycli
# docker run -ti --rm -v $PWD/.dorycli:/home/dory/.dorycli doryengine/dorycli:v1.6.6-alpine bash

# docker save -o dorycli__v1.6.6-alpine doryengine/dorycli:v1.6.6-alpine
# scp -r dorycli__v1.6.6-alpine root@itdev-master03:/root/docker-images/
# scp -r dorycli__v1.6.6-alpine root@gditdev-master03:/root/docker-images/
