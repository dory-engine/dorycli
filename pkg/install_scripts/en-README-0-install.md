# install dory with kubernetes

## summary

1. please follow `README-0-install.md` to install dory by manual
2. please follow `README-1-config.md` to connect dory after install
3. if install fail, please follow `README-2-reset.md` to stop all dory services and install again

{{ if eq $.dory.imageRepo.type "harbor" }}
## load harbor images on all kubernetes hosts

1. please load following images on all kubernetes hosts:
  - {{ $.imageRepoDomainName }}/public/alpine:3.17.2-dory

{{ if $.dory.imageRepo.internal.version }}
2. please load following harbor images on all kubernetes hosts:
  - goharbor/harbor-core:{{ $.dory.imageRepo.internal.version }}
  - goharbor/harbor-db:{{ $.dory.imageRepo.internal.version }}
  - goharbor/harbor-jobservice:{{ $.dory.imageRepo.internal.version }}
  - goharbor/harbor-portal:{{ $.dory.imageRepo.internal.version }}
  - goharbor/harbor-registryctl:{{ $.dory.imageRepo.internal.version }}
  - goharbor/redis-photon:{{ $.dory.imageRepo.internal.version }}
  - goharbor/registry-photon:{{ $.dory.imageRepo.internal.version }}
{{- end }}
{{- end }}

- please ensure that the directory {{ $.kubernetes.pvPath }} in {{ $.kubernetes.pvType }} already exists, otherwise project-data-pod will fail to start

## == create dory components files and copy to kubernetes shared storage, and create relative deployment in kubernetes

```shell script
sh create-dory-files.sh
```

{{ if eq $.dory.imageRepo.type "harbor" }}
## == {{ $.dory.imageRepo.type }} initial settings

```shell script
{{ $certPath := "" }}{{- if eq $.kubernetes.runtime "docker" }}{{ $certPath = "/etc/docker" }}{{- else if eq $.kubernetes.runtime "containerd" }}{{ $certPath = "/etc/containerd" }}{{- else if eq $.kubernetes.runtime "crio" }}{{ $certPath = "/etc/containers" }}{{- end }}
{{- if $.imageRepoInternal }}
# copy {{ $.dory.imageRepo.type }} server ({{ $.imageRepoIp }}) certificates to all kubernetes nodes {{ $certPath }}/certs.d/{{ $.imageRepoDomainName }} directory
scp -r {{ $certPath }}/certs.d root@${KUBERNETES_HOST}:{{ $certPath }}/
{{- else }}
# copy {{ $.dory.imageRepo.type }} server ({{ $.imageRepoIp }}) certificates to all kubernetes nodes {{ $certPath }}/certs.d/{{ $.imageRepoDomainName }} directory
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

# on current host and all kubernetes nodes add {{ $.dory.imageRepo.type }} domain name in /etc/hosts
vi /etc/hosts
{{ $.imageRepoIp }}  {{ $.imageRepoDomainName }}

# push container images to {{ $.dory.imageRepo.type }}
sh push-images.sh
```
{{- end }}

## == install dory services with kubernetes

```shell script
sh deploy-dory.sh
```

{{- if and $.gitRepoInternal (eq $.dory.gitRepo.type "gitea") }}

## == auto config {{ $.dory.gitRepo.type }}

```shell script
{{ $.cmdRun }} --rm -ti -v $PWD:/src doryengine/python:3.11.2-alpine3.17-dory python /src/gitea-config.py
```
{{- else if and $.gitRepoInternal (eq $.dory.gitRepo.type "gitlab") }}

## == auto config {{ $.dory.gitRepo.type }}

```shell script
sh gitlab-config.sh
```
{{- end }}

{{- if $.artifactRepoInternal }}

## == auto config {{ $.dory.artifactRepo.type }}

```shell script
sh nexus-config.sh
```
{{- end }}

{{- if $.scanCodeRepoInternal }}

## == auto config {{ $.dory.scanCodeRepo.type }}

```shell script
sh sonarqube-config.sh
```
{{- end }}

## == restart dory-engine and dory-console

```shell script
sh restart-dory.sh
```

## connect dory

- please follow `README-1-config.md` to connect dory after install
