# 以kubernetes方式部署完dory之后，必须进行以下设置

## 安装完成后必须进行dory-engine配置

{{- if $.dory.gitRepo.internal.image }}

### 完成 {{ $.dory.gitRepo.type }} 安装并更新dory的config.yaml配置

- url: {{ $.viewURL }}:{{ $.dory.gitRepo.internal.port }}
- 数据存放在以下目录: `{{ $.rootDir }}/{{ $.dory.namespace }}/{{ $.dory.gitRepo.type }}`

{{- if eq $.dory.gitRepo.type "gitea" }}
- 1. 打开gitea的网址，完成gitea安装设置，重点设置 `基础URL*` 和 `管理员账号` ，设置基础URL和管理员的用户名、密码、邮箱
- 2. 登录gitea，打开 `{{ $.viewURL }}:{{ $.dory.gitRepo.internal.port }}/user/settings/applications`，在`创建新Token`创建一个新的访问token。
{{- else if eq $.dory.gitRepo.type "gitlab" }}
- 1. gitlab的root用户密码文件存放在: `{{ $.rootDir }}/{{ $.dory.namespace }}/{{ $.dory.gitRepo.type }}/config/initial_root_password`
- 2. 登录gitlab，打开 `{{ $.viewURL }}:{{ $.dory.gitRepo.internal.port }}/-/profile/personal_access_tokens`，新增一个访问token。
- 3. 打开 `{{ $.viewURL }}:{{ $.dory.gitRepo.internal.port }}/admin/application_settings/network`，勾选以下选项
     Outbound requests -> Allow requests to the local network from web hooks and services -> 勾选，然后点击 "Save changes" 保存配置
{{- end }}
- 3. 记住管理员的 `用户名、密码、邮箱、访问token` 用于更新dory-engine的配置文件中的 {{ $.dory.gitRepo.type }} 设置
- 4. 更新dory-engine配置文件:
  - 配置文件存放在: `{{ $.rootDir }}/{{ $.dory.namespace }}/dory-engine/config/config.yaml`
  - 配置文件中搜索: `PLEASE_INPUT_BY_MANUAL`
  - 更新配置文件中以下代码仓库管理员设置: 
    - gitRepoConfigs.username
    - gitRepoConfigs.name
    - gitRepoConfigs.mail
    - gitRepoConfigs.password
    - gitRepoConfigs.token
{{- end }}

{{- if $.artifactRepoInternal }}

### 更新 {{ $.dory.artifactRepo.type }} 管理员密码，并更新dory的config.yaml配置文件

- url: {{ $.artifactRepoViewUrl }}
- user: admin / {{ $.artifactRepoPassword }} (管理员用户)
- 数据存放在: `{{ $.rootDir }}/{{ $.dory.namespace }}/{{ $.dory.artifactRepo.type }}`

- 1. 打开 {{ $.dory.artifactRepo.type }} 网址，使用admin的默认账号密码登录
- 2. 修改管理员密码: `{{ $.artifactRepoViewUrl }}/#user/account`
- 3. 更新dory-engine配置文件:
  - 配置文件存放在: `{{ $.rootDir }}/{{ $.dory.namespace }}/dory-engine/config/config.yaml`
  - 在配置文件中搜索 `{{ $.artifactRepoPassword }}`
  - 更新以下管理员密码配置: 
    - artifactRepoConfigs.password
{{- end }}

{{- if $.scanCodeRepoInternal }}

### 更新 {{ $.dory.scanCodeRepo.type }} 管理员密码，创建管理员token，并更新dory的config.yaml配置文件

- url: {{ $.scanCodeRepoViewUrl }}

