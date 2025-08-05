# create {{ $.dory.namespace }} namespace and pv
kubectl apply -f {{ $.dory.namespace }}/step01-pv.yaml

# start all dory services with kubernetes
kubectl apply -f {{ $.dory.namespace }}/step02-statefulset.yaml
kubectl apply -f {{ $.dory.namespace }}/step03-service.yaml
kubectl apply -f {{ $.dory.namespace }}/step04-networkpolicy.yaml

# check dory services status
sh pods-ready.sh {{ $.dory.namespace }}
