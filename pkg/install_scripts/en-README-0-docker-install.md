# install dory with docker-compose

## summary

1. please follow `README-0-docker-install.md` to install dory by manual
2. please follow `README-1-docker-config.md` to config dory by manual after install
3. if install fail, please follow `README-2-docker-reset.md` to stop all dory services and install again

## copy all scripts and config files to install root directory

```shell script
# copy all scripts and config files to install root directory
mkdir -p {{ $.rootDir }}/{{ $.dory.namespace }}/dory-engine/dory-data/certs/openldap
cp -rp * {{ $.rootDir }}
cp -r /usr/share/zoneinfo {{ $.rootDir }}/{{ $.dory.namespace }}/dory-engine/dory-data

{{- if $.artifactRepoInternal }}
# extract nexus init data
tar Cxzvf {{ $.rootDir }}/{{ $.dory.namespace }} ../{{ $.dory.nexusInitData }}
chown -R 200:200 {{ $.rootDir }}/{{ $.dory.namespace }}/nexus
ls -alh {{ $.rootDir }}/{{ $.dory.namespace }}/nexus
{{- end }}
```

{{ $certPath := "" }}{{- if eq $.kubernetes.runtime "docker" }}{{ $certPath = "/etc/docker" }}{{- else if eq $.kubernetes.runtime "containerd" }}{{ $certPath = "/etc/containerd" }}{{- else if eq $.kubernetes.runtime "crio" }}{{ $certPath = "/etc/containers" }}{{- end }}
{{- if eq $.dory.imageRepo.type "harbor" }}
## install {{ $.dory.imageRepo.type }}

```shell script
{{- if $.imageRepoInternal }}
# create {{ $.dory.imageRepo.type }} certificates
cd {{ $.rootDir }}/{{ $.dory.imageRepo.internal.namespace }}
sh harbor_certs.sh
ls -alh

# install {{ $.dory.imageRepo.type }}
cd {{ $.rootDir }}/{{ $.dory.imageRepo.internal.namespace }}
chmod a+x common.sh install.sh prepare
sh install.sh
ls -alh

# stop and update {{ $.dory.imageRepo.type }} docker-compose.yml
sleep 5 && docker-compose stop && docker-compose rm -f
export HARBOR_CONFIG_ROOT_PATH=$(echo "{{ $.rootDir }}/{{ $.dory.imageRepo.internal.namespace }}" | sed 's#\/#\\\/#g')
sed -i "s/${HARBOR_CONFIG_ROOT_PATH}/./g" docker-compose.yml
cat docker-compose.yml

# restart {{ $.dory.imageRepo.type }}
docker-compose up -d
sleep 10

# check {{ $.dory.imageRepo.type }} status
docker-compose ps
{{- end }}

{{- if $.imageRepoInternal }}
# copy {{ $.dory.imageRepo.type }} server ({{ $.imageRepoIp }}) certificates to all kubernetes nodes {{ $certPath }}/certs.d/{{ $.imageRepoDomainName }} directory
scp -r /etc/docker/certs.d root@${KUBERNETES_HOST}:{{ $certPath }}/
{{- else }}
# copy {{ $.dory.imageRepo.type }} server ({{ $.imageRepoIp }}) certificates to all kubernetes nodes {{ $certPath }}/certs.d/{{ $.imageRepoDomainName }} directory
# certificates are: ca.crt, {{ $.imageRepoDomainName }}.cert, {{ $.imageRepoDomainName }}.key
{{- end }}

# on current host and all kubernetes nodes add {{ $.dory.imageRepo.type }} domain name in /etc/hosts
vi /etc/hosts
{{ $.hostIP }}  {{ $.dory.imageRepo.internal.domainName }}

# docker login to {{ $.dory.imageRepo.type }}
docker login --username admin --password {{ $.dory.imageRepo.internal.password }} {{ $.dory.imageRepo.internal.domainName }}

# create public, hub, gcr, quay projects in {{ $.dory.imageRepo.type }}
curl -k -X POST -H 'Content-Type: application/json' -d '{"project_name": "public", "public": true}' 'https://admin:{{ $.dory.imageRepo.internal.password }}@{{ $.dory.imageRepo.internal.domainName }}/api/v2.0/projects'
curl -k -X POST -H 'Content-Type: application/json' -d '{"project_name": "hub", "public": true}' 'https://admin:{{ $.dory.imageRepo.internal.password }}@{{ $.dory.imageRepo.internal.domainName }}/api/v2.0/projects'
curl -k -X POST -H 'Content-Type: application/json' -d '{"project_name": "gcr", "public": true}' 'https://admin:{{ $.dory.imageRepo.internal.password }}@{{ $.dory.imageRepo.internal.domainName }}/api/v2.0/projects'
curl -k -X POST -H 'Content-Type: application/json' -d '{"project_name": "quay", "public": true}' 'https://admin:{{ $.dory.imageRepo.internal.password }}@{{ $.dory.imageRepo.internal.domainName }}/api/v2.0/projects'

# push container images to {{ $.dory.imageRepo.type }}
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
{{- end }}

## install dory services with docker-compose

```shell script
# create docker certificates
cd {{ $.rootDir }}/{{ $.dory.namespace }}/{{ $.dory.docker.dockerName }}
sh docker_certs.sh
ls -alh

# create openldap certificates
cd {{ $.rootDir }}/{{ $.dory.namespace }}/{{ $.dory.openldap.serviceName }}/certs
sh openldap_certs.sh
ls -alh
cp ca.crt ldap.crt ldap.key {{ $.rootDir }}/{{ $.dory.namespace }}/dory-engine/dory-data/certs/openldap

# create dory services directory and chown
cd {{ $.rootDir }}/{{ $.dory.namespace }}
mkdir -p mongo-dory
chown -R 999:999 mongo-dory
mkdir -p dory-engine/dory-data
mkdir -p dory-engine/tmp
chown -R 1000:1000 dory-engine
ls -alh

# start all dory services with docker-compose
cd {{ $.rootDir }}/{{ $.dory.namespace }}
ls -alh
docker-compose stop && docker-compose rm -f && docker-compose up -d

# check dory services status
sleep 10
docker-compose ps
```

## create project-data-pod pod in kubernetes

```shell script
# project-data-pod pod is used for create project directory in kuberentes
# create project-data-pod pod in kubernetes
cd {{ $.rootDir }}
kubectl apply -f project-data-pod.yaml
kubectl -n {{ $.dory.namespace }} get pods
```

## dory not config yet

2. please follow `README-1-docker-config.md` to config dory by manual after install
