# install dory with kubernetes

## summary

1. please follow `README-0-kubernetes-install.md` to install dory by manual
2. please follow `README-1-kubernetes-config.md` to config dory by manual after install
3. if install fail, please follow `README-2-kubernetes-reset.md` to stop all dory services and install again

## create install root directories

```shell script
{{- if $.imageRepoInternal }}
# create {{ $.dory.imageRepo.type }} root directory
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

# create openldap certificates
mkdir -p {{ $.rootDir }}/{{ $.dory.namespace }}/dory-engine/dory-data/certs/openldap
cd {{ $.dory.namespace }}/{{ $.dory.openldap.serviceName }}/certs
sh openldap_certs.sh
ls -alh
cp ca.crt ldap.crt ldap.key {{ $.rootDir }}/{{ $.dory.namespace }}/dory-engine/dory-data/certs/openldap
cd ../../.. 

# create dory root directory
mkdir -p {{ $.rootDir }}/{{ $.dory.namespace }}/dory-engine/tmp
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
cp -rp /usr/share/zoneinfo {{ $.rootDir }}/{{ $.dory.namespace }}/dory-engine/dory-data
find {{ $.rootDir }}/{{ $.dory.namespace }}/dory-engine/dory-data/zoneinfo -type f -exec chmod a+r {} \;
find {{ $.rootDir }}/{{ $.dory.namespace }}/dory-engine/dory-data/zoneinfo -type d -exec chmod a+rx {} \;
mkdir -p {{ $.rootDir }}/timezone
echo '{{ $.kubernetes.timezone }}' > {{ $.rootDir }}/timezone/timezone
cp -rp /usr/share/zoneinfo {{ $.rootDir }}/timezone
mkdir -p {{ $.rootDir }}/{{ $.dory.namespace }}/mongo-dory
chown -R 999:999 {{ $.rootDir }}/{{ $.dory.namespace }}/mongo-dory
ls -alh {{ $.rootDir }}/{{ $.dory.namespace }}
chown -R 1000:1000 {{ $.rootDir }}/{{ $.dory.namespace }}/dory-engine
find {{ $.rootDir }}/timezone -type f -exec chmod a+r {} \;
find {{ $.rootDir }}/timezone -type d -exec chmod a+rx {} \;
```

{{ $certPath := "" }}{{- if eq $.kubernetes.runtime "docker" }}{{ $certPath = "/etc/docker" }}{{- else if eq $.kubernetes.runtime "containerd" }}{{ $certPath = "/etc/containerd" }}{{- else if eq $.kubernetes.runtime "crio" }}{{ $certPath = "/etc/containers" }}{{- end }}
{{- if eq $.dory.imageRepo.type "harbor" }}
## {{ $.dory.imageRepo.type }} installation and configuration

