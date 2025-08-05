# 清除kubernetes方式部署的服务

## 当安装出现异常的情况下，清除所有dory服务

{{- if $.imageRepoInternal }}
### 停止并清除 {{ $.dory.imageRepo.type }} 服务

```shell script
helm -n {{ $.dory.imageRepo.internal.namespace }} uninstall {{ $.dory.imageRepo.internal.namespace }}
```
{{- end }}

### 停止并清除所有 dory 服务

```shell script
# 删除所有dory数据
kubectl -n {{ $.dory.namespace }} exec project-data-pod-0 -- sh -c "rm -rf /project-data/*"

# 删除相关名字空间
kubectl delete namespace {{ $.dory.namespace }}
kubectl delete pv {{ $.dory.namespace }}-pv
{{- if $.imageRepoInternal }}
kubectl delete namespace {{ $.dory.imageRepo.internal.namespace }}
kubectl delete pv {{ $.dory.imageRepo.internal.namespace }}-pv
{{- end }}

kubectl delete pv {{ $.dory.namespace }}-project-data-pv
```
