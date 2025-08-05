# 以kubernetes方式部署dory

## 概要

1. 请根据 `README-0-install.md` 的说明手工安装dory
2. 请根据 `README-1-config.md` 的说明在完成安装后访问dory
3. 假如安装失败，请根据 `README-2-reset.md` 的说明停止所有dory服务并重新安装

{{ if eq $.dory.imageRepo.type "harbor" }}
## 把harbor镜像复制到kubernetes所有节点

1. 请在所有kubernetes节点上手工加载以下镜像:
  - {{ $.imageRepoDomainName }}/public/alpine:3.17.2-dory

{{ if $.dory.imageRepo.internal.version }}
2. 请在所有kubernetes节点上手工加载以下harbor镜像:
  - goharbor/harbor-core:{{ $.dory.imageRepo.internal.version }}
  - goharbor/harbor-db:{{ $.dory.imageRepo.internal.version }}
  - goharbor/harbor-jobservice:{{ $.dory.imageRepo.internal.version }}
  - goharbor/harbor-portal:{{ $.dory.imageRepo.internal.version }}
  - goharbor/harbor-registryctl:{{ $.dory.imageRepo.internal.version }}
  - goharbor/redis-photon:{{ $.dory.imageRepo.internal.version }}
  - goharbor/registry-photon:{{ $.dory.imageRepo.internal.version }}
{{- end }}
{{- end }}

- 请保证 {{ $.kubernetes.pvType }} 中目录 {{ $.kubernetes.pvPath }} 已经存在，否则project-data-pod会无法启动

## == 创建dory组件相关文件并复制到kubernetes的共享存储，以及在kubernetes创建相关部署

```shell script
sh create-dory-files.sh
```

{{ if eq $.dory.imageRepo.type "harbor" }}
## == {{ $.dory.imageRepo.type }} 初始化配置

```shell script
{{ $certPath := "" }}{{- if eq $.kubernetes.runtime "docker" }}{{ $certPath = "/etc/docker" }}{{- else if eq $.kubernetes.runtime "containerd" }}{{ $certPath = "/etc/containerd" }}{{- else if eq $.kubernetes.runtime "crio" }}{{ $certPath = "/etc/containers" }}{{- end }}
{{- if $.imageRepoInternal }}
# 把{{ $.dory.imageRepo.type }}服务器({{ $.imageRepoIp }})上的证书复制到所有kubernetes节点的 {{ $certPath }}/certs.d/{{ $.imageRepoDomainName }} 目录
scp -r {{ $certPath }}/certs.d root@${KUBERNETES_HOST}:{{ $certPath }}/
{{- else }}
# 把{{ $.dory.imageRepo.type }}服务器({{ $.imageRepoIp }})上的证书复制到所有kubernetes节点的 {{ $certPath }}/certs.d/{{ $.imageRepoDomainName }} 目录
# 证书文件包括: ca.crt, {{ $.imageRepoDomainName }}.cert, {{ $.imageRepoDomainName }}.key
{{- end }}

{{- if eq $.kubernetes.runtime "containerd" }}

# 设置所有kubernetes节点的containerd的证书路径
vi /etc/containerd/config.toml
      # 寻找以下路径
      [plugins."io.containerd.grpc.v1.cri".registry.configs]
        ...
        # 添加ca.crt的证书路径
        [plugins."io.containerd.grpc.v1.cri".registry.configs."{{ $.imageRepoDomainName }}".tls]
          ca_file = "/etc/containerd/certs.d/{{ $.imageRepoDomainName }}/ca.crt"

# 重启所有kubernetes节点的containerd服务
systemctl restart containerd
{{- end }}

# 在当前主机以及所有kubernetes节点主机上，把 {{ $.dory.imageRepo.type }} 的域名记录添加到 /etc/hosts
vi /etc/hosts
{{ $.imageRepoIp }}  {{ $.imageRepoDomainName }}

# 把之前拉取的容器镜像推送到 {{ $.dory.imageRepo.type }}
sh push-images.sh
```
{{- end }}

## == 把dory组件部署到kubernetes中

```shell script
sh deploy-dory.sh
```

{{- if and $.gitRepoInternal (eq $.dory.gitRepo.type "gitea") }}

## == 自动配置 {{ $.dory.gitRepo.type }}

```shell script
{{ $.cmdRun }} --rm -ti -v $PWD:/src doryengine/python:3.11.2-alpine3.17-dory python /src/gitea-config.py
```
{{- else if and $.gitRepoInternal (eq $.dory.gitRepo.type "gitlab") }}

## == 自动配置 {{ $.dory.gitRepo.type }}

```shell script
sh gitlab-config.sh
```
{{- end }}

{{- if $.artifactRepoInternal }}

## == 自动配置 {{ $.dory.artifactRepo.type }}

```shell script
sh nexus-config.sh
```
{{- end }}

{{- if $.scanCodeRepoInternal }}

## == 自动配置 {{ $.dory.scanCodeRepo.type }}

```shell script
sh sonarqube-config.sh
```
{{- end }}

## == 重启 dory-engine 和 dory-console 服务

```shell script
sh restart-dory.sh
```

## 访问dory

- 请根据 `README-1-config.md` 的说明在完成安装后访问dory
