#!/usr/bin/env python3
"""
用户管理 API 测试脚本
测试 /api/admin/users 端点的 CRUD 操作
"""

import os
import requests
from dotenv import load_dotenv

# 加载环境变量
load_dotenv()

# API 基础 URL
BASE_URL = os.getenv("API_BASE_URL", "http://localhost:40012")

# 全局变量存储 token 和测试数据
access_token = None
test_user_id = None


def print_separator(title=""):
    """打印分隔线"""
    if title:
        print(f"\n{'='*60}")
        print(f"  {title}")
        print(f"{'='*60}")
    else:
        print(f"{'='*60}\n")


def get_captcha():
    """获取验证码 (开发模式)"""
    dev_code = "9999"
    dev_secret = os.getenv("DEV_SECRET", "dev-secret-change-me")
    url = f"{BASE_URL}/api/auth/captcha?code={dev_code}&secret={dev_secret}"

    try:
        response = requests.get(url, timeout=10)
        if response.status_code == 200:
            data = response.json()
            if "data" in data:
                return data["data"].get("id"), data["data"].get("code")
    except Exception as e:
        print(f"❌ 获取验证码失败: {e}")
    return None, None


def login():
    """登录获取 access token"""
    global access_token

    print_separator("登录获取 Token")

    captcha_id, captcha_code = get_captcha()
    if not captcha_id or not captcha_code:
        print("❌ 无法获取验证码")
        return False

    url = f"{BASE_URL}/api/auth/login"
    payload = {
        "account": "admin",
        "password": "admin123",
        "captcha_id": captcha_id,
        "captcha": captcha_code
    }

    try:
        response = requests.post(url, json=payload, timeout=10)
        print(f"状态码: {response.status_code}")

        if response.status_code == 200:
            data = response.json()
            if "data" in data and "access_token" in data["data"]:
                access_token = data["data"]["access_token"]
                print(f"✅ 登录成功!")
                print(f"Token: {access_token[:50]}...")
                return True

        print(f"❌ 登录失败: {response.text}")
        return False

    except Exception as e:
        print(f"❌ 登录异常: {e}")
        return False


def get_headers():
    """获取带 token 的请求头"""
    return {
        "Authorization": f"Bearer {access_token}",
        "Content-Type": "application/json"
    }


def test_list_users():
    """测试获取用户列表"""
    print_separator("测试 1: 获取用户列表")

    url = f"{BASE_URL}/api/admin/users"
    params = {"page": 1, "limit": 10}

    try:
        response = requests.get(url, headers=get_headers(), params=params, timeout=10)
        print(f"状态码: {response.status_code}")

        if response.status_code == 200:
            data = response.json()
            print(f"✅ 获取成功!")
            print(f"响应数据: {data}")

            users = data.get("data", [])
            pagination = data.get("meta", {})
            print(f"\n用户数量: {len(users)}")
            print(f"总数: {pagination.get('total', 0)}")

            if users:
                print("\n前3个用户:")
                for user in users[:3]:
                    print(f"  - ID: {user.get('id')}, 用户名: {user.get('username')}, 邮箱: {user.get('email')}")
            return True
        else:
            print(f"❌ 请求失败: {response.text}")
            return False

    except Exception as e:
        print(f"❌ 请求异常: {e}")
        return False


def test_create_user():
    """测试创建用户"""
    global test_user_id

    print_separator("测试 2: 创建用户")

    url = f"{BASE_URL}/api/admin/users"
    payload = {
        "username": "testuser123",
        "email": "testuser123@example.com",
        "password": "password123",
        "full_name": "测试用户",
        "status": "active"
    }

    print(f"创建用户: {payload['username']}")

    try:
        response = requests.post(url, headers=get_headers(), json=payload, timeout=10)
        print(f"状态码: {response.status_code}")

        if response.status_code == 200 or response.status_code == 201:
            data = response.json()
            print(f"✅ 创建成功!")
            print(f"响应数据: {data}")

            if "data" in data:
                test_user_id = data["data"].get("id")
                print(f"\n新用户ID: {test_user_id}")
                print(f"用户名: {data['data'].get('username')}")
                print(f"邮箱: {data['data'].get('email')}")
            return True
        else:
            print(f"❌ 创建失败: {response.text}")
            return False

    except Exception as e:
        print(f"❌ 请求异常: {e}")
        return False


def test_get_user():
    """测试获取单个用户详情"""
    print_separator("测试 3: 获取用户详情")

    if not test_user_id:
        print("❌ 没有测试用户ID，跳过此测试")
        return False

    url = f"{BASE_URL}/api/admin/users/{test_user_id}"

    try:
        response = requests.get(url, headers=get_headers(), timeout=10)
        print(f"状态码: {response.status_code}")

        if response.status_code == 200:
            data = response.json()
            print(f"✅ 获取成功!")
            print(f"响应数据: {data}")

            if "data" in data:
                user = data["data"]
                print(f"\n用户详情:")
                print(f"  ID: {user.get('id')}")
                print(f"  用户名: {user.get('username')}")
                print(f"  邮箱: {user.get('email')}")
                print(f"  状态: {user.get('status')}")
                print(f"  角色: {user.get('roles', [])}")
            return True
        else:
            print(f"❌ 获取失败: {response.text}")
            return False

    except Exception as e:
        print(f"❌ 请求异常: {e}")
        return False


