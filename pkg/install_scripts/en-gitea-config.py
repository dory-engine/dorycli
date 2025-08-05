from html.parser import HTMLParser
from urllib import request, parse
import time
import base64
import json


class InstallFormParser(HTMLParser):
    def __init__(self):
        super().__init__()
        self.in_form = False
        self.form_action = '/'
        self.form_method = 'post'
        self.form_data = {}

    def handle_starttag(self, tag, attrs):
        attrs = dict(attrs)

        if tag == 'form':
            self.in_form = True
            self.form_action = attrs.get('action', '/')
            self.form_method = attrs.get('method', 'post').lower()

        if self.in_form and tag == 'input':
            name = attrs.get('name')
            if not name:
                return
            input_type = attrs.get('type', 'text')

            if input_type == 'checkbox':
                if 'checked' in attrs:
                    self.form_data[name] = attrs.get('value', 'on')
            else:
                self.form_data[name] = attrs.get('value', '')

    def handle_endtag(self, tag):
        if tag == 'form':
            self.in_form = False


class LoginFormChecker(HTMLParser):
    def __init__(self):
        super().__init__()
        self.found_user_input = False

    def handle_starttag(self, tag, attrs):
        attrs = dict(attrs)
        if tag == 'input' and attrs.get('name') == 'user_name':
            self.found_user_input = True


# config value
base_url = '{{ $.gitRepoViewUrl }}'
username = 'root'
password = '{{ $.gitRepoPassword }}'
mail = 'root@example.com'

# === step 1: request install page ===
try:
    print(f"📥 get Gitea install page: {base_url}")
    with request.urlopen(base_url) as resp:
        html = resp.read().decode('utf-8')
except Exception as e:
    print("❌ can not access install page:", e)
    exit(1)

# === step 2: parse default form data ===
parser = InstallFormParser()
parser.feed(html)

print("📋 get default install form data:")
parser.form_data['app_url'] = f'{base_url}/'
parser.form_data['admin_name'] = username
parser.form_data['admin_email'] = mail
parser.form_data['admin_passwd'] = password
parser.form_data['admin_confirm_passwd'] = password
for k, v in parser.form_data.items():
    print(f"{k} = {v}")

# === step 3: post install form ===
post_url = base_url + parser.form_action
post_data = parse.urlencode(parser.form_data).encode('utf-8')
headers = {'Content-Type': 'application/x-www-form-urlencoded'}

print(f"\n🚀 post form url: {post_url}")
try:
    req = request.Request(post_url, data=post_data, headers=headers, method='POST')
    with request.urlopen(req) as resp:
        body = resp.read().decode('utf-8')
        print(f"✅ post form data success, status code: {resp.status}")
        print("📄 body:", body[:300])
except Exception as e:
    print("❌ post form data fail:", e)
    exit(1)

# === step 4: query login page every 5 seconds ===
login_url = f'{base_url}/user/login'
print("\n⏳ check Gitea redirect to login page...")

for attempt in range(0, 30):
    try:
        with request.urlopen(login_url) as resp:
            login_html = resp.read().decode('utf-8')

        login_parser = LoginFormChecker()
        login_parser.feed(login_html)

        if login_parser.found_user_input:
            print(f"\n🎉 get login page ({attempt} try), Gitea install success!")
            break
        else:
            print(f"{attempt} try: login page not ready, waiting...")
    except Exception as e:
        print(f"{attempt} try check fail: {e}")
    time.sleep(5)
else:
    print("\n❌ timeout but login page not available, install fail")

auth_str = f'{username}:{password}'
auth_encoded = base64.b64encode(auth_str.encode()).decode()

url = f'{base_url}/api/v1/users/{username}/tokens'
headers = {
    'Authorization': f'Basic {auth_encoded}',
    'Accept': 'application/json',
    'Content-Type': 'application/json',
}

payload = {
    "name": "dory",
    "scopes": [
        "all",
        "write:activitypub",
        "write:admin",
        "write:issue",
        "write:misc",
        "write:notification",
        "write:organization",
        "write:package",
        "write:repository",
        "write:user"
    ]
}

# === request access token ===
data_bytes = json.dumps(payload).encode('utf-8')
req = request.Request(url, data=data_bytes, headers=headers, method='POST')
token_sha1 = ""
try:
    with request.urlopen(req) as resp:
        resp_data = resp.read().decode('utf-8')
        resp_json = json.loads(resp_data)
        token_sha1 = resp_json.get('sha1')

        if token_sha1:
            print(f"✅ Access Token (sha1): {token_sha1}")
        else:
            print("❌ request without sha1, full response is: ")
            print(resp_json)

except Exception as e:
    print("❌ request fail:", e)


def replace_config_file(file_path, str1, str2):
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()

        updated_content = content.replace(str1, str2)

        if content == updated_content:
            print(f"⚠️ {str1} not found")
        else:
            with open(file_path, 'w', encoding='utf-8') as f:
                f.write(updated_content)
            print(f"✅ {str1} replace success")

    except Exception as e:
        print("❌ replace fail: ", e)


replace_config_file('/src/install-data/{{ $.dory.namespace }}/dory-engine/config/config.yaml', 'GIT_REPO_USERNAME', username)
replace_config_file('/src/install-data/{{ $.dory.namespace }}/dory-engine/config/config.yaml', 'GIT_REPO_NAME', username)
replace_config_file('/src/install-data/{{ $.dory.namespace }}/dory-engine/config/config.yaml', 'GIT_REPO_MAIL@example.com', mail)
replace_config_file('/src/install-data/{{ $.dory.namespace }}/dory-engine/config/config.yaml', 'GIT_REPO_TOKEN', token_sha1)
