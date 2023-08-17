# 设置当前节点以及所有kubernetes节点连接 {{ $.dory.imageRepo.type }}

- 1. 添加以下 {{ $.dory.imageRepo.type }} 域名记录到当前节点以及所有kubernetes节点的 /etc/hosts 文件

```shell script
vi /etc/hosts
{{ $.imageRepoIp }}  {{ $.imageRepoDomainName }}
```

{{ $certPath := "" }}{{- if eq $.kubernetes.runtime "docker" }}{{ $certPath = "/etc/docker" }}{{- else if eq $.kubernetes.runtime "containerd" }}{{ $certPath = "/etc/containerd" }}{{- else if eq $.kubernetes.runtime "crio" }}{{ $certPath = "/etc/containers" }}{{- end }}
- 2. 把{{ $.dory.imageRepo.type }}服务器({{ $.imageRepoIp }})上的证书复制到所有kubernetes节点的 {{ $certPath }}/certs.d/{{ $.imageRepoDomainName }} 目录

```shell script
{{- if $.imageRepoInternal }}
scp -r {{ $certPath }}/certs.d root@${KUBERNETES_HOST}:{{ $certPath }}/
{{- else }}
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
```

# 完成当前节点以及所有kubernetes节点连接 {{ $.dory.imageRepo.type }}后，请输入 [YES] 继续安装，输入 [NO] 取消安装