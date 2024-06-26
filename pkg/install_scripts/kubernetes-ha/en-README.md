# High availability kubernetes cluster deployment

- Please refer to the documentation for installation details: [https://github.com/cookeem/kubeadm-ha](https://github.com/cookeem/kubeadm-ha/blob/master/README-EN.md)

## The directory structure is as follows

```text
{{- range $_, $host := $.masterHosts }}
├── {{ $host.hostname }} # Please copy this directory to {{ $host.hostname }} node host
│   ├── keepalived # docker-compose file and configuration file directory of keepalived service
│   │   ├── check_apiserver.sh # kubernetes apiserver check script for keepalived
│   │   ├── docker-compose.yaml # Use 'docker-compose up -d' to start the keepalived service
│   │   └── keepalived.conf # keepalived configuration file
│   └── nginx-lb # docker-compose file and configuration file directory of nginx-lb service
│       ├── docker-compose.yaml # Use 'docker-compose up -d' to start nginx-lb service
│       └── nginx-lb.conf # nginx-lb configuration file
{{- end }}
└── kubeadm-config.yaml # kubeadm high availability cluster initialization configuration file
```

## Execute the following command to start the load balancer of the kubernetes high-availability cluster on each master node

```bash
# Set the path of the kubernetes high-availability cluster load balancer of each master node
export LB_DIR=/data/k8s-lb
{{ range $i, $host := $.masterHosts }}
# Copy the load balancer configuration file to {{ $host.hostname }} node
ssh {{ $host.hostname }} mkdir -p ${LB_DIR} && \
scp -r {{ $host.hostname }}/nginx-lb {{ $host.hostname }}/keepalived root@{{ $host.hostname }}:${LB_DIR} && \
pwd

# Start load balancer on {{ $host.hostname }} node
ssh {{ $host.hostname }} "cd ${LB_DIR}/keepalived/ && docker-compose stop && docker-compose rm -f && docker-compose up -d" && \
ssh {{ $host.hostname }} "cd ${LB_DIR}/nginx-lb/ && docker-compose stop && docker-compose rm -f && docker-compose up -d" && \
pwd
{{ end }}
{{ $firstHost := first $.masterHosts }}

# Execute kubernetes controll-plane initialization on the first master node
kubeadm init --config=kubeadm-config.yaml --upload-certs
```