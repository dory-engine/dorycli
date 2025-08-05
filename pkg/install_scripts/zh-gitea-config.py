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


# 配置项
base_url = '{{ $.gitRepoViewUrl }}'
username = 'root'
password = '{{ $.gitRepoPassword }}'
mail = 'root@example.com'

# === 步骤 1: 请求安装页面 ===
try:
    print(f"📥 获取 Gitea 安装页面: {base_url}")
    with request.urlopen(base_url) as resp:
        html = resp.read().decode('utf-8')
except Exception as e:
    print("❌ 无法访问安装页面:", e)
    exit(1)

# === 步骤 2: 解析默认表单字段 ===
parser = InstallFormParser()
parser.feed(html)

print("📋 获取到默认表单字段：")
parser.form_data['app_url'] = f'{base_url}/'
parser.form_data['admin_name'] = username
parser.form_data['admin_email'] = mail
parser.form_data['admin_passwd'] = password
parser.form_data['admin_confirm_passwd'] = password
for k, v in parser.form_data.items():
    print(f"{k} = {v}")

# === 步骤 3: 提交安装表单 ===
post_url = base_url + parser.form_action
post_data = parse.urlencode(parser.form_data).encode('utf-8')
headers = {'Content-Type': 'application/x-www-form-urlencoded'}

print(f"\n🚀 提交表单到: {post_url}")
try:
    req = request.Request(post_url, data=post_data, headers=headers, method='POST')
    with request.urlopen(req) as resp:
        body = resp.read().decode('utf-8')
        print(f"✅ 表单提交成功，状态码: {resp.status}")
        print("📄 返回页面片段:", body[:300])
except Exception as e:
    print("❌ 表单提交失败:", e)
    exit(1)

# === 步骤 4: 每5秒轮询检查登录页是否出现 ===
login_url = f'{base_url}/user/login'
print("\n⏳ 检查 Gitea 是否跳转到登录页...")

for attempt in range(0, 30):
    try:
        with request.urlopen(login_url) as resp:
            login_html = resp.read().decode('utf-8')

        login_parser = LoginFormChecker()
        login_parser.feed(login_html)

        if login_parser.found_user_input:
            print(f"\n🎉 检测到登录页（第 {attempt} 次尝试），Gitea 安装成功！")
            break
        else:
            print(f"第 {attempt} 次检查：登录页尚未准备好，继续等待...")
    except Exception as e:
        print(f"第 {attempt} 次检查失败：{e}")
    time.sleep(5)
else:
    print("\n❌ 超时未检测到登录页，可能安装失败。")

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

# === 构造请求 ===
data_bytes = json.dumps(payload).encode('utf-8')
req = request.Request(url, data=data_bytes, headers=headers, method='POST')
token_sha1 = ""
try:
    with request.urlopen(req) as resp:
        resp_data = resp.read().decode('utf-8')
        resp_json = json.loads(resp_data)
        token_sha1 = resp_json.get('sha1')

        if token_sha1:
            print(f"✅ Access Token（sha1）: {token_sha1}")
        else:
            print("❌ 响应中未包含 sha1 字段。完整响应：")
            print(resp_json)

except Exception as e:
    print("❌ 请求失败:", e)


def replace_config_file(file_path, str1, str2):
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()

        updated_content = content.replace(str1, str2)

        if content == updated_content:
            print(f"⚠️ 没有找到需要替换的 {str1}")
        else:
            with open(file_path, 'w', encoding='utf-8') as f:
                f.write(updated_content)
            print(f"✅ 已成功替换 {str1}")

    except Exception as e:
        print("❌ 替换失败: ", e)


replace_config_file('/src/install-data/{{ $.dory.namespace }}/dory-engine/config/config.yaml', 'GIT_REPO_USERNAME', username)
replace_config_file('/src/install-data/{{ $.dory.namespace }}/dory-engine/config/config.yaml', 'GIT_REPO_NAME', username)
replace_config_file('/src/install-data/{{ $.dory.namespace }}/dory-engine/config/config.yaml', 'GIT_REPO_MAIL@example.com', mail)
replace_config_file('/src/install-data/{{ $.dory.namespace }}/dory-engine/config/config.yaml', 'GIT_REPO_TOKEN', token_sha1)
