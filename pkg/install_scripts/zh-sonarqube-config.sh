export PASSWORD={{ $.scanCodeRepoPassword }}

# 修改admin密码
curl -X 'GET' -u admin:admin \
  {{ $.scanCodeRepoViewUrl }}/api/authentication/validate

export TOKEN=$(curl -X 'POST' -u admin:admin \
  {{ $.scanCodeRepoViewUrl }}/api/user_tokens/generate \
  -d 'name=dory' | jq -r .token)
echo "token: $TOKEN"

curl -X 'POST' -u $TOKEN: \
  {{ $.scanCodeRepoViewUrl }}/api/users/change_password \
  -d "login=admin&password=$PASSWORD&previousPassword=admin"

# 创建 admin token
curl -X 'GET' -u admin:$PASSWORD \
  {{ $.scanCodeRepoViewUrl }}/api/authentication/validate

# 设置项目默认不可见
curl -X 'POST' -u $TOKEN: \
  {{ $.scanCodeRepoViewUrl }}/api/projects/update_default_visibility \
  -d "projectVisibility=private"

# 更新dory-engine的配置文件config.yaml
sed -i "s/SCAN_CODE_REPO_TOKEN/$TOKEN/g" install-data/{{ $.dory.namespace }}/dory-engine/config/config.yaml
