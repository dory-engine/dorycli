# copy dory-engine config file to share storage
kubectl -n {{ $.dory.namespace }} cp install-data/{{ $.dory.namespace }}/dory-engine/config/config.yaml project-data-pod-0:/project-data/{{ $.dory.namespace }}/dory-engine/config/

# restart dory-engine and dory-console service
kubectl -n {{ $.dory.namespace }} scale statefulsets.apps dory-engine dory-console --replicas 0
kubectl -n {{ $.dory.namespace }} scale statefulsets.apps dory-engine dory-console --replicas 1

# waiting for dory-engine-0 dory-console-0 ready
sh pods-ready.sh {{ $.dory.namespace }}
