export PASSWORD={{ $.scanCodeRepoPassword }}

# change admin password
curl -X 'GET' -u admin:admin \
  {{ $.scanCodeRepoViewUrl }}/api/authentication/validate

export TOKEN=$(curl -X 'POST' -u admin:admin \
  {{ $.scanCodeRepoViewUrl }}/api/user_tokens/generate \
  -d 'name=dory' | jq -r .token)
echo "token: $TOKEN"

curl -X 'POST' -u $TOKEN: \
  {{ $.scanCodeRepoViewUrl }}/api/users/change_password \
  -d "login=admin&password=$PASSWORD&previousPassword=admin"

# create admin token
curl -X 'GET' -u admin:$PASSWORD \
  {{ $.scanCodeRepoViewUrl }}/api/authentication/validate

# set projects invisible by default
curl -X 'POST' -u $TOKEN: \
  {{ $.scanCodeRepoViewUrl }}/api/projects/update_default_visibility \
  -d "projectVisibility=private"

# update dory-engine config file config.yaml
sed -i "s/SCAN_CODE_REPO_TOKEN/$TOKEN/g" install-data/{{ $.dory.namespace }}/dory-engine/config/config.yaml
