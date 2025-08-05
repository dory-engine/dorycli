export NEUXS_INIT_PASSWORD=$(kubectl -n {{ $.dory.namespace }} exec nexus-0 -- cat /nexus-data/admin.password)
echo $NEUXS_INIT_PASSWORD
export NEUXS_PASSWORD={{ $.artifactRepoPassword }}

# 更新 admin 密码
curl -X 'PUT' -u "admin:$NEUXS_INIT_PASSWORD" \
  '{{ $.artifactRepoViewUrl }}/service/rest/v1/security/users/admin/change-password' \
  -H 'accept: application/json' \
  -H 'Content-Type: text/plain' \
  -d "$NEUXS_PASSWORD"

# 允许匿名访问
curl -X 'PUT' -u "admin:$NEUXS_PASSWORD" \
  '{{ $.artifactRepoViewUrl }}/service/rest/v1/security/anonymous' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "enabled": true,
  "userId": "anonymous",
  "realmName": "NexusAuthorizingRealm"
}'

# 创建 docker proxy 仓库
curl -X 'POST' -u "admin:$NEUXS_PASSWORD" \
  '{{ $.artifactRepoViewUrl }}/service/rest/v1/repositories/docker/proxy' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "name": "docker-hub",
  "online": true,
  "storage": {
    "blobStoreName": "default",
    "strictContentTypeValidation": true
  },
  "proxy": {
    "remoteUrl": "https://registry-1.docker.io",
    "contentMaxAge": 1440,
    "metadataMaxAge": 1440
  },
  "negativeCache": {
    "enabled": true,
    "timeToLive": 1440
  },
  "httpClient": {
    "blocked": false,
    "autoBlock": true,
    "connection": {
      "enableCircularRedirects": false,
      "enableCookies": false,
      "useTrustStore": false
    }
  },
  "docker": {
    "v1Enabled": true,
    "forceBasicAuth": true,
    "httpPort": 1443
  },
  "dockerProxy": {
    "indexType": "HUB",
    "cacheForeignLayers": false,
    "foreignLayerUrlWhitelist": []
  }
}'

curl -X 'POST' -u "admin:$NEUXS_PASSWORD" \
  '{{ $.artifactRepoViewUrl }}/service/rest/v1/repositories/docker/proxy' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "name": "docker-gcr",
  "online": true,
  "storage": {
    "blobStoreName": "default",
    "strictContentTypeValidation": true
  },
  "proxy": {
    "remoteUrl": "https://registry.cn-hangzhou.aliyuncs.com",
    "contentMaxAge": 1440,
    "metadataMaxAge": 1440
  },
  "negativeCache": {
    "enabled": true,
    "timeToLive": 1440
  },
  "httpClient": {
    "blocked": false,
    "autoBlock": true,
    "connection": {
      "enableCircularRedirects": false,
      "enableCookies": false,
      "useTrustStore": false
    }
  },
  "docker": {
    "v1Enabled": true,
    "forceBasicAuth": true,
    "httpPort": 1444
  },
  "dockerProxy": {
    "indexType": "REGISTRY",
    "cacheForeignLayers": false,
    "foreignLayerUrlWhitelist": []
  }
}'

curl -X 'POST' -u "admin:$NEUXS_PASSWORD" \
  '{{ $.artifactRepoViewUrl }}/service/rest/v1/repositories/docker/proxy' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "name": "docker-quay",
  "online": true,
  "storage": {
    "blobStoreName": "default",
    "strictContentTypeValidation": true
  },
  "proxy": {
    "remoteUrl": "https://registry.cn-hangzhou.aliyuncs.com",
    "contentMaxAge": 1440,
    "metadataMaxAge": 1440
  },
  "negativeCache": {
    "enabled": true,
    "timeToLive": 1440
  },
  "httpClient": {
    "blocked": false,
    "autoBlock": true,
    "connection": {
      "enableCircularRedirects": false,
      "enableCookies": false,
      "useTrustStore": false
    }
  },
  "docker": {
    "v1Enabled": true,
    "forceBasicAuth": true,
    "httpPort": 1445
  },
  "dockerProxy": {
    "indexType": "HUB",
    "cacheForeignLayers": false,
    "foreignLayerUrlWhitelist": []
  }
}'

# 创建 apt proxy 仓库
curl -X 'POST' -u "admin:$NEUXS_PASSWORD" \
  '{{ $.artifactRepoViewUrl }}/service/rest/v1/repositories/apt/proxy' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "name": "apt-aliyun",
  "online": true,
  "storage": {
    "blobStoreName": "default",
    "strictContentTypeValidation": true
  },
  "proxy": {
    "remoteUrl": "https://mirrors.aliyun.com/",
    "contentMaxAge": 1440,
    "metadataMaxAge": 1440
  },
  "negativeCache": {
    "enabled": true,
    "timeToLive": 1440
  },
  "httpClient": {
    "blocked": false,
    "autoBlock": true,
    "connection": {
      "retries": 0,
      "timeout": 60,
      "enableCircularRedirects": false,
      "enableCookies": false,
      "useTrustStore": false
    }
  },
  "apt": {
    "distribution": "bullseye",
    "flat": false
  }
}'

