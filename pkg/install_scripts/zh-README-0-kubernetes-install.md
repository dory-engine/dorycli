# 以kubernetes方式部署dory

## 概要

1. 请根据 `README-0-kubernetes-install.md` 的说明手工安装dory
2. 请根据 `README-1-kubernetes-config.md` 的说明在完成安装后手工设置dory
3. 假如安装失败，请根据 `README-2-kubernetes-reset.md` 的说明停止所有dory服务并重新安装

## 创建相关根目录

```shell script
{{- if $.imageRepoInternal }}
# 创建 {{ $.dory.imageRepo.type }} 相关目录并设置目录权限
mkdir -p {{ $.rootDir }}/{{ $.dory.imageRepo.internal.namespace }}/database
mkdir -p {{ $.rootDir }}/{{ $.dory.imageRepo.internal.namespace }}/jobservice
mkdir -p {{ $.rootDir }}/{{ $.dory.imageRepo.internal.namespace }}/redis
mkdir -p {{ $.rootDir }}/{{ $.dory.imageRepo.internal.namespace }}/registry
chown -R 999:999 {{ $.rootDir }}/{{ $.dory.imageRepo.internal.namespace }}/database
chown -R 10000:10000 {{ $.rootDir }}/{{ $.dory.imageRepo.internal.namespace }}/jobservice
chown -R 999:999 {{ $.rootDir }}/{{ $.dory.imageRepo.internal.namespace }}/redis
chown -R 10000:10000 {{ $.rootDir }}/{{ $.dory.imageRepo.internal.namespace }}/registry
ls -alh {{ $.rootDir }}/{{ $.dory.imageRepo.internal.namespace }}
{{- end }}

# 创建 openldap 的自签名证书
mkdir -p {{ $.rootDir }}/{{ $.dory.namespace }}/dory-engine/dory-data/certs/openldap
cd {{ $.dory.namespace }}/{{ $.dory.openldap.serviceName }}/certs
sh openldap_certs.sh
ls -alh
cp ca.crt ldap.crt ldap.key {{ $.rootDir }}/{{ $.dory.namespace }}/dory-engine/dory-data/certs/openldap
cd ../../.. 

# 创建 dory 组件相关目录并设置目录权限
mkdir -p {{ $.rootDir }}/{{ $.dory.namespace }}/dory-engine/logs
cp -rp {{ $.dory.namespace }}/{{ $.dory.docker.dockerName }} {{ $.rootDir }}/{{ $.dory.namespace }}/
cp -rp {{ $.dory.namespace }}/{{ $.dory.openldap.serviceName }} {{ $.rootDir }}/{{ $.dory.namespace }}/
cp -rp {{ $.dory.namespace }}/dory-engine {{ $.rootDir }}/{{ $.dory.namespace }}/
{{- if and (eq $.dory.gitRepo.type "gitlab") $.dory.gitRepo.internal.image }}
cp -rp {{ $.dory.namespace }}/nginx-gitlab {{ $.rootDir }}/{{ $.dory.namespace }}/
{{- end }}
{{- if and (eq $.dory.scanCodeRepo.type "sonarqube") $.dory.scanCodeRepo.internal.image }}
mkdir -p {{ $.rootDir }}/{{ $.dory.namespace }}/sonarqube-web/data
mkdir -p {{ $.rootDir }}/{{ $.dory.namespace }}/sonarqube-web/extensions
mkdir -p {{ $.rootDir }}/{{ $.dory.namespace }}/sonarqube-web/logs
mkdir -p {{ $.rootDir }}/{{ $.dory.namespace }}/sonarqube-web/temp
chown -R 1000:1000 {{ $.rootDir }}/{{ $.dory.namespace }}/sonarqube-web
{{- end }}
mkdir -p {{ $.rootDir }}/{{ $.dory.namespace }}/mongo-dory
chown -R 999:999 {{ $.rootDir }}/{{ $.dory.namespace }}/mongo-dory
ls -alh {{ $.rootDir }}/{{ $.dory.namespace }}
chown -R 1000:1000 {{ $.rootDir }}/{{ $.dory.namespace }}/dory-engine
```

