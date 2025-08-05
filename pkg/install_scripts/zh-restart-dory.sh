# 把dory-engine配置文件复制到共享存储
kubectl -n {{ $.dory.namespace }} cp install-data/{{ $.dory.namespace }}/dory-engine/config/config.yaml project-data-pod-0:/project-data/{{ $.dory.namespace }}/dory-engine/config/

# 重启 dory-engine 和 dory-console 服务
kubectl -n {{ $.dory.namespace }} scale statefulsets.apps dory-engine dory-console --replicas 0
kubectl -n {{ $.dory.namespace }} scale statefulsets.apps dory-engine dory-console --replicas 1

# 等待 dory-engine-0 dory-console-0 pod处于ready状态
sh pods-ready.sh {{ $.dory.namespace }}