# 创建 go proxy 仓库
curl -X 'POST' -u "admin:$NEUXS_PASSWORD" \
  '{{ $.artifactRepoViewUrl }}/service/rest/v1/repositories/go/proxy' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "name": "go-aliyun",
  "online": true,
  "storage": {
    "blobStoreName": "default",
    "strictContentTypeValidation": true
  },
  "proxy": {
    "remoteUrl": "https://mirrors.aliyun.com/goproxy",
    "contentMaxAge": 1440,
    "metadataMaxAge": 1440
  },
  "negativeCache": {
    "enabled": true,
    "timeToLive": 1440
  },
  "httpClient": {
    "blocked": false,
    "autoBlock": true,
    "connection": {
      "retries": 0,
      "timeout": 60,
      "enableCircularRedirects": false,
      "enableCookies": false,
      "useTrustStore": false
    }
  }
}'

curl -X 'POST' -u "admin:$NEUXS_PASSWORD" \
  '{{ $.artifactRepoViewUrl }}/service/rest/v1/repositories/go/proxy' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "name": "go-goproxy",
  "online": true,
  "storage": {
    "blobStoreName": "default",
    "strictContentTypeValidation": true
  },
  "proxy": {
    "remoteUrl": "https://goproxy.cn",
    "contentMaxAge": 1440,
    "metadataMaxAge": 1440
  },
  "negativeCache": {
    "enabled": true,
    "timeToLive": 1440
  },
  "httpClient": {
    "blocked": false,
    "autoBlock": true,
    "connection": {
      "retries": 0,
      "timeout": 60,
      "enableCircularRedirects": false,
      "enableCookies": false,
      "useTrustStore": false
    }
  }
}'

curl -X 'POST' -u "admin:$NEUXS_PASSWORD" \
  '{{ $.artifactRepoViewUrl }}/service/rest/v1/repositories/go/group' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "name": "go-group-public",
  "online": true,
  "storage": {
    "blobStoreName": "default",
    "strictContentTypeValidation": true
  },
  "group": {
    "memberNames": [
      "go-goproxy",
      "go-aliyun"
    ]
  }
}'

# 创建 maven proxy 仓库
curl -X 'POST' -u "admin:$NEUXS_PASSWORD" \
  '{{ $.artifactRepoViewUrl }}/service/rest/v1/repositories/maven/proxy' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "name": "maven-aliyun-central",
  "online": true,
  "storage": {
    "blobStoreName": "default",
    "strictContentTypeValidation": true
  },
  "proxy": {
    "remoteUrl": "https://maven.aliyun.com/repository/central",
    "contentMaxAge": 1440,
    "metadataMaxAge": 1440
  },
  "negativeCache": {
    "enabled": true,
    "timeToLive": 1440
  },
  "httpClient": {
    "blocked": false,
    "autoBlock": true,
    "connection": {
      "retries": 0,
      "timeout": 60,
      "enableCircularRedirects": false,
      "enableCookies": false,
      "useTrustStore": false
    }
  },
  "maven": {
    "versionPolicy": "RELEASE",
    "layoutPolicy": "PERMISSIVE",
    "contentDisposition": "INLINE"
  }
}'

curl -X 'POST' -u "admin:$NEUXS_PASSWORD" \
  '{{ $.artifactRepoViewUrl }}/service/rest/v1/repositories/maven/proxy' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "name": "maven-aliyun-google",
  "online": true,
  "storage": {
    "blobStoreName": "default",
    "strictContentTypeValidation": true
  },
  "proxy": {
    "remoteUrl": "https://maven.aliyun.com/repository/google",
    "contentMaxAge": 1440,
    "metadataMaxAge": 1440
  },
  "negativeCache": {
    "enabled": true,
    "timeToLive": 1440
  },
  "httpClient": {
    "blocked": false,
    "autoBlock": true,
    "connection": {
      "retries": 0,
      "timeout": 60,
      "enableCircularRedirects": false,
      "enableCookies": false,
      "useTrustStore": false
    }
  },
  "maven": {
    "versionPolicy": "RELEASE",
    "layoutPolicy": "PERMISSIVE",
    "contentDisposition": "INLINE"
  }
}'

