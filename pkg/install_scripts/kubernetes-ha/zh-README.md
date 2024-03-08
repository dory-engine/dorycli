# 高可用kubernetes集群部署

- 安装详细参见文档: [https://github.com/cookeem/kubeadm-ha](https://github.com/cookeem/kubeadm-ha/blob/master/README.md)

## 目录结构如下

```text
{{- range $_, $host := $.masterHosts }}
├── {{ $host.hostname }} # 请把该目录复制到 {{ $host.hostname }} 节点主机
│   ├── keepalived # keepalived服务的docker-compose文件以及配置文件目录
│   │   ├── check_apiserver.sh # keepalived的kubernetes apiserver检查脚本
│   │   ├── docker-compose.yaml # 使用 'docker-compose up -d' 启动keepalived服务
│   │   └── keepalived.conf # keepalived配置文件
│   └── nginx-lb # nginx-lb服务的docker-compose文件以及配置文件目录
│       ├── docker-compose.yaml # 使用 'docker-compose up -d' 启动nginx-lb服务
│       └── nginx-lb.conf # nginx-lb配置文件
{{- end }}
└── kubeadm-config.yaml # kubeadm的高可用集群初始化配置文件
```

## 执行以下命令，在各个master节点启动kubernetes高可用集群的load balancer

```bash
# 设置各个master节点的kubernetes高可用集群load balancer的路径
export LB_DIR=/data/k8s-lb
{{ range $i, $host := $.masterHosts }}
# 把load balancer配置文件复制到 {{ $host.hostname }} 节点上
ssh {{ $host.hostname }} mkdir -p ${LB_DIR}
scp -r {{ $host.hostname }}/nginx-lb {{ $host.hostname }}/keepalived root@{{ $host.hostname }}:${LB_DIR}

# 在 {{ $host.hostname }} 节点上启动load balancer
ssh {{ $host.hostname }} "cd ${LB_DIR}/keepalived/ && docker-compose stop && docker-compose rm -f && docker-compose up -d"
ssh {{ $host.hostname }} "cd ${LB_DIR}/nginx-lb/ && docker-compose stop && docker-compose rm -f && docker-compose up -d"
{{ end }}
{{ $firstHost := first $.masterHosts }}

# 在第一个master节点执行kubernetes controll-plane 初始化
kubeadm init --config=kubeadm-config.yaml --upload-certs
```