{{ $certPath := "" }}{{- if eq $.kubernetes.runtime "docker" }}{{ $certPath = "/etc/docker" }}{{- else if eq $.kubernetes.runtime "containerd" }}{{ $certPath = "/etc/containerd" }}{{- else if eq $.kubernetes.runtime "crio" }}{{ $certPath = "/etc/containers" }}{{- end }}
{{- if eq $.dory.imageRepo.type "harbor" }}
## {{ $.dory.imageRepo.type }} 安装配置

```shell script
{{- if $.imageRepoInternal }}
# 创建 {{ $.dory.imageRepo.type }} 名字空间与pv
kubectl delete ns {{ $.dory.imageRepo.internal.namespace }}
kubectl delete pv {{ $.dory.imageRepo.internal.namespace }}-pv
kubectl apply -f {{ $.dory.imageRepo.internal.namespace }}/step01-namespace-pv.yaml

# 使用helm安装 {{ $.dory.imageRepo.type }}
helm install -n {{ $.dory.imageRepo.internal.namespace }} {{ $.dory.imageRepo.internal.namespace }} {{ $.dory.imageRepo.type }}
helm -n {{ $.dory.imageRepo.internal.namespace }} list

# 等待所有 {{ $.dory.imageRepo.type }} 服务状态为 ready
kubectl -n {{ $.dory.imageRepo.internal.namespace }} get pods -o wide

# 创建 {{ $.dory.imageRepo.type }} 自签名证书并复制到 {{ $certPath }}/certs.d
sh {{ $.dory.imageRepo.internal.namespace }}/harbor_update_docker_certs.sh
ls -alh {{ $certPath }}/certs.d/{{ $.imageRepoDomainName }}
{{- end }}

{{- if $.imageRepoInternal }}
# 把{{ $.dory.imageRepo.type }}服务器({{ $.imageRepoIp }})上的证书复制到所有kubernetes节点的 {{ $certPath }}/certs.d/{{ $.imageRepoDomainName }} 目录
scp -r {{ $certPath }}/certs.d root@${KUBERNETES_HOST}:{{ $certPath }}/
{{- else }}
# 把{{ $.dory.imageRepo.type }}服务器({{ $.imageRepoIp }})上的证书复制到所有kubernetes节点的 {{ $certPath }}/certs.d/{{ $.imageRepoDomainName }} 目录
# 证书文件包括: ca.crt, {{ $.imageRepoDomainName }}.cert, {{ $.imageRepoDomainName }}.key
{{- end }}

# 在当前主机以及所有kubernetes节点主机上，把 {{ $.dory.imageRepo.type }} 的域名记录添加到 /etc/hosts
vi /etc/hosts
{{ $.imageRepoIp }}  {{ $.imageRepoDomainName }}

# 设置 {{ $.kubernetes.runtime }} 客户端登录到 {{ $.dory.imageRepo.type }}
{{ $.cmdLogin }} --username {{ $.imageRepoUsername }} --password {{ $.imageRepoPassword }} {{ $.imageRepoDomainName }}

# 在 {{ $.dory.imageRepo.type }} 中创建 public, hub, gcr, quay 四个项目
curl -k -X POST -H 'Content-Type: application/json' -d '{"project_name": "public", "public": true}' 'https://{{ $.imageRepoUsername }}:{{ $.imageRepoPassword }}@{{ $.imageRepoDomainName }}/api/v2.0/projects'
curl -k -X POST -H 'Content-Type: application/json' -d '{"project_name": "hub", "public": true}' 'https://{{ $.imageRepoUsername }}:{{ $.imageRepoPassword }}@{{ $.imageRepoDomainName }}/api/v2.0/projects'
curl -k -X POST -H 'Content-Type: application/json' -d '{"project_name": "gcr", "public": true}' 'https://{{ $.imageRepoUsername }}:{{ $.imageRepoPassword }}@{{ $.imageRepoDomainName }}/api/v2.0/projects'
curl -k -X POST -H 'Content-Type: application/json' -d '{"project_name": "quay", "public": true}' 'https://{{ $.imageRepoUsername }}:{{ $.imageRepoPassword }}@{{ $.imageRepoDomainName }}/api/v2.0/projects'

# 把之前拉取的容器镜像推送到 {{ $.dory.imageRepo.type }}
{{- range $_, $image := $.dockerImages }}
{{ $.cmdTag }} {{ if $image.dockerFile }}{{ $image.target }}{{ else }}{{ $image.source }}{{ end }} {{ $.imageRepoDomainName }}/{{ $image.target }}
{{- if $image.arm64 }}
{{ $.cmdTag }} {{ if $image.dockerFile }}{{ $image.target }}-arm64v8{{ else }}{{ $image.arm64 }}{{ end }} {{ $.imageRepoDomainName }}/{{ $image.target }}-arm64v8
{{- end }}
{{- end }}

{{- range $_, $image := $.dockerImages }}
{{ $.cmdPush }} {{ $.imageRepoDomainName }}/{{ $image.target }}
{{- if $image.arm64 }}
{{ $.cmdPush }} {{ $.imageRepoDomainName }}/{{ $image.target }}-arm64v8
{{- end }}
{{- end }}
```
{{- end }}

