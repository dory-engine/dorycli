# 以docker-compose方式部署dory

## 概要

1. 请根据 `README-0-docker-install.md` 的说明手工安装dory
2. 请根据 `README-1-docker-config.md` 的说明在完成安装后手工设置dory
3. 假如安装失败，请根据 `README-2-docker-reset.md` 的说明停止所有dory服务并重新安装

## 复制所有脚本、配置文件到安装目录

```shell script
# 复制所有脚本、配置文件到安装目录
mkdir -p {{ $.rootDir }}/{{ $.dory.namespace }}/dory-engine/dory-data/certs/openldap
cp -rp * {{ $.rootDir }}
cp -r /usr/share/zoneinfo {{ $.rootDir }}/{{ $.dory.namespace }}/dory-engine/dory-data

{{- if $.artifactRepoInternal }}
# 解压nexus初始化数据
tar Cxzvf {{ $.rootDir }}/{{ $.dory.namespace }} ../{{ $.dory.nexusInitData }}
chown -R 200:200 {{ $.rootDir }}/{{ $.dory.namespace }}/nexus
ls -alh {{ $.rootDir }}/{{ $.dory.namespace }}/nexus
{{- end }}

# 解压trivy漏洞库
tar Cxzvf {{ $.rootDir }}/{{ $.dory.namespace }}/dory-engine/dory-data ../{{ $.dory.trivyDb }}
ls -alh {{ $.rootDir }}/{{ $.dory.namespace }}/dory-engine/dory-data/trivy
```

## {{ $.dory.imageRepo.type }} 安装配置

```shell script
{{- if $.imageRepoInternal }}
# 创建 {{ $.dory.imageRepo.type }} 自签名证书
cd {{ $.rootDir }}/{{ $.dory.imageRepo.internal.namespace }}
sh harbor_certs.sh
ls -alh

# 安装 {{ $.dory.imageRepo.type }}
cd {{ $.rootDir }}/{{ $.dory.imageRepo.internal.namespace }}
chmod a+x common.sh install.sh prepare
sh install.sh
ls -alh

# 停止并更新 {{ $.dory.imageRepo.type }} 的 docker-compose.yml 文件
sleep 5 && docker-compose stop && docker-compose rm -f
export HARBOR_CONFIG_ROOT_PATH=$(echo "{{ $.rootDir }}/{{ $.dory.imageRepo.internal.namespace }}" | sed 's#\/#\\\/#g')
sed -i "s/${HARBOR_CONFIG_ROOT_PATH}/./g" docker-compose.yml
cat docker-compose.yml

# 重启 {{ $.dory.imageRepo.type }}
docker-compose up -d
sleep 10

# 检查 {{ $.dory.imageRepo.type }} 状态
docker-compose ps
{{- end }}
{{ $certPath := "" }}{{- if eq $.kubernetes.runtime "docker" }}{{ $certPath = "/etc/docker" }}{{- else if eq $.kubernetes.runtime "containerd" }}{{ $certPath = "/etc/containerd" }}{{- else if eq $.kubernetes.runtime "crio" }}{{ $certPath = "/etc/containers" }}{{- end }}

{{- if $.imageRepoInternal }}
# 把{{ $.dory.imageRepo.type }}服务器({{ $.imageRepoIp }})上的证书复制到所有kubernetes节点的 {{ $certPath }}/certs.d/{{ $.imageRepoDomainName }} 目录
scp -r /etc/docker/certs.d root@${KUBERNETES_HOST}:{{ $certPath }}/
{{- else }}
# 把{{ $.dory.imageRepo.type }}服务器({{ $.imageRepoIp }})上的证书复制到所有kubernetes节点的 {{ $certPath }}/certs.d/{{ $.imageRepoDomainName }} 目录
# 证书文件包括: ca.crt, {{ $.imageRepoDomainName }}.cert, {{ $.imageRepoDomainName }}.key
{{- end }}

# 在当前主机以及所有kubernetes节点主机上，把 {{ $.dory.imageRepo.type }} 的域名记录添加到 /etc/hosts
vi /etc/hosts
{{ $.imageRepoIp }}  {{ $.imageRepoDomainName }}

# 设置docker客户端登录到 {{ $.dory.imageRepo.type }}
docker login --username {{ $.imageRepoUsername }} --password {{ $.imageRepoPassword }} {{ $.imageRepoDomainName }}

# 在 {{ $.dory.imageRepo.type }} 中创建 public, hub, gcr, quay 四个项目
curl -k -X POST -H 'Content-Type: application/json' -d '{"project_name": "public", "public": true}' 'https://{{ $.imageRepoUsername }}:{{ $.imageRepoPassword }}@{{ $.imageRepoDomainName }}/api/v2.0/projects'
curl -k -X POST -H 'Content-Type: application/json' -d '{"project_name": "hub", "public": true}' 'https://{{ $.imageRepoUsername }}:{{ $.imageRepoPassword }}@{{ $.imageRepoDomainName }}/api/v2.0/projects'
curl -k -X POST -H 'Content-Type: application/json' -d '{"project_name": "gcr", "public": true}' 'https://{{ $.imageRepoUsername }}:{{ $.imageRepoPassword }}@{{ $.imageRepoDomainName }}/api/v2.0/projects'
curl -k -X POST -H 'Content-Type: application/json' -d '{"project_name": "quay", "public": true}' 'https://{{ $.imageRepoUsername }}:{{ $.imageRepoPassword }}@{{ $.imageRepoDomainName }}/api/v2.0/projects'

# 把之前拉取的容器镜像推送到 {{ $.dory.imageRepo.type }}
{{- range $_, $image := $.dockerImages }}
docker tag {{ if $image.dockerFile }}{{ $image.target }}{{ else }}{{ $image.source }}{{ end }} {{ $.imageRepoDomainName }}/{{ $image.target }}
{{- if $image.arm64 }}
docker tag {{ if $image.dockerFile }}{{ $image.target }}-arm64v8{{ else }}{{ $image.arm64 }}{{ end }} {{ $.imageRepoDomainName }}/{{ $image.target }}-arm64v8
{{- end }}
{{- end }}

{{- range $_, $image := $.dockerImages }}
docker push {{ $.imageRepoDomainName }}/{{ $image.target }}
{{- if $image.arm64 }}
docker push {{ $.imageRepoDomainName }}/{{ $image.target }}-arm64v8
{{- end }}
{{- end }}
```

