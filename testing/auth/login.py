#!/usr/bin/env python3
"""
Login API Test Script

This script tests the authentication flow for the Go DDD Template application:
1. Get captcha (using dev mode for simplicity)
2. Login with admin credentials
3. Verify the access token by calling /api/auth/me

Usage:
    uv run testing/auth/login.py
"""

import requests
import json
import sys
from typing import Dict, Any


class Colors:
    """ANSI color codes for terminal output"""
    HEADER = '\033[95m'
    OKBLUE = '\033[94m'
    OKCYAN = '\033[96m'
    OKGREEN = '\033[92m'
    WARNING = '\033[93m'
    FAIL = '\033[91m'
    ENDC = '\033[0m'
    BOLD = '\033[1m'
    UNDERLINE = '\033[4m'


class APIClient:
    """Simple API client for testing"""

    def __init__(self, base_url: str = "http://localhost:40012"):
        self.base_url = base_url
        self.session = requests.Session()

    def get_captcha(self, dev_mode: bool = True) -> Dict[str, Any]:
        """Get captcha (using dev mode by default for testing)"""
        url = f"{self.base_url}/api/auth/captcha"
        if dev_mode:
            url += "?code=9999&secret=dev-secret-change-me"

        response = self.session.get(url)
        response.raise_for_status()
        return response.json()

    def login(self, username: str, password: str, captcha_id: str, captcha: str) -> Dict[str, Any]:
        """Login with credentials"""
        url = f"{self.base_url}/api/auth/login"
        payload = {
            "login": username,
            "password": password,
            "captcha_id": captcha_id,
            "captcha": captcha
        }

        response = self.session.post(url, json=payload)
        response.raise_for_status()
        return response.json()

    def get_current_user(self, access_token: str) -> Dict[str, Any]:
        """Get current user info using access token"""
        url = f"{self.base_url}/api/auth/me"
        headers = {"Authorization": f"Bearer {access_token}"}

        response = self.session.get(url, headers=headers)
        response.raise_for_status()
        return response.json()


def print_test_header(title: str):
    """Print a formatted test header"""
    print(f"\n{Colors.BOLD}{Colors.HEADER}{'='*60}{Colors.ENDC}")
    print(f"{Colors.BOLD}{Colors.HEADER}{title:^60}{Colors.ENDC}")
    print(f"{Colors.BOLD}{Colors.HEADER}{'='*60}{Colors.ENDC}\n")


def print_success(message: str):
    """Print success message"""
    print(f"{Colors.OKGREEN}✓ {message}{Colors.ENDC}")


def print_error(message: str):
    """Print error message"""
    print(f"{Colors.FAIL}✗ {message}{Colors.ENDC}")


def print_info(message: str):
    """Print info message"""
    print(f"{Colors.OKCYAN}ℹ {message}{Colors.ENDC}")


def print_json(data: Dict[str, Any], indent: int = 2):
    """Print JSON data with formatting"""
    print(json.dumps(data, indent=indent, ensure_ascii=False))


def main():
    """Main test function"""
    print_test_header("Login API Test")

    # Initialize client
    client = APIClient()

    # Test credentials
    username = "admin"
    password = "admin123"

    try:
        # Step 1: Get Captcha
        print_info("Step 1: Getting captcha (dev mode)...")
        captcha_response = client.get_captcha(dev_mode=True)

        if "data" not in captcha_response:
            print_error("Failed to get captcha")
            sys.exit(1)

        captcha_data = captcha_response["data"]
        captcha_id = captcha_data["id"]
        captcha_code = captcha_data.get("code", "9999")  # In dev mode, code is returned

        print_success(f"Captcha obtained: {captcha_id}")
        print_info(f"  Captcha Code: {captcha_code}")
        print_info(f"  Expires At: {captcha_data.get('expire_at')}")

        # Step 2: Login
        print_info(f"\nStep 2: Logging in as '{username}'...")
        login_response = client.login(username, password, captcha_id, captcha_code)

        if "data" not in login_response:
            print_error("Login failed")
            sys.exit(1)

        login_data = login_response["data"]
        access_token = login_data.get("access_token")

        print_success("Login successful!")
        print_info(f"  User ID: {login_data.get('user', {}).get('user_id')}")
        print_info(f"  Username: {login_data.get('user', {}).get('username')}")
        print_info(f"  Token Type: {login_data.get('token_type')}")
        print_info(f"  Access Token: {access_token[:50]}..." if access_token else "  Access Token: None")

        # Step 3: Verify Token
        print_info("\nStep 3: Verifying access token...")
        user_response = client.get_current_user(access_token)

        if "data" not in user_response:
            print_error("Failed to get user info")
            sys.exit(1)

        user_data = user_response["data"]
        print_success("Token verified successfully!")
        print_info(f"  ID: {user_data.get('id')}")
        print_info(f"  Username: {user_data.get('username')}")
        print_info(f"  Email: {user_data.get('email')}")
        print_info(f"  Full Name: {user_data.get('full_name')}")
        print_info(f"  Status: {user_data.get('status')}")
        print_info(f"  Roles: {', '.join([r['name'] for r in user_data.get('roles', [])])}")

        # Summary
        print_test_header("Test Summary")
        print_success("All tests passed!")
        print_info("\nTest Results:")
        print_info("  ✓ Captcha generation")
        print_info("  ✓ User login")
        print_info("  ✓ Token verification")
        print_info("  ✓ User profile retrieval")

        print(f"\n{Colors.OKGREEN}{Colors.BOLD}SUCCESS: Login API is working correctly!{Colors.ENDC}\n")

        return 0

    except requests.exceptions.ConnectionError:
        print_error("Connection failed: Is the server running on http://localhost:40012?")
        return 1
    except requests.exceptions.HTTPError as e:
        print_error(f"HTTP Error: {e}")
        if e.response is not None:
            print_info("Response:")
            try:
                print_json(e.response.json())
            except:
                print(e.response.text)
        return 1
    except Exception as e:
        print_error(f"Unexpected error: {e}")
        import traceback
        traceback.print_exc()
        return 1


if __name__ == "__main__":
    sys.exit(main())
