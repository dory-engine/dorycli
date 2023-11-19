# kubernetes环境部署要求

{{- if eq $.mode "docker" }}
## docker模式安装注意事项

- docker模式安装的主机不能是kubernetes的master节点，docker模式启动的harbor会与kubernetes的master节点的ingress controller抢夺443端口，会引起harbor访问异常
- 要使用docker模式安装，请保证你至少有两台主机
{{- end }}

## 系统硬件资源需求

### 安装DORY核心组件 (默认安装)

- cpus: 1核
- memory: 1G
- storage: 2G

### 安装所有可选组件 (完整安装)

- cpus: 4核
- memory: 16G
- storage: 60G

## 在kubernetes集群中创建管理token

- [注意] 请保证本机的kubectl能够管理目标kubernetes集群

- kubernetes管理token用于dory连接kubernetes集群并发布应用，必须在dory配置文件中设置

```shell script
# 创建管理员serviceaccount
kubectl create serviceaccount -n kube-system admin-user --dry-run=client -o yaml | kubectl apply -f -

# 创建管理员clusterrolebinding
kubectl create clusterrolebinding admin-user --clusterrole=cluster-admin --serviceaccount=kube-system:admin-user --dry-run=client -o yaml | kubectl apply -f -

# 需要手动创建serviceaccount的secret
cat << EOF | kubectl apply -f -
apiVersion: v1
kind: Secret
metadata:
  name: admin-user-secret
  namespace: kube-system
  annotations:
    kubernetes.io/service-account.name: admin-user
type: kubernetes.io/service-account-token
EOF

# 获取kubernetes管理token `kubernetes.token`
# kubernetes管理token需要在dory安装过程进行设置
kubectl -n kube-system get secret admin-user-secret -o jsonpath='{ .data.token }' | base64 -d

# 获取kubernetes集群管理员ca证书的base64编码字符串 `kubernetes.caCrtBase64`
# kubernetes集群管理员ca证书需要在dory安装过程进行设置
kubectl config view --raw -o=jsonpath='{.clusters[0].cluster.certificate-authority-data}'
```

## x86架构和arm64架构容器镜像交叉构建需求 (可选)

- {{ if eq $.mode "docker" }}本节点{{ else }}所有部署DORY的节点{{ end }}都安装`qemu-user-static`，以保证这些节点都能够运行x86架构和arm64架构容器镜像
- 文档参见以下链接: https://github.com/multiarch/qemu-user-static

- 升级linux操作系统内核为最新版本，保证linux kernel支持`binfmt_misc`

- 安装qemu
```shell script
# apt-get方式安装qemu
apt-get install qemu

# yum方式安装qemu
yum install -y qemu
```

- 执行qemu-user-static，让节点都能够运行x86架构和arm64架构容器镜像
```shell script
{{- if or (eq $.mode "docker") (eq $.runtime "docker") }}
docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
{{- else if eq $.runtime "containerd" }}
nerdctl -n k8s.io run --rm --privileged multiarch/qemu-user-static --reset -p yes
{{- else if eq $.runtime "crio" }}
podman run --rm --privileged multiarch/qemu-user-static --reset -p yes
{{- end }}
```

- 验证节点是否支持运行x86架构和arm64架构容器镜像
```shell script
{{- if or (eq $.mode "docker") (eq $.runtime "docker") }}
docker run --rm -t arm64v8/alpine:latest uname -m
{{- else if eq $.runtime "containerd" }}
nerdctl -n k8s.io run --rm -t arm64v8/alpine:latest uname -m
{{- else if eq $.runtime "crio" }}
podman run --rm -t arm64v8/alpine:latest uname -m
{{- end }}
```

{{- if eq $.runtime "containerd" }}
## 所有kubernetes节点设置containerd的镜像仓库自签名证书路径

```shell script
# # 查找并添加config_path配置
# vi /etc/containerd/config.toml
#     [plugins."io.containerd.grpc.v1.cri".registry]
#       config_path = "/etc/containerd/certs.d"
#             
# # 创建镜像仓库证书目录
# mkdir -p /etc/containerd/certs.d
# 
# # 重启containerd
# systemctl restart containerd
```
{{- end }}

## kubernetes-dashboard

