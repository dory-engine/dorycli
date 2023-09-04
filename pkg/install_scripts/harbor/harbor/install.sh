#!/bin/bash

set -e

DIR="$(cd "$(dirname "$0")" && pwd)"
source $DIR/common.sh

set +o noglob

usage=$'Please set hostname and other necessary attributes in harbor.yml first. DO NOT use localhost or 127.0.0.1 for hostname, because Harbor needs to be accessed by external clients.
Please set --with-notary if needs enable Notary in Harbor, and set ui_url_protocol/ssl_cert/ssl_cert_key in harbor.yml bacause notary must run under https. 
Please set --with-trivy if needs enable Trivy in Harbor.
Please do NOT set --with-chartmuseum, as chartmusuem has been deprecated and removed.'
item=0

# notary is not enabled by default
with_notary=$false
# clair is deprecated
with_clair=$false
# trivy is not enabled by default
with_trivy=$false

# flag to using docker compose v1 or v2, default would using v1 docker-compose
DOCKER_COMPOSE=docker-compose

while [ $# -gt 0 ]; do
        case $1 in
            --help)
            note "$usage"
            exit 0;;
            --with-notary)
            with_notary=true;;
            --with-clair)
            with_clair=true;;
            --with-trivy)
            with_trivy=true;;
            *)
            note "$usage"
            exit 1;;
        esac
        shift || true
done

if [ $with_clair ]
then
    error "Clair is deprecated please remove it from installation arguments !!!"
    exit 1
fi

workdir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $workdir

h2 "[Step $item]: checking if docker is installed ..."; let item+=1
check_docker

h2 "[Step $item]: checking docker-compose is installed ..."; let item+=1
check_dockercompose

if [ -f harbor*.tar.gz ]
then
    h2 "[Step $item]: loading Harbor images ..."; let item+=1
    docker load -i ./harbor*.tar.gz
fi
echo ""

h2 "[Step $item]: preparing environment ...";  let item+=1
if [ -n "$host" ]
then
    sed "s/^hostname: .*/hostname: $host/g" -i ./harbor.yml
fi

h2 "[Step $item]: preparing harbor configs ...";  let item+=1
prepare_para=
if [ $with_notary ] 
then
    prepare_para="${prepare_para} --with-notary"
fi
if [ $with_trivy ]
then
    prepare_para="${prepare_para} --with-trivy"
fi

./prepare $prepare_para
echo ""

if [ -n "$DOCKER_COMPOSE ps -q"  ]
    then
        note "stopping existing Harbor instance ..." 
        $DOCKER_COMPOSE down -v
fi
echo ""

h2 "[Step $item]: starting Harbor ..."
if [ $with_notary ]
then
    warn "
    Notary will be deprecated as of Harbor v2.6.0 and start to be removed in v2.8.0 or later.
    You can use cosign for signature instead since Harbor v2.5.0.
    Please see discussion here for more details. https://github.com/goharbor/harbor/discussions/16612"
fi

$DOCKER_COMPOSE up -d

success $"----Harbor has been installed and started successfully.----"
