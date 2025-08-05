# 如何访问dory

## 如何访问dory以及相关组件的文件

- dory以及相关组件的文件位于kubernetes中部署的 `project-data-pod` 容器中的 `/project-data/` 目录

```shell script
# 执行以下命令进入project-data-pod，可以查看dory的配置文件
kubectl -n {{ $.dory.namespace }} exec -ti project-data-pod-0 -- ash
cd /project-data
```

## trivy镜像扫描漏洞库更新

- 如果需要启用镜像扫描功能，请执行trivy漏洞库更新

```shell script
# 下载trivy漏洞库
{{ $.cmdRun }} --rm -v trivy:/root/.cache/trivy aquasec/trivy:0.37.2 image --download-db-only
{{ $.cmdRun }} --rm -v trivy:/root/.cache/trivy aquasec/trivy:0.37.2 image --download-java-db-only

# 把trivy漏洞库上传到kubernetes共享存储
kubectl -n {{ $.dory.namespace }} cp trivy project-data-pod-0:/project-data/{{ $.dory.namespace }}/dory-engine/dory-data/
kubectl -n {{ $.dory.namespace }} exec project-data-pod-0 -- chown -R 1000:1000 /project-data/{{ $.dory.namespace }}/dory-engine/dory-data/trivy
```

## 访问各个dory服务

### dory-engine 管理界面

- url: {{ $.viewURL }}:{{ $.dory.doryengine.port }}
- 管理员用户: {{ $.account.adminUser.username }}
- 管理员账号密码存放在: `/project-data/{{ $.dory.namespace }}/dory-engine/dory-data/admin.password`
- dory-engine数据和配置存放在: `/project-data/{{ $.dory.namespace }}/dory-engine`
- dory-engine配置文件存放在: `/project-data/{{ $.dory.namespace }}/dory-engine/config/config.yaml`

{{- if $.imageRepoInternal }}

### {{ $.dory.imageRepo.type }} 容器镜像仓库

- url: https://{{ $.imageRepoDomainName }}
- 管理员账号: admin / {{ $.imageRepoPassword }}
- 数据存放在: `/project-data/{{ $.dory.imageRepo.internal.namespace }}`
{{- end }}

{{- if $.gitRepoInternal }}

### {{ $.dory.gitRepo.type }} 代码仓库

- url: {{ $.viewURL }}:{{ $.dory.gitRepo.internal.port }}
- 管理员账号: root / {{ $.gitRepoPassword }}
- 数据存放在: `/project-data/{{ $.dory.namespace }}/{{ $.dory.gitRepo.type }}`
{{- end }}

{{- if $.artifactRepoInternal }}

### {{ $.dory.artifactRepo.type }} 依赖与制品仓库

- url: {{ $.artifactRepoViewUrl }}
- 管理员账号: admin / {{ $.artifactRepoPassword }}
- 公共用户账号: {{ $.artifactRepoPublicUser }} / {{ $.artifactRepoPublicPassword }}
- docker.io镜像代理地址: {{ $.artifactRepoIp }}:{{ $.artifactRepoPortHub }}
- gcr.io镜像代理地址: {{ $.artifactRepoIp }}:{{ $.artifactRepoPortGcr }}
- quay.io镜像代理地址: {{ $.artifactRepoIp }}:{{ $.artifactRepoPortQuay }}

- 数据存放在: `/project-data/{{ $.dory.namespace }}/{{ $.dory.artifactRepo.type }}`
{{- end }}

{{- if $.scanCodeRepoInternal }}

### {{ $.dory.scanCodeRepo.type }} 代码扫描仓库

- url: {{ $.scanCodeRepoViewUrl }}
- 管理员账号: admin / {{ $.scanCodeRepoPassword }}
{{- end }}

### openldap 账号管理中心

- url: {{ $.viewURL | replace "http://" "https://" }}:{{ $.dory.openldap.port }}
- 管理员用户: cn=admin,{{ $.dory.openldap.baseDN }} / {{ $.dory.openldap.password }}

{{- if $.demoDatabaseInternal }}

### 项目演示数据库

- jdbc 连接 url: {{ $.demoDatabaseUrl }}
- 用户: {{ $.demoDatabaseUsername }} / {{ $.demoDatabasePassword }}
{{- end }}

{{- if $.demoHostInternal }}

### 项目演示ssh主机

- ssh 命令: `ssh -p {{ $.demoHostPort }} root@{{ $.demoHostAddr }}`
- 密码: {{ $.demoHostPassword }}
- 演示ssh主机暴露的web服务 url:  http://{{ $.demoHostAddr }}:{{ $.demoHostNodePortWeb }}
{{- end }}

### 注意，本目录非常重要，建议保留
