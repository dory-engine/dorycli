export DORY_DOCKER_NAME={{ $.dory.docker.dockerName }}
export DORY_DOCKER_NAMESPACE={{ $.dory.namespace }}

sudo rm -rf docker-certs/
mkdir -p docker-certs/
cd docker-certs/

openssl genrsa -out ca.key 4096
openssl req -x509 -new -nodes -sha512 -days 3650 -subj "/CN=${DORY_DOCKER_NAME}" -key ca.key -out ca.crt
openssl genrsa -out tls.key 4096
openssl req -sha512 -new -subj "/CN=${DORY_DOCKER_NAME}" -key tls.key -out tls.csr
cat << EOF > v3.ext
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
extendedKeyUsage = serverAuth, clientAuth
subjectAltName = @alt_names

[alt_names]
DNS.1=${DORY_DOCKER_NAME}
DNS.2=*.${DORY_DOCKER_NAME}
DNS.3=*.${DORY_DOCKER_NAME}.${DORY_DOCKER_NAMESPACE}
DNS.4=*.${DORY_DOCKER_NAMESPACE}
DNS.5=localhost
EOF
openssl x509 -req -sha512 -days 3650 -extfile v3.ext -CA ca.crt -CAkey ca.key -CAcreateserial -in tls.csr -out tls.crt
# echo "[INFO] # check docker certificates info"
# echo "[CMD] openssl x509 -noout -text -in tls.crt"
# openssl x509 -noout -text -in tls.crt
cd ..
sudo chown -R 1000:1000 docker-certs
