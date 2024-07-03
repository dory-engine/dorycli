# kubernetes prerequisite

{{- if eq $.mode "docker" }}
## caution about docker install mode 

- docker mode install machine and kubernetes master can't be the same machine, because in docker mode harbor services and kubernetes master ingress controller use the same 443 TLS port, port conflicts will cause dory connect harbor failed
- make sure you have 2 hosts for docker mode install
{{- end }}

## hardware requirement

### Install DORY core components (default installation)

- cpus: 1 core
- memory: 1G
- storage: 2G

### Install all optional components (full installation)

- cpus: 4 cores
- memory: 16G
- storage: 60G

## create kubernetes admin token

- [note] Please ensure that the local kubectl can manage the target kubernetes cluster

- kubernetes admin token is for dory to deploy project applications in kubernetes cluster, you must set it in dory's config file

```shell script
# create admin serviceaccount
kubectl create serviceaccount -n kube-system admin-user --dry-run=client -o yaml | kubectl apply -f -

# create administrator clusterrolebinding
kubectl create clusterrolebinding admin-user --clusterrole=cluster-admin --serviceaccount=kube-system:admin-user --dry-run=client -o yaml | kubectl apply -f -

# Need to manually create the secret of serviceaccount
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

# Get the kubernetes management token `kubernetes.token`
# The kubernetes management token needs to be set during the dory installation process
kubectl -n kube-system get secret admin-user-secret -o jsonpath='{ .data.token }' | base64 -d

# Get the kubernetes cluster administrator ca certificate base64 encoded string `kubernetes.caCrtBase64`
# The kubernetes cluster administrator ca certificate needs to be set during the dory installation process
kubectl config view --raw -o=jsonpath='{.clusters[0].cluster.certificate-authority-data}'
```

## X86 architecture and arm64 architecture container image cross-build requirements

- {{ if eq $.mode "docker" }}this node{{ else }}all nodes deploying DORY{{ end }} install `qemu-user-static` to ensure that these nodes can run x86 architecture and arm64 architecture container images
- see the following link for documentation: https://github.com/multiarch/qemu-user-static

- upgrade the linux operating system kernel to the latest version, make sure linux kernel support `binfmt_misc`

- install qemu
```shell script
# apt-get way to install qemu
apt-get install qemu

# yum way to install qemu
yum install -y qemu
```

- execute qemu-user-static, so that both nodes can run x86 architecture and arm64 architecture container images
```shell script
{{- if or (eq $.mode "docker") (eq $.runtime "docker") }}
docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
{{- else if eq $.runtime "containerd" }}
nerdctl -n k8s.io run --rm --privileged multiarch/qemu-user-static --reset -p yes
{{- else if eq $.runtime "crio" }}
podman run --rm --privileged multiarch/qemu-user-static --reset -p yes
{{- end }}
```

- verify that the node supports running x86 architecture and arm64 architecture container images
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
## set containerd image repository self signed certificates path on all kubernetes nodes

```shell script
# # find and add config_path settings
# vi /etc/containerd/config.toml
#     [plugins."io.containerd.grpc.v1.cri".registry]
#       config_path = "/etc/containerd/certs.d"
# 
# # create image repository certificates directory
# mkdir -p /etc/containerd/certs.d
# 
# # restart containerd service
# systemctl restart containerd
```
{{- end }}

## kubernetes-dashboard

- to manager your project pods in kubernetes, we recommend to use `kubernetes-dashboard`
- read official repository README.md to learn more: [kubernetes-dashboard](https://github.com/kubernetes/dashboard)

- install:
```shell script
# install kubernetes-dashboard
kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.7.0/aio/deploy/recommended.yaml

# expose kubernetes-dashboard service in nodePort 30000
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

- to use kubernetes [ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/), you need to install an ingress controller, we recommend to use `traefik`
- read official website to learn more: [traefik](https://doc.traefik.io/traefik/)

- install traefik in kubernetes master nodes:
```shell script
# fetch traefik helm repo
helm repo add traefik https://traefik.github.io/charts
helm fetch traefik/traefik --untar

# install traefik in kubernetes as daemonset on master nodes
cat << EOF > traefik.yaml
deployment:
  kind: DaemonSet
image:
  name: traefik
  tag: v2.10.5
ports:
  web:
    hostPort: 80
  websecure:
    hostPort: 443
service:
  type: ClusterIP
EOF

# install traefik
kubectl create namespace traefik --dry-run=client -o yaml | kubectl apply -f -
helm install -n traefik traefik traefik/ -f traefik.yaml

# check install
helm -n traefik list
kubectl -n traefik get pods -o wide
kubectl -n traefik get services -o wide
```

## metrics-server

- to use kubernetes [horizontal pod autoscale](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/), must install `metrics-server`
- read official repository README.md to learn more: [metrics-server](https://github.com/kubernetes-sigs/metrics-server)

- install:
```shell script
# pull container image
{{ $.cmdImagePull }} registry.aliyuncs.com/google_containers/metrics-server:v0.6.4
{{ $.cmdImageTag }} registry.aliyuncs.com/google_containers/metrics-server:v0.6.4 registry.k8s.io/metrics-server/metrics-server:v0.6.4

# get metrics-server install yaml
curl -O -L https://github.com/kubernetes-sigs/metrics-server/releases/download/v0.6.4/components.yaml
# add --kubelet-insecure-tls args
sed -i 's/- args:/- args:\n        - --kubelet-insecure-tls/g' components.yaml
# install metrics-server
kubectl apply -f components.yaml


# waiting for metrics-server ready
kubectl -n kube-system get pods -l=k8s-app=metrics-server

# get pods metrics
kubectl top pods -A
```

## istio

- to use the hybrid grayscale publishing capabilities of the service mesh, you need to deploy the `istio` service mesh
- read official website to learn more: [istio.io](https://istio.io/latest/docs/)

- install:
```shell script
# install istioctl, for example here use version v1.19.3, istioctl client tool download page: https://github.com/istio/istio/releases/tag/1.19.3

# down load istioctl
wget https://github.com/istio/istio/releases/download/1.19.3/istioctl-1.19.3-linux-amd64.tar.gz
tar zxvf istioctl-1.19.3-linux-amd64.tar.gz

# move istioctl to $PATH directory
mv istioctl /usr/bin/

# check istioctl version
istioctl version

# use istioctl to deploy istio in kubernetes
istioctl install --set profile=demo \
--set values.gateways.istio-ingressgateway.type=ClusterIP \
--set values.global.imagePullPolicy=IfNotPresent \
--set values.global.proxy_init.resources.limits.cpu=100m \
--set values.global.proxy_init.resources.limits.memory=100Mi \
--set values.global.proxy.resources.limits.cpu=100m \
--set values.global.proxy.resources.limits.memory=100Mi

# check istio deploy status
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

## sonarqube linux kernel param settings (optional, if not install sonarqube please ignore)

```shell script
# sonarqube required linux kernel param: vm.max_map_count = 262144

cat <<EOF >  /etc/sysctl.d/sonarqube.conf
vm.max_map_count = 262144
EOF

# run sysctl to apply linux kernel param in all sonarqube nodes
sysctl --system
```
