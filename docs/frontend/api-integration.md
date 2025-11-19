# API 集成

前端 API 层集中在 `web/src/api/`，与后端模块保持一致：`auth/` 负责认证与用户自助、`user/` 负责个人访问令牌、`admin/` 负责后台管理。以下示例均来自当前代码。

## Axios 客户端（`src/api/auth/client.ts`）

```ts
import axios from "axios";
import { getAccessToken, getRefreshToken, saveAccessToken, saveRefreshToken, clearAuthTokens } from "@/utils/auth";
import type { ApiResponse, AuthResponse } from "@/types/auth";

const API_BASE_URL = "/api/auth";

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
});

apiClient.interceptors.request.use((config) => {
  const token = getAccessToken();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      const refreshToken = getRefreshToken();
      if (refreshToken) {
        const { data } = await axios.post<ApiResponse<AuthResponse>>(`${API_BASE_URL}/refresh`, {
          refresh_token: refreshToken,
        });
        if (data.data) {
          saveAccessToken(data.data.access_token);
          saveRefreshToken(data.data.refresh_token);
          originalRequest.headers.Authorization = `Bearer ${data.data.access_token}`;
          return apiClient(originalRequest);
        }
      }
      clearAuthTokens();
      window.location.href = "/#/auth/login";
    }
    return Promise.reject(error);
  },
);
```

- 客户端使用固定的 `/api/auth` 前缀，同时覆盖 `/api/auth/user/*` 与 `/api/auth/user/tokens`。
- 刷新令牌失败会清理本地状态并回到登录页，确保 401 循环被阻断。

## 认证 API（`src/api/auth/auth.ts`）

```ts
import { apiClient } from "./client";
import { saveAccessToken, saveRefreshToken, clearAuthTokens } from "@/utils/auth";
import type { LoginRequest, RegisterRequest, AuthResponse, ApiResponse } from "@/types/auth";

export const login = async (req: LoginRequest): Promise<AuthResponse> => {
  const { data } = await apiClient.post<ApiResponse<AuthResponse>>("/login", req);
  if (data.data) {
    saveAccessToken(data.data.access_token);
    saveRefreshToken(data.data.refresh_token);
    return data.data;
  }
  throw new Error(data.error || "Login failed");
};

export const register = async (req: RegisterRequest): Promise<AuthResponse> => {
  const { data } = await apiClient.post<ApiResponse<AuthResponse>>("/register", req);
  if (data.data) {
    saveAccessToken(data.data.access_token);
    saveRefreshToken(data.data.refresh_token);
    return data.data;
  }
  throw new Error(data.error || "Registration failed");
};

export const refreshToken = async (refreshToken: string): Promise<AuthResponse> => {
  const { data } = await apiClient.post<ApiResponse<AuthResponse>>("/refresh", {
    refresh_token: refreshToken,
  });
  if (data.data) {
    return data.data;
  }
  throw new Error(data.error || "Token refresh failed");
};

export const logout = () => {
  clearAuthTokens();
};
```

- `login`/`register` 均保存 token，`logout` 只清理客户端缓存。
- 带验证码或 2FA 的登陆流程由 `PlatformAuthAPI` 在同目录下实现，仍然共享 `apiClient`。

## 用户自助 API（`src/api/auth/user.ts`）

```ts
import { apiClient } from "./client";
import type { User, ApiResponse } from "@/types/auth";

export const getCurrentUser = async (): Promise<User> => {
  const { data } = await apiClient.get<ApiResponse<User>>("/me");
  if (data.data) {
    return data.data;
  }
  throw new Error(data.error || "Failed to get user info");
};

export interface ChangePasswordRequest {
  old_password: string;
  new_password: string;
}

export const changePassword = async (params: ChangePasswordRequest): Promise<void> => {
  await apiClient.put("/user/me/password", params);
};

export interface UpdateProfileRequest {
  full_name?: string;
  avatar?: string;
  bio?: string;
}

export const updateProfile = async (params: UpdateProfileRequest): Promise<User> => {
  const { data } = await apiClient.put<ApiResponse<User>>("/user/me", params);
  if (data.data) {
    return data.data;
  }
  throw new Error(data.error || "更新个人资料失败");
};
```

- Path 全部基于 `/api/auth`，因此 `PUT /user/me` 实际请求 `/api/auth/user/me`。
- 错误直接抛出以便 Store/组件统一处理。

## Personal Access Token API（`src/api/user/tokens.ts`）

```ts
import { apiClient } from "../auth/client";
import type { PersonalAccessToken, CreateTokenRequest, CreateTokenResponse } from "@/types/user";
import type { ApiResponse } from "@/types/auth";

export const listTokens = async (): Promise<PersonalAccessToken[]> => {
  const { data } = await apiClient.get<ApiResponse<PersonalAccessToken[]>>("/user/tokens");
  if (data.data) {
    return data.data;
  }
  throw new Error(data.error || "获取 Token 列表失败");
};

export const createToken = async (params: CreateTokenRequest): Promise<CreateTokenResponse> => {
  const { data } = await apiClient.post<ApiResponse<CreateTokenResponse>>("/user/tokens", params);
  if (data.data) {
    return data.data;
  }
  throw new Error(data.error || "创建 Token 失败");
};

export const revokeToken = async (id: number): Promise<void> => {
  await apiClient.delete(`/user/tokens/${id}`);
};
```

- 成功创建时仅返回一次明文 token（`CreateTokenResponse`），需要立即显示给用户。
- 其余字段（前缀、权限、过期时间）与后端 PAT 模块完全一致。

## API 聚合导出

`src/api/index.ts` 统一 re-export：

```ts
export * from "./auth"; // auth/auth.ts + auth/user.ts + auth/client.ts + platformAuth.ts
export * from "./user";
export type {
  CaptchaData,
  LoginRequest,
  PlatformLoginRequest,
  AuthResponse,
  LoginResult,
  User,
} from "@/types";
```

组件或 Store 只需一次导入：

```ts
import { login, getCurrentUser, listTokens } from "@/api";

await login({ login: "admin", password: "secret" });
const profile = await getCurrentUser();
const tokens = await listTokens();
```

## Store / 组件调用示例

`src/stores/auth.ts` 直接引用上述 API：

```ts
const response = await login(credentials);
currentUser.value = response.user;

const user = await getCurrentUser();
currentUser.value = user;
```

`src/pages/user/tokens/index.vue` 也可通过 `listTokens` / `createToken` 构建 PAT UI，而无需重复拼接 URL。

---

通过集中式客户端 + 模块化 API，前端与后端 DDD 模块保持一一对应，Token 注入、刷新和错误处理都在同一位置完成，避免在组件层散落重复逻辑。
