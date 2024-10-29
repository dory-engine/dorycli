{{ $certPath := "" }}{{- if eq $.kubernetes.runtime "docker" }}{{ $certPath = "/etc/docker" }}{{- else if eq $.kubernetes.runtime "containerd" }}{{ $certPath = "/etc/containerd" }}{{- else if eq $.kubernetes.runtime "crio" }}{{ $certPath = "/etc/containers" }}{{- end }}
export CERT_PATH="{{ $certPath }}/certs.d"
rm -rf ${CERT_PATH}/{{ $.imageRepoDomainName }}
mkdir -p ${CERT_PATH}/{{ $.imageRepoDomainName }}
export INGRESS_SECRET_NAME=$(kubectl -n {{ $.dory.imageRepo.internal.namespace }} get secrets | grep "ingress" | awk '{print $1}')
kubectl -n {{ $.dory.imageRepo.internal.namespace }} get secrets ${INGRESS_SECRET_NAME} -o jsonpath='{ .data.ca\.crt }'  | base64 -d > ${CERT_PATH}/{{ $.imageRepoDomainName }}/ca.crt
kubectl -n {{ $.dory.imageRepo.internal.namespace }} get secrets ${INGRESS_SECRET_NAME} -o jsonpath='{ .data.tls\.crt }' | base64 -d > ${CERT_PATH}/{{ $.imageRepoDomainName }}/{{ $.imageRepoDomainName }}.cert
kubectl -n {{ $.dory.imageRepo.internal.namespace }} get secrets ${INGRESS_SECRET_NAME} -o jsonpath='{ .data.tls\.key }' | base64 -d > ${CERT_PATH}/{{ $.imageRepoDomainName }}/{{ $.imageRepoDomainName }}.key
