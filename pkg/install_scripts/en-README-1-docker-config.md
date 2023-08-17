# config dory after install in docker

## dory-engine settings after installed

### create directory in kubernetes shared storage

- create directory

{{- if $.kubernetes.pvConfigLocal.localPath }}
```shell script
# create directory in kubernetes local storage
mkdir -p {{ $.kubernetes.pvConfigLocal.localPath }}/timezone
rm -rf {{ $.kubernetes.pvConfigLocal.localPath }}/timezone/*
# copy timezone files to this folder
echo '{{ $.kubernetes.timezone }}' > {{ $.kubernetes.pvConfigLocal.localPath }}/timezone/timezone
cp -rp /usr/share/zoneinfo {{ $.kubernetes.pvConfigLocal.localPath }}/timezone
```
{{- else if $.kubernetes.pvConfigNfs.nfsPath }}
```shell script
# create directory in kubernetes nfs storage
mkdir -p {{ $.kubernetes.pvConfigNfs.nfsPath }}/timezone
rm -rf {{ $.kubernetes.pvConfigNfs.nfsPath }}/timezone/*
# copy timezone files to this folder
echo '{{ $.kubernetes.timezone }}' > {{ $.kubernetes.pvConfigNfs.nfsPath }}/timezone/timezone
cp -rp /usr/share/zoneinfo {{ $.kubernetes.pvConfigNfs.nfsPath }}/timezone
```
{{- else if $.kubernetes.pvConfigCephfs.cephPath }}
```shell script
# create directory in kubernetes cephfs storage
mkdir -p {{ $.kubernetes.pvConfigCephfs.cephPath }}/timezone
rm -rf {{ $.kubernetes.pvConfigCephfs.cephPath }}/timezone/*
# copy timezone files to this folder
echo '{{ $.kubernetes.timezone }}' > {{ $.kubernetes.pvConfigCephfs.cephPath }}/timezone/timezone
cp -rp /usr/share/zoneinfo {{ $.kubernetes.pvConfigCephfs.cephPath }}/timezone
```
{{- end }}

- restart project-data-pod-0 pods

```shell script
kubectl -n {{ $.dory.namespace }} delete pods project-data-pod-0

# check project-data-pod-0 pod status is ready
kubectl -n {{ $.dory.namespace }} get pods project-data-pod-0
```

{{- if $.dory.gitRepo.internal.image }}
### finish {{ $.dory.gitRepo.type }} install and update dory config.yaml

- url: {{ $.viewURL }}:{{ $.dory.gitRepo.internal.port }}
- data located at: `{{ $.rootDir }}/{{ $.dory.namespace }}/{{ $.dory.gitRepo.type }}`

{{- if eq $.dory.gitRepo.type "gitea" }}
- 1. open gitea url finish gitea install, at `Gitea Base URL*` and `Administrator Account Settings` set gitea base url and admin username / password / email
- 2. login to gitea, open `{{ $.viewURL }}:{{ $.dory.gitRepo.internal.port }}/user/settings/applications`, at `Generate New Token` generate a new token.
{{- else if eq $.dory.gitRepo.type "gitlab" }}
- 1. gitlab password file located at: `{{ $.rootDir }}/{{ $.dory.namespace }}/{{ $.dory.gitRepo.type }}/config/initial_root_password`
- 2. login to gitlab, open `{{ $.viewURL }}:{{ $.dory.gitRepo.internal.port }}/-/profile/personal_access_tokens`, add a personal access token.
- 3. open `{{ $.viewURL }}:{{ $.dory.gitRepo.internal.port }}/admin/application_settings/network`, find following options and check it
     Outbound requests -> Allow requests to the local network from web hooks and services -> check it, and "Save changes"
{{- end }}
- 3. copy admin `username / password / email / token` to update dory-engine config file {{ $.dory.gitRepo.type }} settings
- 4. update dory-engine config file:
  - config file located at: `{{ $.rootDir }}/{{ $.dory.namespace }}/dory-engine/config/config.yaml`
  - search `PLEASE_INPUT_BY_MANUAL` in config file
  - update following admin user settings: 
    - gitRepoConfigs.username
    - gitRepoConfigs.name
    - gitRepoConfigs.mail
    - gitRepoConfigs.password
    - gitRepoConfigs.token
{{- end }}
    
{{- if $.artifactRepoInternal }}
### update {{ $.dory.artifactRepo.type }} admin password and update dory config.yaml

- url: {{ $.artifactRepoViewUrl }}
- user: admin / {{ $.artifactRepoPassword }} (admin user)
- data located at: `{{ $.rootDir }}/{{ $.dory.namespace }}/{{ $.dory.artifactRepo.type }}`