## 把dory组件部署到kubernetes中

```shell script
# 创建 {{ $.dory.namespace }} 组件的名字空间与pv
kubectl delete ns {{ $.dory.namespace }}
kubectl delete pv {{ $.dory.namespace }}-pv
kubectl apply -f {{ $.dory.namespace }}/step01-namespace-pv.yaml

# 创建 docker executor 自签名证书
sh {{ $.dory.namespace }}/{{ $.dory.docker.dockerName }}/docker_certs.sh
kubectl -n {{ $.dory.namespace }} create secret generic {{ $.dory.docker.dockerName }}-tls --from-file=certs/ca.crt --from-file=certs/tls.crt --from-file=certs/tls.key --dry-run=client -o yaml | kubectl apply -f -
kubectl -n {{ $.dory.namespace }} describe secret {{ $.dory.docker.dockerName }}-tls
rm -rf certs

# 复制容器镜像仓库证书到 docker executor 配置目录
cp -rp {{ $certPath }}/certs.d/{{ $.imageRepoDomainName }} {{ $.rootDir }}/{{ $.dory.namespace }}/{{ $.dory.docker.dockerName }}

{{- if $.artifactRepoInternal }}
# 解压nexus初始化数据
tar Cxzvf {{ $.rootDir }}/{{ $.dory.namespace }} ../{{ $.dory.nexusInitData }}
chown -R 200:200 {{ $.rootDir }}/{{ $.dory.namespace }}/nexus
ls -alh {{ $.rootDir }}/{{ $.dory.namespace }}/nexus
{{- end }}

# 在kubernetes中部署dory组件
kubectl apply -f {{ $.dory.namespace }}/step02-statefulset.yaml
kubectl apply -f {{ $.dory.namespace }}/step03-service.yaml
kubectl apply -f {{ $.dory.namespace }}/step04-networkpolicy.yaml

# 检查dory服务状态
kubectl -n {{ $.dory.namespace }} get pods -o wide
```

## 在kubernetes中创建project-data-alpine pod

```shell script
# project-data-pod pod 用于创建项目的应用文件目录
# 在kubernetes中创建project-data-alpine pod
kubectl apply -f project-data-pod.yaml
kubectl -n {{ $.dory.namespace }} get pods
```

## 请继续完成dory的配置

2. 请根据 `README-1-kubernetes-config.md` 的说明在完成安装后手工设置dory