def test_update_user():
    """测试更新用户"""
    print_separator("测试 4: 更新用户")

    if not test_user_id:
        print("❌ 没有测试用户ID，跳过此测试")
        return False

    url = f"{BASE_URL}/api/admin/users/{test_user_id}"
    payload = {
        "full_name": "测试用户（已更新）",
        "status": "active"
    }

    print(f"更新用户ID: {test_user_id}")
    print(f"更新数据: {payload}")

    try:
        response = requests.put(url, headers=get_headers(), json=payload, timeout=10)
        print(f"状态码: {response.status_code}")

        if response.status_code == 200:
            data = response.json()
            print(f"✅ 更新成功!")
            print(f"响应数据: {data}")

            if "data" in data:
                user = data["data"]
                print(f"\n更新后的用户信息:")
                print(f"  全名: {user.get('full_name')}")
                print(f"  状态: {user.get('status')}")
            return True
        else:
            print(f"❌ 更新失败: {response.text}")
            return False

    except Exception as e:
        print(f"❌ 请求异常: {e}")
        return False


def test_search_users():
    """测试搜索用户"""
    print_separator("测试 5: 搜索用户")

    url = f"{BASE_URL}/api/admin/users"
    params = {
        "page": 1,
        "limit": 10,
        "search": "testuser"
    }

    print(f"搜索关键词: {params['search']}")

    try:
        response = requests.get(url, headers=get_headers(), params=params, timeout=10)
        print(f"状态码: {response.status_code}")

        if response.status_code == 200:
            data = response.json()
            print(f"✅ 搜索成功!")

            users = data.get("data", [])
            print(f"\n找到 {len(users)} 个用户:")
            for user in users:
                print(f"  - 用户名: {user.get('username')}, 邮箱: {user.get('email')}")
            return True
        else:
            print(f"❌ 搜索失败: {response.text}")
            return False

    except Exception as e:
        print(f"❌ 请求异常: {e}")
        return False


def test_assign_roles():
    """测试分配角色"""
    print_separator("测试 6: 分配角色")

    if not test_user_id:
        print("❌ 没有测试用户ID，跳过此测试")
        return False

    # 先获取可用角色列表
    roles_url = f"{BASE_URL}/api/admin/roles"
    try:
        roles_response = requests.get(roles_url, headers=get_headers(), params={"page": 1, "limit": 10}, timeout=10)
        if roles_response.status_code == 200:
            roles_data = roles_response.json()
            roles = roles_data.get("data", [])
            if roles:
                # 取第一个角色进行测试
                role_id = roles[0]["id"]
                print(f"使用角色ID: {role_id} ({roles[0].get('display_name', 'N/A')})")

                # 分配角色
                url = f"{BASE_URL}/api/admin/users/{test_user_id}/roles"
                payload = {"role_ids": [role_id]}

                response = requests.put(url, headers=get_headers(), json=payload, timeout=10)
                print(f"状态码: {response.status_code}")

                if response.status_code == 200:
                    data = response.json()
                    print(f"✅ 角色分配成功!")
                    print(f"响应数据: {data}")

                    if "data" in data:
                        user = data["data"]
                        print(f"\n用户角色: {user.get('roles', [])}")
                    return True
                else:
                    print(f"❌ 分配失败: {response.text}")
                    return False
            else:
                print("❌ 没有可用角色")
                return False
        else:
            print(f"❌ 获取角色列表失败")
            return False

    except Exception as e:
        print(f"❌ 请求异常: {e}")
        return False


def test_delete_user():
    """测试删除用户"""
    print_separator("测试 7: 删除用户")

    if not test_user_id:
        print("❌ 没有测试用户ID，跳过此测试")
        return False

    url = f"{BASE_URL}/api/admin/users/{test_user_id}"

    print(f"删除用户ID: {test_user_id}")

    try:
        response = requests.delete(url, headers=get_headers(), timeout=10)
        print(f"状态码: {response.status_code}")

        if response.status_code == 200 or response.status_code == 204:
            print(f"✅ 删除成功!")

            # 验证用户已删除
            verify_response = requests.get(url, headers=get_headers(), timeout=10)
            if verify_response.status_code == 404:
                print(f"✅ 验证成功: 用户已被删除")
            else:
                print(f"⚠️  用户可能仍然存在 (状态码: {verify_response.status_code})")
            return True
        else:
            print(f"❌ 删除失败: {response.text}")
            return False

    except Exception as e:
        print(f"❌ 请求异常: {e}")
        return False


def main():
    """主测试流程"""
    print_separator("🧪 用户管理 API 测试套件")
    print(f"API Base URL: {BASE_URL}")

    # 步骤 0: 登录
    if not login():
        print("\n❌ 登录失败，终止测试")
        return

    # 测试统计
    tests = [
        ("获取用户列表", test_list_users),
        ("创建用户", test_create_user),
        ("获取用户详情", test_get_user),
        ("更新用户", test_update_user),
        ("搜索用户", test_search_users),
        ("分配角色", test_assign_roles),
        ("删除用户", test_delete_user),
    ]

    passed = 0
    failed = 0

    for test_name, test_func in tests:
        try:
            result = test_func()
            if result:
                passed += 1
            else:
                failed += 1
        except Exception as e:
            print(f"❌ 测试异常: {e}")
            failed += 1

    # 总结
    print_separator("📊 测试总结")
    print(f"总测试数: {len(tests)}")
    print(f"✅ 通过: {passed}")
    print(f"❌ 失败: {failed}")
    print(f"成功率: {passed / len(tests) * 100:.1f}%")
    print_separator()


if __name__ == "__main__":
    main()
