#!/usr/bin/env python3
"""
登录 API 测试脚本
测试 POST /api/v1/auth/login 端点的各种场景
"""

import os
import requests
from dotenv import load_dotenv

# 加载环境变量
load_dotenv()

# API 基础 URL
BASE_URL = os.getenv("API_BASE_URL", "http://localhost:40012")


def print_separator(title=""):
    """打印分隔线"""
    if title:
        print(f"\n{'='*60}")
        print(f"  {title}")
        print(f"{'='*60}")
    else:
        print(f"{'='*60}\n")


def get_captcha():
    """
    获取验证码 (使用开发模式)

    Returns:
        tuple: (captcha_id, captcha_answer) 或 (None, None) 如果失败
    """
    # 使用开发模式获取固定验证码
    dev_code = "9999"
    dev_secret = os.getenv("DEV_SECRET", "dev-secret-change-me")
    url = f"{BASE_URL}/api/auth/captcha?code={dev_code}&secret={dev_secret}"

    try:
        print(f"📡 正在获取验证码 (开发模式): {url}")
        response = requests.get(url, timeout=10)

        print(f"   状态码: {response.status_code}")

        if response.status_code == 200:
            data = response.json()

            # 检查响应格式
            if "data" in data:
                captcha_id = data["data"].get("id")
                captcha_code = data["data"].get("code")  # 开发模式返回 code 字段

                if captcha_id and captcha_code:
                    print(f"   ✅ 获取成功")
                    print(f"   验证码ID: {captcha_id}")
                    print(f"   验证码答案: {captcha_code}")
                    return captcha_id, captcha_code
                else:
                    print(f"   ❌ 响应数据格式不正确: {data}")
            else:
                print(f"   ❌ 响应格式错误: {data}")
        else:
            print(f"   ❌ 请求失败: {response.text}")

    except requests.exceptions.RequestException as e:
        print(f"   ❌ 请求异常: {e}")

    return None, None


def test_login(account, password, captcha_id, captcha, test_name=""):
    """
    测试登录接口

    Args:
        account: 手机号/用户名/邮箱
        password: 密码
        captcha_id: 验证码ID
        captcha: 验证码
        test_name: 测试用例名称

    Returns:
        dict: 响应数据
    """
    url = f"{BASE_URL}/api/auth/login"

    payload = {"account": account, "password": password, "captcha_id": captcha_id, "captcha": captcha}

    print_separator(test_name if test_name else "登录测试")

    print(f"📡 正在测试登录: {url}")
    print(f"   请求参数:")
    print(f"      account: {account}")
    print(f"      password: {'*' * len(password)}")
    print(f"      captcha_id: {captcha_id}")
    print(f"      captcha: {captcha}")

    try:
        response = requests.post(url, json=payload, timeout=10)

        print(f"\n   响应状态码: {response.status_code}")

        try:
            data = response.json()
            print(f"   响应数据: {data}")

            # 分析响应
            if response.status_code == 200:
                if "data" in data:
                    response_data = data["data"]

                    # 检查是否需要 2FA
                    if response_data.get("requires_2fa"):
                        print(f"\n   🔐 需要双因素认证 (2FA)")
                        print(f"   Session Token: {response_data.get('session_token')}")
                    else:
                        # 正常登录成功
                        print(f"\n   ✅ 登录成功!")
                        print(f"   消息: {data.get('message', 'N/A')}")
                        if "access_token" in response_data:
                            print(f"   Access Token: {response_data['access_token'][:50]}...")
                        if "refresh_token" in response_data:
                            print(f"   Refresh Token: {response_data['refresh_token'][:50]}...")
                        if "expires_in" in response_data:
                            print(f"   过期时间: {response_data['expires_in']} 秒")
                        if "user" in response_data:
                            user = response_data["user"]
                            print(f"   用户信息:")
                            print(f"      User ID: {user.get('user_id')}")
                            print(f"      Username: {user.get('username')}")
                else:
                    print(f"\n   ⚠️  响应格式异常")
            elif response.status_code == 401:
                # 新格式：{"code": 401, "message": "..."}
                print(f"\n   ❌ 认证失败: {data.get('message', '未知错误')}")
            else:
                print(f"\n   ❌ 请求失败")

            return data

        except ValueError:
            print(f"   ❌ 响应不是有效的 JSON: {response.text}")
            return None

    except requests.exceptions.RequestException as e:
        print(f"   ❌ 请求异常: {e}")
        return None


def main():
    """主测试流程"""
    print_separator("🧪 登录 API 测试套件")
    print(f"API Base URL: {BASE_URL}")

    # 测试 1: 使用默认管理员账户登录 (正常流程)
    print_separator("测试 1: 默认管理员账户登录")
    captcha_id, captcha_answer = get_captcha()

    if captcha_id and captcha_answer:
        test_login(account="admin", password="admin123", captcha_id=captcha_id, captcha=captcha_answer, test_name="测试 1: 使用正确凭证登录")
    else:
        print("❌ 无法获取验证码,跳过测试 1")

    # 测试 2: 错误的密码
    print_separator("测试 2: 错误密码")
    captcha_id, captcha_answer = get_captcha()

    if captcha_id and captcha_answer:
        test_login(account="admin", password="wrong_password", captcha_id=captcha_id, captcha=captcha_answer, test_name="测试 2: 错误的密码")
    else:
        print("❌ 无法获取验证码,跳过测试 2")

    # 测试 3: 错误的验证码
    print_separator("测试 3: 错误验证码")
    captcha_id, _ = get_captcha()

    if captcha_id:
        test_login(account="admin", password="admin123", captcha_id=captcha_id, captcha="0000", test_name="测试 3: 错误的验证码")  # 错误的验证码
    else:
        print("❌ 无法获取验证码,跳过测试 3")

    # 测试 4: 不存在的用户
    print_separator("测试 4: 不存在的用户")
    captcha_id, captcha_answer = get_captcha()

    if captcha_id and captcha_answer:
        test_login(account="nonexistent_user", password="password123", captcha_id=captcha_id, captcha=captcha_answer, test_name="测试 4: 不存在的用户")
    else:
        print("❌ 无法获取验证码,跳过测试 4")

    print_separator("✅ 测试完成")


if __name__ == "__main__":
    main()