curl -X 'POST' -u "admin:$NEUXS_PASSWORD" \
  '{{ $.artifactRepoViewUrl }}/service/rest/v1/repositories/maven/proxy' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "name": "maven-aliyun-group",
  "online": true,
  "storage": {
    "blobStoreName": "default",
    "strictContentTypeValidation": true
  },
  "proxy": {
    "remoteUrl": "https://maven.aliyun.com/nexus/content/groups/public",
    "contentMaxAge": 1440,
    "metadataMaxAge": 1440
  },
  "negativeCache": {
    "enabled": true,
    "timeToLive": 1440
  },
  "httpClient": {
    "blocked": false,
    "autoBlock": true,
    "connection": {
      "retries": 0,
      "timeout": 60,
      "enableCircularRedirects": false,
      "enableCookies": false,
      "useTrustStore": false
    }
  },
  "maven": {
    "versionPolicy": "RELEASE",
    "layoutPolicy": "PERMISSIVE",
    "contentDisposition": "INLINE"
  }
}'

curl -X 'POST' -u "admin:$NEUXS_PASSWORD" \
  '{{ $.artifactRepoViewUrl }}/service/rest/v1/repositories/maven/proxy' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "name": "maven-aliyun-jcenter",
  "online": true,
  "storage": {
    "blobStoreName": "default",
    "strictContentTypeValidation": true
  },
  "proxy": {
    "remoteUrl": "https://maven.aliyun.com/repository/jcenter",
    "contentMaxAge": 1440,
    "metadataMaxAge": 1440
  },
  "negativeCache": {
    "enabled": true,
    "timeToLive": 1440
  },
  "httpClient": {
    "blocked": false,
    "autoBlock": true,
    "connection": {
      "retries": 0,
      "timeout": 60,
      "enableCircularRedirects": false,
      "enableCookies": false,
      "useTrustStore": false
    }
  },
  "maven": {
    "versionPolicy": "RELEASE",
    "layoutPolicy": "PERMISSIVE",
    "contentDisposition": "INLINE"
  }
}'

curl -X 'POST' -u "admin:$NEUXS_PASSWORD" \
  '{{ $.artifactRepoViewUrl }}/service/rest/v1/repositories/maven/proxy' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "name": "maven-aliyun-public",
  "online": true,
  "storage": {
    "blobStoreName": "default",
    "strictContentTypeValidation": true
  },
  "proxy": {
    "remoteUrl": "https://maven.aliyun.com/repository/public",
    "contentMaxAge": 1440,
    "metadataMaxAge": 1440
  },
  "negativeCache": {
    "enabled": true,
    "timeToLive": 1440
  },
  "httpClient": {
    "blocked": false,
    "autoBlock": true,
    "connection": {
      "retries": 0,
      "timeout": 60,
      "enableCircularRedirects": false,
      "enableCookies": false,
      "useTrustStore": false
    }
  },
  "maven": {
    "versionPolicy": "RELEASE",
    "layoutPolicy": "PERMISSIVE",
    "contentDisposition": "INLINE"
  }
}'

curl -X 'POST' -u "admin:$NEUXS_PASSWORD" \
  '{{ $.artifactRepoViewUrl }}/service/rest/v1/repositories/maven/group' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "name": "maven-group-public",
  "online": true,
  "storage": {
    "blobStoreName": "default",
    "strictContentTypeValidation": true
  },
  "group": {
    "memberNames": [
      "maven-aliyun-public",
      "maven-aliyun-group",
      "maven-aliyun-central",
      "maven-aliyun-google",
      "maven-aliyun-jcenter",
      "maven-central"
    ]
  }
}'

# 创建 npm proxy 仓库
curl -X 'POST' -u "admin:$NEUXS_PASSWORD" \
  '{{ $.artifactRepoViewUrl }}/service/rest/v1/repositories/npm/proxy' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "name": "npm-cnpmjs",
  "online": true,
  "storage": {
    "blobStoreName": "default",
    "strictContentTypeValidation": true
  },
  "proxy": {
    "remoteUrl": "https://r.cnpmjs.org",
    "contentMaxAge": 1440,
    "metadataMaxAge": 1440
  },
  "negativeCache": {
    "enabled": true,
    "timeToLive": 1440
  },
  "httpClient": {
    "blocked": false,
    "autoBlock": true,
    "connection": {
      "retries": 0,
      "timeout": 60,
      "enableCircularRedirects": false,
      "enableCookies": false,
      "useTrustStore": false
    }
  },
  "npm": {
    "removeNonCataloged": false,
    "removeQuarantined": false
  }
}'

