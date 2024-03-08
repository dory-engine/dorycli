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
kubectl delete namespace {{ $.dory.namespace }}
kubectl delete pv {{ $.dory.namespace }}-pv
kubectl delete pv project-data-pv
{{- if $.imageRepoInternal }}
kubectl delete namespace {{ $.dory.imageRepo.internal.namespace }}
kubectl delete pv {{ $.dory.imageRepo.internal.namespace }}-pv
{{- end }}
```

## 所有dory组件的数据存放位置

- 所有dory组件的数据存放在: `{{ $.rootDir }}`

```shell script
# 重新安装前，请清理dory组件数据
rm -rf {{ $.rootDir }}/*
```
