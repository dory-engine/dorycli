# create public, hub, gcr, quay projects in {{ $.dory.imageRepo.type }}
curl -k -X POST -u {{ $.imageRepoUsername }}:{{ $.imageRepoPassword }} -H 'Content-Type: application/json' -d '{"project_name": "public", "public": true}' 'https://{{ $.imageRepoDomainName }}/api/v2.0/projects'
curl -k -X POST -u {{ $.imageRepoUsername }}:{{ $.imageRepoPassword }} -H 'Content-Type: application/json' -d '{"project_name": "hub", "public": true}' 'https://{{ $.imageRepoDomainName }}/api/v2.0/projects'
curl -k -X POST -u {{ $.imageRepoUsername }}:{{ $.imageRepoPassword }} -H 'Content-Type: application/json' -d '{"project_name": "gcr", "public": true}' 'https://{{ $.imageRepoDomainName }}/api/v2.0/projects'
curl -k -X POST -u {{ $.imageRepoUsername }}:{{ $.imageRepoPassword }} -H 'Content-Type: application/json' -d '{"project_name": "quay", "public": true}' 'https://{{ $.imageRepoDomainName }}/api/v2.0/projects'

# set configuration only admin can create projects
curl -k -X PUT -u {{ $.imageRepoUsername }}:{{ $.imageRepoPassword }} -H 'Content-Type: application/json' -d '{"project_creation_restriction": "adminonly"}' 'https://{{ $.imageRepoDomainName }}/api/v2.0/configurations'

# {{ $.kubernetes.runtime }} client login to {{ $.dory.imageRepo.type }}
{{ $.cmdLogin }} --username {{ $.imageRepoUsername }} --password {{ $.imageRepoPassword }} {{ $.imageRepoDomainName }}

# push container images to {{ $.dory.imageRepo.type }}
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
