FROM alpine:3.19.1

LABEL maintainer="cookeem"
LABEL email="cookeem@qq.com"
LABEL version="v1.6.2"

COPY dorycli /usr/bin
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk --update add ca-certificates bash-completion bash git tree htop curl zip jq && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/* && \
    adduser -h /home/dory -s /bin/bash -u 1000 -D dory && \
    mkdir -p /etc/bash_completion.d/
RUN dorycli completion bash > /etc/bash_completion.d/dorycli
COPY .bashrc /home/dory/
COPY .bashrc /root/
WORKDIR /home/dory
USER dory

# docker build -t doryengine/dorycli:v1.6.2-alpine .
# docker push doryengine/dorycli:v1.6.2-alpine

# 创建外部目录保存.dorycli/config.yaml
# mkdir -p .dorycli && sudo chown -R 1000:1000 .dorycli
# docker run -ti --rm -v $PWD/.dorycli:/home/dory/.dorycli doryengine/dorycli:v1.6.2-alpine bash