## 使用docker-compose方式安装dory组件

```shell script
# 创建 docker executor 自签名证书
cd {{ $.rootDir }}/{{ $.dory.namespace }}/{{ $.dory.docker.dockerName }}
sh docker_certs.sh
ls -alh

# 创建 openldap 自签名证书
cd {{ $.rootDir }}/{{ $.dory.namespace }}/{{ $.dory.openldap.serviceName }}/certs
sh openldap_certs.sh
ls -alh
cp ca.crt ldap.crt ldap.key {{ $.rootDir }}/{{ $.dory.namespace }}/dory-engine/dory-data/certs/openldap

# 创建 dory 组件目录并设置权限
cd {{ $.rootDir }}/{{ $.dory.namespace }}
mkdir -p mongo-dory
chown -R 999:999 mongo-dory
mkdir -p dory-engine/dory-data
mkdir -p dory-engine/tmp
chown -R 1000:1000 dory-engine
ls -alh

# 使用docker-compose启动所有dory组件
cd {{ $.rootDir }}/{{ $.dory.namespace }}
ls -alh
docker-compose stop && docker-compose rm -f && docker-compose up -d

# 检查dory组件的状态
sleep 10
docker-compose ps
```

## 在kubernetes中创建project-data-alpine pod

```shell script
# project-data-pod pod 用于创建项目的应用文件目录
# 在kubernetes中创建project-data-alpine pod
cd {{ $.rootDir }}
kubectl apply -f project-data-pod.yaml
kubectl -n {{ $.dory.namespace }} get pods
```

## 请继续完成dory的配置

2. 请根据 `README-1-docker-config.md` 的说明在完成安装后手工设置dory
