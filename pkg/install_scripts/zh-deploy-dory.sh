# 创建 {{ $.dory.namespace }} 组件的名字空间与pv
kubectl apply -f {{ $.dory.namespace }}/step01-pv.yaml

# 在kubernetes中部署dory组件
kubectl apply -f {{ $.dory.namespace }}/step02-statefulset.yaml
kubectl apply -f {{ $.dory.namespace }}/step03-service.yaml
kubectl apply -f {{ $.dory.namespace }}/step04-networkpolicy.yaml

# 检查dory服务状态
sh pods-ready.sh {{ $.dory.namespace }}
