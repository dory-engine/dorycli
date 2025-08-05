# how to connect dory

## how to access dory and related components files

- the files of dory and related components are located in the `/project-data/` directory of the `project-data-pod` container deployed in Kubernetes

```shell script
# execute the following command to enter project-data-pod and check the configuration file of dory
kubectl -n {{ $.dory.namespace }} exec -ti project-data-pod-0 -- ash
cd /project-data
```

### trivy scan image vulnerabilities database update

- If you need to enable the image scanning function, please perform trivy vulnerability library update

```shell
# download trivy vulnerabilities database
{{ $.cmdRun }} --rm -v trivy:/root/.cache/trivy aquasec/trivy:0.37.2 image --download-db-only
{{ $.cmdRun }} --rm -v trivy:/root/.cache/trivy aquasec/trivy:0.37.2 image --download-java-db-only

# copy trivy vulnerabilities database to kubernetes shared storage
kubectl -n {{ $.dory.namespace }} cp trivy project-data-pod-0:/project-data/{{ $.dory.namespace }}/dory-engine/dory-data/
kubectl -n {{ $.dory.namespace }} exec project-data-pod-0 -- chown -R 1000:1000 /project-data/{{ $.dory.namespace }}/dory-engine/dory-data/trivy
```

## connect your dory

### dory-engine dashboard

- url: {{ $.viewURL }}:{{ $.dory.doryengine.port }}
- user: {{ $.account.adminUser.username }}
- password file located at: `/project-data/{{ $.dory.namespace }}/dory-engine/dory-data/admin.password`
- dory-engine data located at: `/project-data/{{ $.dory.namespace }}/dory-engine`
- dory-engine config file located at: `/project-data/{{ $.dory.namespace }}/dory-engine/config/config.yaml`

{{- if $.imageRepoInternal }}

### {{ $.dory.imageRepo.type }} image repository

- url: https://{{ $.imageRepoDomainName }}
- admin user: admin / {{ $.imageRepoPassword }}
- data located at: `/project-data/{{ $.dory.imageRepo.internal.namespace }}`
{{- end }}

{{- if $.dory.gitRepo.internal.image }}

### {{ $.dory.gitRepo.type }} git repository

- url: {{ $.viewURL }}:{{ $.dory.gitRepo.internal.port }}
- admin user: root / {{ $.gitRepoPassword }}
- data located at: `/project-data/{{ $.dory.namespace }}/{{ $.dory.gitRepo.type }}`
{{- end }}

{{- if $.artifactRepoInternal }}

### {{ $.dory.artifactRepo.type }} artifact and dependency repository

- url: {{ $.artifactRepoViewUrl }}
- admin user: admin / {{ $.artifactRepoPassword }}
- public user: {{ $.artifactRepoPublicUser }} / {{ $.artifactRepoPublicPassword }}
- docker.io registry proxy: {{ $.artifactRepoIp }}:{{ $.artifactRepoPortHub }}
- gcr.io registry proxy: {{ $.artifactRepoIp }}:{{ $.artifactRepoPortGcr }}
- quay.io registry proxy: {{ $.artifactRepoIp }}:{{ $.artifactRepoPortQuay }}

- data located at: `/project-data/{{ $.dory.namespace }}/{{ $.dory.artifactRepo.type }}`
{{- end }}

{{- if $.scanCodeRepoInternal }}

### {{ $.dory.scanCodeRepo.type }} scan code repository

- url: {{ $.scanCodeRepoViewUrl }}
- admin user: admin / {{ $.scanCodeRepoPassword }}
{{- end }}

### openldap account management

- url: {{ $.viewURL | replace "http://" "https://" }}:{{ $.dory.openldap.port }}
- user: cn=admin,{{ $.dory.openldap.baseDN }} / {{ $.dory.openldap.password }}

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

### caution: this folder is very important, please keep it
