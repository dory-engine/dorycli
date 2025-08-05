#!/bin/bash

CONTAINER="gitlab-0"
GITLAB_ROOT_PASSWORD=$(kubectl -n {{ $.dory.namespace }} exec $CONTAINER -- cat /etc/gitlab/initial_root_password 2>/dev/null | grep "^Password:" | awk -F': ' '{print $2}')

if [ -n "$GITLAB_ROOT_PASSWORD" ]; then
    echo "✅ root初始密码: $GITLAB_ROOT_PASSWORD"
else
    echo "❌ 未能获取root初始密码"
fi

{{ $.cmdRun }} --rm -ti -v $PWD:/src doryengine/python:3.11.2-alpine3.17-dory python /src/gitlab-config.py --password $GITLAB_ROOT_PASSWORD
