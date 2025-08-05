{{- $certPath := "" }}{{- if eq $.kubernetes.runtime "docker" }}{{ $certPath = "/etc/docker" }}{{- else if eq $.kubernetes.runtime "containerd" }}{{ $certPath = "/etc/containerd" }}{{- else if eq $.kubernetes.runtime "crio" }}{{ $certPath = "/etc/containers" }}{{- end }}
## create project-data-pod pod in kubernetes
# project-data-pod pod is used for copy project directory and files in kuberentes
# create project-data-pod pod in kubernetes
kubectl apply -f {{ $.dory.namespace }}/step01-namespace.yaml
kubectl apply -f project-data-pv.yaml
kubectl apply -f project-data-pod.yaml

# waiting project-data-pod ready
sh pods-ready.sh {{ $.dory.namespace }}

## create install root directories
rm -rf install-data
rm -rf install-data.tar.gz

{{- if $.imageRepoInternal }}
# create {{ $.dory.imageRepo.type }} root directory
mkdir -p install-data/{{ $.dory.imageRepo.internal.namespace }}/database
mkdir -p install-data/{{ $.dory.imageRepo.internal.namespace }}/jobservice
mkdir -p install-data/{{ $.dory.imageRepo.internal.namespace }}/redis
mkdir -p install-data/{{ $.dory.imageRepo.internal.namespace }}/registry
chown -R 999:999 install-data/{{ $.dory.imageRepo.internal.namespace }}/database
chown -R 10000:10000 install-data/{{ $.dory.imageRepo.internal.namespace }}/jobservice
chown -R 999:999 install-data/{{ $.dory.imageRepo.internal.namespace }}/redis
chown -R 10000:10000 install-data/{{ $.dory.imageRepo.internal.namespace }}/registry
find install-data/{{ $.dory.imageRepo.internal.namespace }} -type d -exec chmod a+rx {} \;
find install-data/{{ $.dory.imageRepo.internal.namespace }} -type f -exec chmod a+r {} \;
ls -alh install-data/{{ $.dory.imageRepo.internal.namespace }}
{{- end }}

# create openldap certificates
mkdir -p install-data/{{ $.dory.namespace }}/dory-engine/dory-data/certs/openldap
cd {{ $.dory.namespace }}/{{ $.dory.openldap.serviceName }}/certs
sh openldap_certs.sh
ls -alh
cd ../../..
cp {{ $.dory.namespace }}/{{ $.dory.openldap.serviceName }}/certs/ca.crt install-data/{{ $.dory.namespace }}/dory-engine/dory-data/certs/openldap
cp {{ $.dory.namespace }}/{{ $.dory.openldap.serviceName }}/certs/ldap.crt install-data/{{ $.dory.namespace }}/dory-engine/dory-data/certs/openldap
cp {{ $.dory.namespace }}/{{ $.dory.openldap.serviceName }}/certs/ldap.key install-data/{{ $.dory.namespace }}/dory-engine/dory-data/certs/openldap

# create dory root directory
mkdir -p install-data/{{ $.dory.namespace }}/dory-engine/logs
cp -rp {{ $.dory.namespace }}/{{ $.dory.docker.dockerName }} install-data/{{ $.dory.namespace }}/
cp -rp {{ $.dory.namespace }}/{{ $.dory.openldap.serviceName }} install-data/{{ $.dory.namespace }}/
cp -rp {{ $.dory.namespace }}/dory-engine install-data/{{ $.dory.namespace }}/
{{- if and (eq $.dory.gitRepo.type "gitlab") $.dory.gitRepo.internal.image }}
cp -rp {{ $.dory.namespace }}/nginx-gitlab install-data/{{ $.dory.namespace }}/
{{- end }}
{{- if $.scanCodeRepoInternal }}
mkdir -p install-data/{{ $.dory.namespace }}/sonarqube-web/data
mkdir -p install-data/{{ $.dory.namespace }}/sonarqube-web/extensions
mkdir -p install-data/{{ $.dory.namespace }}/sonarqube-web/logs
mkdir -p install-data/{{ $.dory.namespace }}/sonarqube-web/temp
chown -R 1000:1000 install-data/{{ $.dory.namespace }}/sonarqube-web
{{- end }}
{{ if $.artifactRepoInternal }}
mkdir -p install-data/{{ $.dory.namespace }}/nexus
chown -R 200:200 install-data/{{ $.dory.namespace }}/nexus
{{- end }}
mkdir -p install-data/{{ $.dory.namespace }}/mongo-dory
chown -R 999:999 install-data/{{ $.dory.namespace }}/mongo-dory
ls -alh install-data/{{ $.dory.namespace }}
chown -R 1000:1000 install-data/{{ $.dory.namespace }}/dory-engine
find install-data/{{ $.dory.namespace }} -type d -exec chmod a+rx {} \;
find install-data/{{ $.dory.namespace }} -type f -exec chmod a+r {} \;

## create dory components files

# create docker executor self signed certificates
sh {{ $.dory.namespace }}/{{ $.dory.docker.dockerName }}/docker_certs.sh

# create docker executor self signed certificates
kubectl -n {{ $.dory.namespace }} create secret generic {{ $.dory.docker.dockerName }}-tls --from-file=docker-certs/ca.crt --from-file=docker-certs/tls.crt --from-file=docker-certs/tls.key --dry-run=client -o yaml | kubectl apply -f -
kubectl -n {{ $.dory.namespace }} describe secret {{ $.dory.docker.dockerName }}-tls

# copy install-data all files to kubernetes shared storage by pod project-data-pod
tar -czvf install-data.tar.gz -C install-data .
kubectl -n {{ $.dory.namespace }} cp install-data.tar.gz project-data-pod-0:/project-data/
kubectl -n {{ $.dory.namespace }} exec project-data-pod-0 -- tar zxvf /project-data/install-data.tar.gz -C /project-data/
kubectl -n {{ $.dory.namespace }} exec project-data-pod-0 -- tree -u -g -L 3 /project-data/
kubectl -n {{ $.dory.namespace }} exec project-data-pod-0 -- rm -f /project-data/install-data.tar.gz

{{- if $.imageRepoInternal }}
# create {{ $.dory.imageRepo.type }} namespace and pv
kubectl apply -f {{ $.dory.imageRepo.internal.namespace }}/step01-namespace.yaml
kubectl apply -f {{ $.dory.imageRepo.internal.namespace }}/step01-pv.yaml

# install {{ $.dory.imageRepo.type }}
helm install -n {{ $.dory.imageRepo.internal.namespace }} {{ $.dory.imageRepo.internal.namespace }} {{ $.dory.imageRepo.type }}
helm -n {{ $.dory.imageRepo.internal.namespace }} list

# waiting for all {{ $.dory.imageRepo.type }} services ready
sh pods-ready.sh {{ $.dory.imageRepo.internal.namespace }}

# create {{ $.dory.imageRepo.type }} self signed certificates and copy to {{ $certPath }}/certs.d
sh {{ $.dory.imageRepo.internal.namespace }}/harbor_update_docker_certs.sh
ls -alh {{ $certPath }}/certs.d/{{ $.imageRepoDomainName }}

# copy container image repository certificates in docker executor settings directory
cp -rp {{ $certPath }}/certs.d/{{ $.imageRepoDomainName }} install-data/{{ $.dory.namespace }}/{{ $.dory.docker.dockerName }}
kubectl -n {{ $.dory.namespace }} cp {{ $certPath }}/certs.d/{{ $.imageRepoDomainName }} project-data-pod-0:/project-data/{{ $.dory.namespace }}/{{ $.dory.docker.dockerName }}
{{- end }}
