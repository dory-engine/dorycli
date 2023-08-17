# reset kubernetes installation

## remove all dory service when install failure

{{- if $.imageRepoInternal }}
### stop and remove {{ $.dory.imageRepo.type }} services

```shell script
helm -n {{ $.dory.imageRepo.internal.namespace }} uninstall {{ $.dory.imageRepo.internal.namespace }}
```
{{- end }}

### stop and remove dory services

```shell script
kubectl delete namespace {{ $.dory.namespace }}
kubectl delete pv {{ $.dory.namespace }}-pv
kubectl delete pv {{ $.dory.namespace }}-timezone-pv
kubectl delete pv project-data-pv
kubectl delete pv project-data-timezone-pv
{{- if $.imageRepoInternal }}
kubectl delete namespace {{ $.dory.imageRepo.internal.namespace }}
kubectl delete pv {{ $.dory.imageRepo.internal.namespace }}-pv
kubectl delete pv {{ $.dory.imageRepo.internal.namespace }}-timezone-pv
{{- end }}
```

## about dory services data

- dory services data located at: `{{ $.rootDir }}`

```shell script
# before reinstall, please remove dory services data first
rm -rf {{ $.rootDir }}/*
```
