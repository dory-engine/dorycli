# for docker mode install only
export HARBOR_CONFIG_DOMAIN_NAME={{ $.imageRepoDomainName }}
export HARBOR_CONFIG_CERT_PATH={{ $.dory.imageRepo.internal.certsDir }}

rm -rf ${HARBOR_CONFIG_CERT_PATH}/
mkdir -p ${HARBOR_CONFIG_CERT_PATH}/
cd ${HARBOR_CONFIG_CERT_PATH}/

openssl genrsa -out ca.key 4096
openssl req -x509 -new -nodes -sha512 -days 3650 -subj "/CN=${HARBOR_CONFIG_DOMAIN_NAME}" -key ca.key -out ca.crt
openssl genrsa -out ${HARBOR_CONFIG_DOMAIN_NAME}.key 4096
openssl req -sha512 -new -subj "/CN=${HARBOR_CONFIG_DOMAIN_NAME}" -key ${HARBOR_CONFIG_DOMAIN_NAME}.key -out ${HARBOR_CONFIG_DOMAIN_NAME}.csr
cat > v3.ext <<-EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[alt_names]
DNS.1=${HARBOR_CONFIG_DOMAIN_NAME}
EOF
openssl x509 -req -sha512 -days 3650 -extfile v3.ext -CA ca.crt -CAkey ca.key -CAcreateserial -in ${HARBOR_CONFIG_DOMAIN_NAME}.csr -out ${HARBOR_CONFIG_DOMAIN_NAME}.crt
openssl x509 -inform PEM -in ${HARBOR_CONFIG_DOMAIN_NAME}.crt -out ${HARBOR_CONFIG_DOMAIN_NAME}.cert
# echo "[INFO] # check harbor certificates info"
# echo "[CMD] openssl x509 -noout -text -in ${HARBOR_CONFIG_DOMAIN_NAME}.crt"
# openssl x509 -noout -text -in ${HARBOR_CONFIG_DOMAIN_NAME}.crt

# update /etc/docker/certs.d/ harbor certificates
echo "update docker harbor certificates"
rm -rf /etc/docker/certs.d/${HARBOR_CONFIG_DOMAIN_NAME}/
mkdir -p /etc/docker/certs.d/${HARBOR_CONFIG_DOMAIN_NAME}/
cp ${HARBOR_CONFIG_DOMAIN_NAME}.cert /etc/docker/certs.d/${HARBOR_CONFIG_DOMAIN_NAME}/
cp ${HARBOR_CONFIG_DOMAIN_NAME}.key /etc/docker/certs.d/${HARBOR_CONFIG_DOMAIN_NAME}/
cp ca.crt /etc/docker/certs.d/${HARBOR_CONFIG_DOMAIN_NAME}/
ls -al /etc/docker/certs.d/${HARBOR_CONFIG_DOMAIN_NAME}/

echo "update kubernetes {{ $.kubernetes.runtime }} harbor certificates"
{{ $certPath := "" }}{{- if eq $.kubernetes.runtime "docker" }}{{ $certPath = "/etc/docker" }}{{- else if eq $.kubernetes.runtime "containerd" }}{{ $certPath = "/etc/containerd" }}{{- else if eq $.kubernetes.runtime "crio" }}{{ $certPath = "/etc/containers" }}{{- end }}
rm -rf {{ $certPath }}/certs.d/${HARBOR_CONFIG_DOMAIN_NAME}/
mkdir -p {{ $certPath }}/certs.d/${HARBOR_CONFIG_DOMAIN_NAME}/
cp ${HARBOR_CONFIG_DOMAIN_NAME}.cert {{ $certPath }}/certs.d/${HARBOR_CONFIG_DOMAIN_NAME}/
cp ${HARBOR_CONFIG_DOMAIN_NAME}.key {{ $certPath }}/certs.d/${HARBOR_CONFIG_DOMAIN_NAME}/
cp ca.crt {{ $certPath }}/certs.d/${HARBOR_CONFIG_DOMAIN_NAME}/
ls -al {{ $certPath }}/certs.d/${HARBOR_CONFIG_DOMAIN_NAME}/

cd ..
