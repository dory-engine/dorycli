# set current host and all kubernetes nodes to connect {{ $.dory.imageRepo.type }}

- 1. add following {{ $.dory.imageRepo.type }} domain name in /etc/hosts record for current host and all kubernetes nodes

```shell script
vi /etc/hosts
{{ $.imageRepoIp }}  {{ $.imageRepoDomainName }}
```

{{ $certPath := "" }}{{- if eq $.kubernetes.runtime "docker" }}{{ $certPath = "/etc/docker" }}{{- else if eq $.kubernetes.runtime "containerd" }}{{ $certPath = "/etc/containerd" }}{{- else if eq $.kubernetes.runtime "crio" }}{{ $certPath = "/etc/containers" }}{{- end }}
- 2. copy {{ $.dory.imageRepo.type }} server ({{ $.imageRepoIp }}) certificates to all kubernetes nodes {{ $certPath }}/certs.d/{{ $.imageRepoDomainName }} directory

```shell script
{{- if $.imageRepoInternal }}
scp -r {{ $certPath }}/certs.d root@${KUBERNETES_HOST}:{{ $certPath }}/
{{- else }}
# certificates are: ca.crt, {{ $.imageRepoDomainName }}.cert, {{ $.imageRepoDomainName }}.key
{{- end }}


{{- if eq $.kubernetes.runtime "containerd" }}
# set all kubernetes nodes containerd certificates path
vi /etc/containerd/config.toml
      # search following config
      [plugins."io.containerd.grpc.v1.cri".registry.configs]
        ...
        # add harbor ca.crt path settings
        [plugins."io.containerd.grpc.v1.cri".registry.configs."{{ $.imageRepoDomainName }}".tls]
          ca_file = "/etc/containerd/certs.d/{{ $.imageRepoDomainName }}/ca.crt"

# restart all kubernetes nodes containerd service
systemctl restart containerd
{{- end }}
```

# after finish set current host and all kubernetes nodes to connect {{ $.dory.imageRepo.type }}, please input [YES] to go on, input [NO] to cancel