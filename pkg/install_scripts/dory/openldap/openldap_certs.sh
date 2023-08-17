export DORY_OPENLDAP_NAME={{ $.dory.openldap.serviceName }}
export DORY_OPENLDAP_NAMESPACE={{ $.dory.namespace }}

openssl dhparam -out dhparam.pem 2048
openssl genrsa -out ca.key 4096
openssl req -x509 -new -nodes -sha512 -days 3650 -subj "/CN=${DORY_OPENLDAP_NAME}" -key ca.key -out ca.crt
openssl genrsa -out ldap.key 4096
openssl req -sha512 -new -subj "/CN=${DORY_OPENLDAP_NAME}" -key ldap.key -out ldap.csr
cat << EOF > v3.ext
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
extendedKeyUsage = serverAuth, clientAuth
subjectAltName = @alt_names

[alt_names]
DNS.1=${DORY_OPENLDAP_NAME}-0
DNS.2=*.${DORY_OPENLDAP_NAME}
DNS.3=${DORY_OPENLDAP_NAME}
DNS.4=*.${DORY_OPENLDAP_NAMESPACE}
DNS.5=localhost
EOF
openssl x509 -req -sha512 -days 3650 -extfile v3.ext -CA ca.crt -CAkey ca.key -CAcreateserial -in ldap.csr -out ldap.crt
chown -R 911:911 ..