- 1. {{ $.dory.scanCodeRepo.type }} 默认 admin 密码是 `admin`，首先需要更新默认管理员密码
- 2. 打开 `{{ $.scanCodeRepoViewUrl }}/account/security/`，创建管理员token，token类型选择 `User Token`, 过期时间选择 `No expiration`
- 3. 复制管理员的 `token` 用于更新dory-engine的配置文件中的 {{ $.dory.scanCodeRepo.type }} 设置
  - 配置文件存放在: `{{ $.rootDir }}/{{ $.dory.namespace }}/dory-engine/config/config.yaml`
  - 配置文件中搜索: `{{ $.scanCodeRepoToken }}`
  - 更新配置文件中以下代码扫描仓库管理员设置: 
    - scanCodeRepoConfigs.token
- 4. 出于安全考虑，设置项目默认查看权限为`私有`
  - 打开 `{{ $.scanCodeRepoViewUrl }}/admin/projects_management`，`Default visibility of new projects`设置为`Private`
{{- end }}

### trivy漏洞库更新

- 如果需要启用镜像扫描功能，请执行trivy漏洞库更新

```shell
docker run --rm -v {{ $.rootDir }}/{{ $.dory.namespace }}/dory-engine/dory-data/trivy:/root/.cache/trivy aquasec/trivy:0.37.2 image --download-db-only
docker run --rm -v {{ $.rootDir }}/{{ $.dory.namespace }}/dory-engine/dory-data/trivy:/root/.cache/trivy aquasec/trivy:0.37.2 image --download-java-db-only
chown -R 1000:1000 {{ $.rootDir }}/{{ $.dory.namespace }}/dory-engine/dory-data/trivy
```

### 重启 dory-engine 和 dory-console 服务

- 1. 重启 dory-engine 和 dory-console 服务

```shell script
kubectl -n {{ $.dory.namespace }} delete pods dory-engine-0 dory-console-0

# 等待 dory-engine-0 dory-console-0 pod处于ready状态
kubectl -n {{ $.dory.namespace }} get pods -o wide -w
```

## 访问各个dory服务

### dory-console 管理界面

- url: {{ $.viewURL }}:{{ $.dory.doryengine.port }}
- 管理员用户: {{ $.account.adminUser.username }}
- 管理员账号密码存放在: `{{ $.rootDir }}/{{ $.dory.namespace }}/dory-engine/dory-data/admin.password`
- dory-engine数据和配置存放在: `{{ $.rootDir }}/{{ $.dory.namespace }}/dory-engine`

{{- if $.dory.gitRepo.internal.image }}

### {{ $.dory.gitRepo.type }} 代码仓库

- url: {{ $.viewURL }}:{{ $.dory.gitRepo.internal.port }}
- 数据存放在: `{{ $.rootDir }}/{{ $.dory.namespace }}/{{ $.dory.gitRepo.type }}`
{{- end }}

{{- if $.artifactRepoInternal }}

### {{ $.dory.artifactRepo.type }} 依赖与制品仓库

- url: {{ $.artifactRepoViewUrl }}
- 公共用户账号: {{ $.artifactRepoPublicUser }} / {{ $.artifactRepoPublicPassword }}
- docker.io镜像代理地址: {{ $.artifactRepoIp }}:{{ $.artifactRepoPortHub }}
- gcr.io镜像代理地址: {{ $.artifactRepoIp }}:{{ $.artifactRepoPortGcr }}
- quay.io镜像代理地址: {{ $.artifactRepoIp }}:{{ $.artifactRepoPortQuay }}
{{- end }}

{{- if $.imageRepoInternal }}

### {{ $.dory.imageRepo.type }} 容器镜像仓库

- url: https://{{ $.imageRepoDomainName }}
- user: admin / {{ $.imageRepoPassword }} (管理员用户)
- 数据存放在: `{{ $.rootDir }}/{{ $.dory.imageRepo.internal.namespace }}`
{{- end }}

{{- if $.scanCodeRepoInternal }}

### {{ $.dory.scanCodeRepo.type }} 代码扫描仓库

- url: {{ $.scanCodeRepoViewUrl }}
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

### 注意，本目录非常重要，本目录为安装过程配置文件以及说明文件目录，建议保留
