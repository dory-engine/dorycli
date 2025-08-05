import urllib.request
import urllib.parse
import http.cookiejar
import datetime
from html.parser import HTMLParser
import json
import argparse


class TokenParser(HTMLParser):
    def __init__(self):
        super().__init__()
        self.token = None
        self.in_input = False

    def handle_starttag(self, tag, attrs):
        attrs = dict(attrs)
        if tag == "meta" and attrs.get("name") == "csrf-token":
            self.token = attrs.get("content")


class PATParser(HTMLParser):
    def __init__(self):
        super().__init__()
        self.token = None

    def handle_starttag(self, tag, attrs):
        attrs = dict(attrs)
        if tag == "input" and attrs.get("id") == "created-personal-access-token":
            self.token = attrs.get("value")


def get_authenticity_token(html):
    parser = TokenParser()
    parser.feed(html)
    return parser.token


def extract_pat_token(html):
    parser = PATParser()
    parser.feed(html)
    return parser.token


def make_request(opener, url, data=None):
    req = urllib.request.Request(url, data=data)
    return opener.open(req)


def change_root_password(gitlab_url, access_token, new_password):
    url = f"{gitlab_url}/api/v4/users/1"
    headers = {
        "Content-Type": "application/json",
        "PRIVATE-TOKEN": access_token
    }
    data = {
        "password": new_password
    }
    request_data = json.dumps(data).encode()

    req = urllib.request.Request(url, data=request_data, headers=headers, method="PUT")
    try:
        with urllib.request.urlopen(req) as response:
            if response.status == 200:
                print(f"✅ update root password success: {new_password}")
            else:
                print(f"⚠️ status code: {response.status}")
    except urllib.error.HTTPError as e:
        print(f"❌ http error: {e.code} - {e.reason}")
        print(e.read().decode())
    except Exception as e:
        print(f"❌ other error: {e}")


def update_application_settings(gitlab_url, access_token):
    url = f"{gitlab_url}/api/v4/application/settings"
    headers = {
        "Content-Type": "application/json",
        "PRIVATE-TOKEN": access_token
    }
    setting_data = {
        "allow_local_requests_from_web_hooks_and_services": True
    }
    data = json.dumps(setting_data).encode()

    req = urllib.request.Request(url, data=data, headers=headers, method="PUT")

    try:
        with urllib.request.urlopen(req) as response:
            resp_body = response.read().decode()
            resp_data = json.loads(resp_body)
            print(f"✅ application settings updated allow_local_requests_from_web_hooks_and_services: {resp_data['allow_local_requests_from_web_hooks_and_services']}")
    except urllib.error.HTTPError as e:
        print(f"❌ http error: {e.code} - {e.reason}")
        print(e.read().decode())
    except Exception as e:
        print(f"❌ other error: {e}")


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
            print(f"✅ {str1} replaced")

    except Exception as e:
        print("❌ replace error: ", e)


def main():
    # 配置项
    gitlab_url = "{{ $.gitRepoViewUrl }}"
    username = "root"
    token_name = "dory-token"
    scopes = ["api"]
    expire_days = 365
    new_password = "{{ $.gitRepoPassword }}"

    parser = argparse.ArgumentParser()
    parser.add_argument('--password', type=str, required=True, help='GitLab password')
    args = parser.parse_args()

    password = args.password

    cj = http.cookiejar.CookieJar()
    opener = urllib.request.build_opener(urllib.request.HTTPCookieProcessor(cj))

    login_url = f"{gitlab_url}/users/sign_in"
    login_page = make_request(opener, login_url)
    login_html = login_page.read().decode()
    auth_token = get_authenticity_token(login_html)
    if not auth_token:
        print("❌ get login page authenticity_token failed")
        return

    login_data = urllib.parse.urlencode({
        "user[login]": username,
        "user[password]": password,
        "authenticity_token": auth_token
    }).encode()
    login_response = make_request(opener, login_url, login_data)
    login_result = login_response.read().decode()
    if "Invalid Login or password" in login_result:
        print("❌ username or password error")
        return

    pat_page_url = f"{gitlab_url}/-/user_settings/personal_access_tokens"
    pat_page = make_request(opener, pat_page_url)
    pat_html = pat_page.read().decode()
    auth_token2 = get_authenticity_token(pat_html)
    if not auth_token2:
        print("❌ get create token page authenticity_token failed")
        return

    expire_at = (datetime.datetime.now() + datetime.timedelta(days=expire_days)).strftime("%Y-%m-%d")
    pat_data = {
        "authenticity_token": auth_token2,
        "personal_access_token[name]": token_name,
        "personal_access_token[expires_at]": expire_at,
    }
    for scope in scopes:
        pat_data["personal_access_token[scopes][]"] = scope

    pat_encoded = urllib.parse.urlencode(pat_data).encode()
    pat_response = make_request(opener, pat_page_url, pat_encoded)
    pat_result = pat_response.read().decode()

    access_token = ""
    try:
        data = json.loads(pat_result)
        access_token = data.get("token") or data.get("new_token")
        if access_token:
            print(f"✅ Access Token created: {access_token}")
        else:
            print("❌ Access Token not found")
    except json.JSONDecodeError:
        print("❌ response not JSON")

    change_root_password(gitlab_url, access_token, new_password)

    update_application_settings(gitlab_url, access_token)

    replace_config_file('/src/install-data/{{ $.dory.namespace }}/dory-engine/config/config.yaml', 'GIT_REPO_USERNAME', 'root')
    replace_config_file('/src/install-data/{{ $.dory.namespace }}/dory-engine/config/config.yaml', 'GIT_REPO_NAME', 'Administrator')
    replace_config_file('/src/install-data/{{ $.dory.namespace }}/dory-engine/config/config.yaml', 'GIT_REPO_MAIL@example.com', 'gitlab_admin@example.com')
    replace_config_file('/src/install-data/{{ $.dory.namespace }}/dory-engine/config/config.yaml', 'GIT_REPO_TOKEN', access_token)


if __name__ == "__main__":
    main()
