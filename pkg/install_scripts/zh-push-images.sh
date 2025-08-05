# 在 {{ $.dory.imageRepo.type }} 中创建 public, hub, gcr, quay 四个项目
curl -k -X POST -u {{ $.imageRepoUsername }}:{{ $.imageRepoPassword }} -H 'Content-Type: application/json' -d '{"project_name": "public", "public": true}' 'https://{{ $.imageRepoDomainName }}/api/v2.0/projects'
curl -k -X POST -u {{ $.imageRepoUsername }}:{{ $.imageRepoPassword }} -H 'Content-Type: application/json' -d '{"project_name": "hub", "public": true}' 'https://{{ $.imageRepoDomainName }}/api/v2.0/projects'
curl -k -X POST -u {{ $.imageRepoUsername }}:{{ $.imageRepoPassword }} -H 'Content-Type: application/json' -d '{"project_name": "gcr", "public": true}' 'https://{{ $.imageRepoDomainName }}/api/v2.0/projects'
curl -k -X POST -u {{ $.imageRepoUsername }}:{{ $.imageRepoPassword }} -H 'Content-Type: application/json' -d '{"project_name": "quay", "public": true}' 'https://{{ $.imageRepoDomainName }}/api/v2.0/projects'

# 设置只有管理员可以创建项目
curl -k -X PUT -u {{ $.imageRepoUsername }}:{{ $.imageRepoPassword }} -H 'Content-Type: application/json' -d '{"project_creation_restriction": "adminonly"}' 'https://{{ $.imageRepoDomainName }}/api/v2.0/configurations'

# 设置 {{ $.kubernetes.runtime }} 客户端登录到 {{ $.dory.imageRepo.type }}
{{ $.cmdLogin }} --username {{ $.imageRepoUsername }} --password {{ $.imageRepoPassword }} {{ $.imageRepoDomainName }}

# 把之前拉取的容器镜像推送到 {{ $.dory.imageRepo.type }}
{{- range $_, $image := $.dockerImages }}
{{ $.cmdTag }} {{ if $image.dockerFile }}{{ $image.built }}{{ else }}{{ $image.source }}{{ end }} {{ $.imageRepoDomainName }}/{{ $image.target }}
{{- if $image.arm64 }}
{{ $.cmdTag }} {{ if $image.dockerFile }}{{ $image.built }}-arm64v8{{ else }}{{ $image.arm64 }}{{ end }} {{ $.imageRepoDomainName }}/{{ $image.target }}-arm64v8
{{- end }}
{{- end }}

{{ range $_, $image := $.dockerImages }}
{{ $.cmdPush }} {{ $.imageRepoDomainName }}/{{ $image.target }}
{{- if $image.arm64 }}
{{ $.cmdPush }} {{ $.imageRepoDomainName }}/{{ $image.target }}-arm64v8
{{- end }}
{{- end }}