```shell script
{{- if $.imageRepoInternal }}
# create {{ $.dory.imageRepo.type }} namespace and pv
kubectl delete ns {{ $.dory.imageRepo.internal.namespace }}
kubectl delete pv {{ $.dory.imageRepo.internal.namespace }}-pv
kubectl delete pv {{ $.dory.imageRepo.internal.namespace }}-timezone-pv
kubectl apply -f {{ $.dory.imageRepo.internal.namespace }}/step01-namespace-pv.yaml

# install {{ $.dory.imageRepo.type }}
helm install -n {{ $.dory.imageRepo.internal.namespace }} {{ $.dory.imageRepo.internal.namespace }} {{ $.dory.imageRepo.type }}
helm -n {{ $.dory.imageRepo.internal.namespace }} list

# waiting for all {{ $.dory.imageRepo.type }} services ready
kubectl -n {{ $.dory.imageRepo.internal.namespace }} get pods -o wide

# create {{ $.dory.imageRepo.type }} self signed certificates and copy to {{ $certPath }}/certs.d
sh {{ $.dory.imageRepo.internal.namespace }}/harbor_update_docker_certs.sh
ls -alh {{ $certPath }}/certs.d/{{ $.dory.imageRepo.internal.domainName }}
{{- end }}

{{- if $.imageRepoInternal }}
# copy {{ $.dory.imageRepo.type }} server ({{ $.imageRepoIp }}) certificates to all kubernetes nodes {{ $certPath }}/certs.d/{{ $.imageRepoDomainName }} directory
scp -r {{ $certPath }}/certs.d root@${KUBERNETES_HOST}:{{ $certPath }}/
{{- else }}
# copy {{ $.dory.imageRepo.type }} server ({{ $.imageRepoIp }}) certificates to all kubernetes nodes {{ $certPath }}/certs.d/{{ $.imageRepoDomainName }} directory
# certificates are: ca.crt, {{ $.imageRepoDomainName }}.cert, {{ $.imageRepoDomainName }}.key
{{- end }}

# on current host and all kubernetes nodes add {{ $.dory.imageRepo.type }} domain name in /etc/hosts
vi /etc/hosts
{{ $.imageRepoIp }}  {{ $.imageRepoDomainName }}

# {{ $.kubernetes.runtime }} client login to {{ $.dory.imageRepo.type }}
{{ $.cmdLogin }} --username {{ $.imageRepoUsername }} --password {{ $.imageRepoPassword }} {{ $.imageRepoDomainName }}

# create public, hub, gcr, quay projects in {{ $.dory.imageRepo.type }}
curl -k -X POST -H 'Content-Type: application/json' -d '{"project_name": "public", "public": true}' 'https://{{ $.imageRepoUsername }}:{{ $.imageRepoPassword }}@{{ $.imageRepoDomainName }}/api/v2.0/projects'
curl -k -X POST -H 'Content-Type: application/json' -d '{"project_name": "hub", "public": true}' 'https://{{ $.imageRepoUsername }}:{{ $.imageRepoPassword }}@{{ $.imageRepoDomainName }}/api/v2.0/projects'
curl -k -X POST -H 'Content-Type: application/json' -d '{"project_name": "gcr", "public": true}' 'https://{{ $.imageRepoUsername }}:{{ $.imageRepoPassword }}@{{ $.imageRepoDomainName }}/api/v2.0/projects'
curl -k -X POST -H 'Content-Type: application/json' -d '{"project_name": "quay", "public": true}' 'https://{{ $.imageRepoUsername }}:{{ $.imageRepoPassword }}@{{ $.imageRepoDomainName }}/api/v2.0/projects'

# push container images to {{ $.dory.imageRepo.type }}
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

## install dory services with kubernetes

```shell script
# create {{ $.dory.namespace }} namespace and pv
kubectl delete ns {{ $.dory.namespace }}
kubectl delete pv {{ $.dory.namespace }}-pv
kubectl delete pv {{ $.dory.namespace }}-timezone-pv
kubectl apply -f {{ $.dory.namespace }}/step01-namespace-pv.yaml

# create docker executor self signed certificates
sh {{ $.dory.namespace }}/{{ $.dory.docker.dockerName }}/docker_certs.sh
kubectl -n {{ $.dory.namespace }} create secret generic {{ $.dory.docker.dockerName }}-tls --from-file=certs/ca.crt --from-file=certs/tls.crt --from-file=certs/tls.key --dry-run=client -o yaml | kubectl apply -f -
kubectl -n {{ $.dory.namespace }} describe secret {{ $.dory.docker.dockerName }}-tls
rm -rf certs

# copy container image repository certificates in docker executor settings directory
cp -rp {{ $certPath }}/certs.d/{{ $.imageRepoDomainName }} {{ $.rootDir }}/{{ $.dory.namespace }}/{{ $.dory.docker.dockerName }}

{{- if $.artifactRepoInternal }}
# extract nexus init data
tar Cxzvf {{ $.rootDir }}/{{ $.dory.namespace }} ../{{ $.dory.nexusInitData }}
chown -R 200:200 {{ $.rootDir }}/{{ $.dory.namespace }}/nexus
ls -alh {{ $.rootDir }}/{{ $.dory.namespace }}/nexus
{{- end }}

# start all dory services with kubernetes
kubectl apply -f {{ $.dory.namespace }}/step02-statefulset.yaml
kubectl apply -f {{ $.dory.namespace }}/step03-service.yaml
kubectl apply -f {{ $.dory.namespace }}/step04-networkpolicy.yaml

# check dory services status
kubectl -n {{ $.dory.namespace }} get pods -o wide
```

## create project-data-pod pod in kubernetes

```shell script
# project-data-pod pod is used for create project directory in kuberentes
# create project-data-pod pod in kubernetes
kubectl apply -f project-data-pod.yaml
kubectl -n {{ $.dory.namespace }} get pods
```

## dory not config yet

2. please follow `README-1-kubernetes-config.md` to config dory by manual after install