- 1. open {{ $.dory.artifactRepo.type }} url, login as admin user
- 2. change admin password, open `{{ $.artifactRepoViewUrl }}/#user/account` and change password
- 3. update dory-engine config file:
  - config file located at: `{{ $.rootDir }}/{{ $.dory.namespace }}/dory-engine/config/config.yaml`
  - search `{{ $.artifactRepoPassword }}` in config file
  - update following admin user password: 
    - artifactRepoConfigs.password
{{- end }}

{{- if $.scanCodeRepoInternal }}
### update {{ $.dory.scanCodeRepo.type }} admin password, create admin token and update dory config.yaml

- url: {{ $.scanCodeRepoViewUrl }}

- 1. {{ $.dory.scanCodeRepo.type }} default admin password is `admin`, update admin password first
- 2. open `{{ $.scanCodeRepoViewUrl }}/account/security/`, at `Generate Tokens` generate a new token, select token type `User Token`, select expires in `No expiration`
- 3. copy admin `token` update dory-engine config file:
  - config file located at: `{{ $.rootDir }}/{{ $.dory.namespace }}/dory-engine/config/config.yaml`
  - search `{{ $.scanCodeRepoToken }}` in config file
  - update following admin user settings: 
    - scanCodeRepoConfigs.token
- 4. for security reason, set project default visibility to `Priviate`
  - open `{{ $.scanCodeRepoViewUrl }}/admin/projects_management`, `Default visibility of new projects` set as `Private`
{{- end }}

### restart dory-engine and dory-console

- 1. restart dory-engine and dory-console

```shell script
cd {{ $.rootDir }}/{{ $.dory.namespace }}
docker stop dory-engine dory-console
docker start dory-engine dory-console
```

## connect your dory

### dory-console admin dashboard

- url: {{ $.viewURL }}:{{ $.dory.doryengine.port }}
- user: {{ $.account.adminUser.username }}
- password file located at: `{{ $.rootDir }}/{{ $.dory.namespace }}/dory-engine/dory-data/admin.password`
- data located at: `{{ $.rootDir }}/{{ $.dory.namespace }}/dory-engine`

{{- if $.dory.gitRepo.internal.image }}
### {{ $.dory.gitRepo.type }} git repository

- url: {{ $.viewURL }}:{{ $.dory.gitRepo.internal.port }}
- data located at: `{{ $.rootDir }}/{{ $.dory.namespace }}/{{ $.dory.gitRepo.type }}`
{{- end }}

{{- if $.artifactRepoInternal }}
### {{ $.dory.artifactRepo.type }} artifact and dependency repository

- url: {{ $.artifactRepoViewUrl }}
- public user: {{ $.artifactRepoPublicUser }} / {{ $.artifactRepoPublicPassword }}
- docker.io image proxy: {{ $.artifactRepoIp }}:{{ $.artifactRepoPortHub }}
- gcr.io image proxy: {{ $.artifactRepoIp }}:{{ $.artifactRepoPortGcr }}
- quay.io image proxy: {{ $.artifactRepoIp }}:{{ $.artifactRepoPortQuay }}
{{- end }}

{{- if $.imageRepoInternal }}
### {{ $.dory.imageRepo.type }} image repository

- url: https://{{ $.imageRepoDomainName }}
- user: admin / {{ $.imageRepoPassword }} (admin user)
- data located at: `{{ $.rootDir }}/{{ $.dory.imageRepo.internal.namespace }}`
{{- end }}

### openldap account management

- url: {{ $.viewURL | replace "http://" "https://" }}:{{ $.dory.openldap.port }}
- user: cn=admin,{{ $.dory.openldap.baseDN }} / {{ $.dory.openldap.password }}

{{- if $.scanCodeRepoInternal }}
### {{ $.dory.scanCodeRepo.type }} scan code repository

- url: {{ $.scanCodeRepoViewUrl }}
{{- end }}

### trivy vulnerabilities database update

```shell
docker run --rm -v {{ $.rootDir }}/{{ $.dory.namespace }}/dory-engine/dory-data/trivy:/root/.cache/trivy aquasec/trivy:0.37.2 image --download-db-only
docker run --rm -v {{ $.rootDir }}/{{ $.dory.namespace }}/dory-engine/dory-data/trivy:/root/.cache/trivy aquasec/trivy:0.37.2 image --download-java-db-only
chown -R 1000:1000 {{ $.rootDir }}/{{ $.dory.namespace }}/dory-engine/dory-data/trivy
```

{{- if $.demoDatabaseInternal }}
### project demo database

- connect jdbc url: {{ $.demoDatabaseUrl }}
- user: {{ $.demoDatabaseUsername }} / {{ $.demoDatabasePassword }}
{{- end }}

{{- if $.demoHostInternal }}
### project demo ssh host

- ssh command: `ssh -p {{ $.demoHostPort }} root@{{ $.demoHostAddr }}`
- password: {{ $.demoHostPassword }}
- demo ssh host expose web service url:  http://{{ $.demoHostAddr }}:{{ $.demoHostNodePortWeb }}
{{- end }}

### caution: this folder is very important, included all config files and readme files, please keep it