- 为了管理kubernetes中部署的应用，推荐使用`kubernetes-dashboard`
- 要了解更多，请阅读官方代码仓库README.md文档: [kubernetes-dashboard](https://github.com/kubernetes/dashboard)

- 安装:
```shell script
# 安装 kubernetes-dashboard
kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.7.0/aio/deploy/recommended.yaml

# 暴露 kubernetes-dashboard 的服务端口为 nodePort 30000
cat << EOF | kubectl apply -f -
apiVersion: v1
kind: Service
metadata:
  labels:
    k8s-app: kubernetes-dashboard
  name: kubernetes-dashboard
  namespace: kubernetes-dashboard
spec:
  ports:
  - nodePort: 30000
    port: 443
    protocol: TCP
    targetPort: 8443
  selector:
    k8s-app: kubernetes-dashboard
  type: NodePort
EOF
```

## traefik (ingress controller)

- 要使用kubernetes的[ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/)功能，必须安装ingress controller，推荐使用`traefik`
- 要了解更多，请阅读官方网站文档: [traefik](https://doc.traefik.io/traefik/)

- 在kubernetes所有master节点部署traefik: 
```shell script
# 拉取 traefik helm repo
helm repo add traefik https://helm.traefik.io/traefik
helm fetch traefik/traefik --untar

# 在kubernetes的master节点以daemonset方式部署traefik
cat << EOF > traefik.yaml
deployment:
  kind: DaemonSet
image:
  name: traefik
  tag: v2.10.5
pilot:
  enabled: true
experimental:
  plugins:
    enabled: true
ports:
  web:
    hostPort: 80
  websecure:
    hostPort: 443
service:
  type: ClusterIP
EOF

# 安装traefik
kubectl create namespace traefik --dry-run=client -o yaml | kubectl apply -f -
helm install -n traefik traefik traefik/ -f traefik.yaml

# 检查安装情况
helm -n traefik list
kubectl -n traefik get pods -o wide
kubectl -n traefik get services -o wide
```

## metrics-server

- 为了使用kubernetes的水平扩展缩容功能[horizontal pod autoscale](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/)，必须安装`metrics-server`
- 要了解更多，请阅读官方代码仓库README.md文档: [metrics-server](https://github.com/kubernetes-sigs/metrics-server)

- install:
```shell script
# 拉取镜像
{{ $.cmdImagePull }} registry.aliyuncs.com/google_containers/metrics-server:v0.6.4
{{ $.cmdImageTag }} registry.aliyuncs.com/google_containers/metrics-server:v0.6.4 registry.k8s.io/metrics-server/metrics-server:v0.6.4

# 获取metrics-server安装yaml
curl -O -L https://github.com/kubernetes-sigs/metrics-server/releases/download/v0.6.4/components.yaml
# 添加--kubelet-insecure-tls参数
sed -i 's/- args:/- args:\n        - --kubelet-insecure-tls/g' components.yaml
# 安装metrics-server
kubectl apply -f components.yaml

# 等待metrics-server正常
kubectl -n kube-system get pods -l=k8s-app=metrics-server

# 查看pod的metrics
kubectl top pods -A

```

## istio

- 要使用服务网格的混合灰度发布能力，需要部署istio服务网格
- 要了解更多，请阅读istio官网文档: [istio.io](https://istio.io/latest/docs/)

- install:
```shell script
# 安装istioctl，这里以v1.19.3为例子，客户端下载地址 https://github.com/istio/istio/releases/tag/1.19.3

# 下载istioctl
wget https://github.com/istio/istio/releases/download/1.19.3/istioctl-1.19.3-linux-amd64.tar.gz
tar zxvf istioctl-1.19.3-linux-amd64.tar.gz

# 把istioctl移动到$PATH对应目录
mv istioctl /usr/bin/

# 确认istioctl版本
istioctl version

# 使用istioctl部署istio到kubernetes
istioctl install --set profile=demo \
--set values.gateways.istio-ingressgateway.type=ClusterIP \
--set values.global.imagePullPolicy=IfNotPresent \
--set values.global.proxy_init.resources.limits.cpu=100m \
--set values.global.proxy_init.resources.limits.memory=100Mi \
--set values.global.proxy.resources.limits.cpu=100m \
--set values.global.proxy.resources.limits.memory=100Mi

# 检查istio部署情况
kubectl -n istio-system get pods,svc
NAME                                       READY   STATUS    RESTARTS   AGE
pod/istio-egressgateway-599c8845c9-lcs68   1/1     Running   0          15h
pod/istio-ingressgateway-69dc56d7f-cscwh   1/1     Running   0          15h
pod/istiod-8c75fcbc9-qv9mn                 1/1     Running   0          15h

NAME                           TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)                                        AGE
service/istio-egressgateway    ClusterIP   10.111.238.12    <none>        80/TCP,443/TCP                                 15h
service/istio-ingressgateway   ClusterIP   10.103.206.173   <none>        15021/TCP,80/TCP,443/TCP,31400/TCP,15443/TCP   15h
service/istiod                 ClusterIP   10.103.41.209    <none>        15010/TCP,15012/TCP,443/TCP,15014/TCP          15h
```

## sonarqube内核参数调整 (可选，不安装sonarqube请忽略)

```shell script
# sonarqube 部署必须设置linux内核参数: vm.max_map_count = 262144

cat <<EOF >  /etc/sysctl.d/sonarqube.conf
vm.max_map_count = 262144
EOF

# 在所有sonarqube节点设置启用sysctl
sysctl --system
```