curl -X 'POST' -u "admin:$NEUXS_PASSWORD" \
  '{{ $.artifactRepoViewUrl }}/service/rest/v1/repositories/npm/proxy' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "name": "npm-taobao",
  "online": true,
  "storage": {
    "blobStoreName": "default",
    "strictContentTypeValidation": true
  },
  "proxy": {
    "remoteUrl": "https://registry.npmmirror.com",
    "contentMaxAge": 1440,
    "metadataMaxAge": 1440
  },
  "negativeCache": {
    "enabled": true,
    "timeToLive": 1440
  },
  "httpClient": {
    "blocked": false,
    "autoBlock": true,
    "connection": {
      "retries": 0,
      "timeout": 60,
      "enableCircularRedirects": false,
      "enableCookies": false,
      "useTrustStore": false
    }
  },
  "npm": {
    "removeNonCataloged": false,
    "removeQuarantined": false
  }
}'

curl -X 'POST' -u "admin:$NEUXS_PASSWORD" \
  '{{ $.artifactRepoViewUrl }}/service/rest/v1/repositories/npm/group' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "name": "npm-group-public",
  "online": true,
  "storage": {
    "blobStoreName": "default",
    "strictContentTypeValidation": true
  },
  "group": {
    "memberNames": [
      "npm-taobao",
      "npm-cnpmjs"
    ]
  }
}'

# 创建 pypi proxy 仓库
curl -X 'POST' -u "admin:$NEUXS_PASSWORD" \
  '{{ $.artifactRepoViewUrl }}/service/rest/v1/repositories/pypi/proxy' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "name": "pypi-aliyun",
  "online": true,
  "storage": {
    "blobStoreName": "default",
    "strictContentTypeValidation": true
  },
  "proxy": {
    "remoteUrl": "http://mirrors.aliyun.com/pypi",
    "contentMaxAge": 1440,
    "metadataMaxAge": 1440
  },
  "negativeCache": {
    "enabled": true,
    "timeToLive": 1440
  },
  "httpClient": {
    "blocked": false,
    "autoBlock": true,
    "connection": {
      "retries": 0,
      "timeout": 60,
      "enableCircularRedirects": false,
      "enableCookies": false,
      "useTrustStore": false
    }
  },
  "pypi": {
    "removeQuarantined": false
  }
}'

curl -X 'POST' -u "admin:$NEUXS_PASSWORD" \
  '{{ $.artifactRepoViewUrl }}/service/rest/v1/repositories/pypi/group' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "name": "pypi-group-public",
  "online": true,
  "storage": {
    "blobStoreName": "default",
    "strictContentTypeValidation": true
  },
  "group": {
    "memberNames": [
      "pypi-aliyun"
    ]
  }
}'

# 创建角色
curl -X 'POST' -u "admin:$NEUXS_PASSWORD" \
  '{{ $.artifactRepoViewUrl }}/service/rest/v1/security/roles' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "id": "anonymous-role",
  "name": "anonymous-role",
  "description": "anonymous-role",
  "privileges": [
    "nx-repository-view-go-go-group-public-read",
    "nx-repository-view-go-go-group-public-browse"
  ]
}'

curl -X 'POST' -u "admin:$NEUXS_PASSWORD" \
  '{{ $.artifactRepoViewUrl }}/service/rest/v1/security/roles' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "id": "public-role",
  "name": "public-role",
  "description": "public-role",
  "privileges": [
    "nx-repository-view-docker-*-browse",
    "nx-repository-view-npm-npm-group-public-browse",
    "nx-repository-view-pypi-pypi-group-public-browse",
    "nx-repository-view-pypi-pypi-group-public-read",
    "nx-repository-view-maven2-maven-group-public-browse",
    "nx-repository-view-go-go-group-public-browse",
    "nx-repository-view-docker-*-read",
    "nx-repository-view-maven2-maven-group-public-read",
    "nx-repository-view-yum-*-browse",
    "nx-repository-view-apt-*-browse",
    "nx-repository-view-go-go-group-public-read",
    "nx-repository-view-apt-*-read",
    "nx-repository-view-npm-npm-group-public-read",
    "nx-repository-view-yum-*-read"
  ]
}'

# 更新匿名用户角色
curl -X 'PUT' -u "admin:$NEUXS_PASSWORD" \
  '{{ $.artifactRepoViewUrl }}/service/rest/v1/security/users/anonymous' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
    "userId": "anonymous",
    "firstName": "Anonymous",
    "lastName": "User",
    "emailAddress": "anonymous@example.org",
    "source": "default",
    "status": "active",
    "readOnly": false,
    "roles": [
      "anonymous-role"
    ]
}'

# 创建 public-user 用户
curl -X 'POST' -u "admin:$NEUXS_PASSWORD" \
  '{{ $.artifactRepoViewUrl }}/service/rest/v1/security/users' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "userId": "public-user",
  "firstName": "public-user",
  "lastName": "public-user",
  "emailAddress": "public-user@example.com",
  "password": "public-user",
  "status": "active",
  "roles": [
    "public-role"
  ]
}'
