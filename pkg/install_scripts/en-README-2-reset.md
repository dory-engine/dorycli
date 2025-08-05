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
# delete all dory data
kubectl -n {{ $.dory.namespace }} exec project-data-pod-0 -- sh -c "rm -rf /project-data/*"

# delete relative namespaces
kubectl delete namespace {{ $.dory.namespace }}
kubectl delete pv {{ $.dory.namespace }}-pv
{{- if $.imageRepoInternal }}
kubectl delete namespace {{ $.dory.imageRepo.internal.namespace }}
kubectl delete pv {{ $.dory.imageRepo.internal.namespace }}-pv
{{- end }}

kubectl delete pv {{ $.dory.namespace }}-project-data-pv
```
