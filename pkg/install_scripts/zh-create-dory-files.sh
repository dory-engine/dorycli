{{- $certPath := "" }}{{- if eq $.kubernetes.runtime "docker" }}{{ $certPath = "/etc/docker" }}{{- else if eq $.kubernetes.runtime "containerd" }}{{ $certPath = "/etc/containerd" }}{{- else if eq $.kubernetes.runtime "crio" }}{{ $certPath = "/etc/containers" }}{{- end }}
## 在kubernetes中创建project-data-pod pod
# project-data-pod pod 用于把创建的文件和目录目录复制到kubernetes
# 在kubernetes中创建project-data-pod
kubectl apply -f {{ $.dory.namespace }}/step01-namespace.yaml
kubectl apply -f project-data-pv.yaml
kubectl apply -f project-data-pod.yaml

# 等待project-data-pod成功启动
sh pods-ready.sh {{ $.dory.namespace }}

## 创建相关根目录
rm -rf install-data
rm -rf install-data.tar.gz

{{- if $.imageRepoInternal }}
# 创建 {{ $.dory.imageRepo.type }} 相关目录并设置目录权限
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

# 创建 openldap 的自签名证书
mkdir -p install-data/{{ $.dory.namespace }}/dory-engine/dory-data/certs/openldap
cd {{ $.dory.namespace }}/{{ $.dory.openldap.serviceName }}/certs
sh openldap_certs.sh
ls -alh
cd ../../..
cp {{ $.dory.namespace }}/{{ $.dory.openldap.serviceName }}/certs/ca.crt install-data/{{ $.dory.namespace }}/dory-engine/dory-data/certs/openldap
cp {{ $.dory.namespace }}/{{ $.dory.openldap.serviceName }}/certs/ldap.crt install-data/{{ $.dory.namespace }}/dory-engine/dory-data/certs/openldap
cp {{ $.dory.namespace }}/{{ $.dory.openldap.serviceName }}/certs/ldap.key install-data/{{ $.dory.namespace }}/dory-engine/dory-data/certs/openldap

# 创建 dory 组件相关目录并设置目录权限
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

## 创建dory组件相关文件

# 创建 docker executor 自签名证书
sh {{ $.dory.namespace }}/{{ $.dory.docker.dockerName }}/docker_certs.sh

# 创建 docker executor 自签名证书
kubectl -n {{ $.dory.namespace }} create secret generic {{ $.dory.docker.dockerName }}-tls --from-file=docker-certs/ca.crt --from-file=docker-certs/tls.crt --from-file=docker-certs/tls.key --dry-run=client -o yaml | kubectl apply -f -
kubectl -n {{ $.dory.namespace }} describe secret {{ $.dory.docker.dockerName }}-tls

# 把install-data目录所有文件通过project-data-pod复制到kubernetes的共享存储中
tar -czvf install-data.tar.gz -C install-data .
kubectl -n {{ $.dory.namespace }} cp install-data.tar.gz project-data-pod-0:/project-data/
kubectl -n {{ $.dory.namespace }} exec project-data-pod-0 -- tar zxvf /project-data/install-data.tar.gz -C /project-data/
kubectl -n {{ $.dory.namespace }} exec project-data-pod-0 -- tree -u -g -L 3 /project-data/
kubectl -n {{ $.dory.namespace }} exec project-data-pod-0 -- rm -f /project-data/install-data.tar.gz

{{- if $.imageRepoInternal }}
# 创建 {{ $.dory.imageRepo.type }} 名字空间与pv
kubectl apply -f {{ $.dory.imageRepo.internal.namespace }}/step01-namespace.yaml
kubectl apply -f {{ $.dory.imageRepo.internal.namespace }}/step01-pv.yaml

# 使用helm安装 {{ $.dory.imageRepo.type }}
helm install -n {{ $.dory.imageRepo.internal.namespace }} {{ $.dory.imageRepo.internal.namespace }} {{ $.dory.imageRepo.type }}
helm -n {{ $.dory.imageRepo.internal.namespace }} list

# 等待所有 {{ $.dory.imageRepo.type }} 服务状态为 ready
sh pods-ready.sh {{ $.dory.imageRepo.internal.namespace }}

# 创建 {{ $.dory.imageRepo.type }} 自签名证书并复制到 {{ $certPath }}/certs.d
sh {{ $.dory.imageRepo.internal.namespace }}/harbor_update_docker_certs.sh
ls -alh {{ $certPath }}/certs.d/{{ $.imageRepoDomainName }}

# 复制容器镜像仓库证书到 docker executor 配置目录
cp -rp {{ $certPath }}/certs.d/{{ $.imageRepoDomainName }} install-data/{{ $.dory.namespace }}/{{ $.dory.docker.dockerName }}
kubectl -n {{ $.dory.namespace }} cp {{ $certPath }}/certs.d/{{ $.imageRepoDomainName }} project-data-pod-0:/project-data/{{ $.dory.namespace }}/{{ $.dory.docker.dockerName }}
{{- end }}
