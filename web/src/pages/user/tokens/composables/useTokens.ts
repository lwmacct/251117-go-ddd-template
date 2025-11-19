/**
 * Personal Access Token 管理 Composable
 */
import { ref } from "vue";
import { UserTokensAPI } from "@/api/user";
import type { PersonalAccessToken, CreateTokenRequest, CreateTokenResponse } from "@/types/user";

export function useTokens() {
  const tokens = ref<PersonalAccessToken[]>([]);
  const loading = ref(false);
  const errorMessage = ref("");
  const successMessage = ref("");

  const fetchTokens = async () => {
    loading.value = true;
    errorMessage.value = "";

    try {
      tokens.value = await UserTokensAPI.listTokens();
    } catch (error: any) {
      errorMessage.value = error.message || "获取 Token 列表失败";
    } finally {
      loading.value = false;
    }
  };

  const createToken = async (data: CreateTokenRequest): Promise<CreateTokenResponse | null> => {
    loading.value = true;
    errorMessage.value = "";
    successMessage.value = "";

    try {
      const response = await UserTokensAPI.createToken(data);
      successMessage.value = "Token 创建成功";
      await fetchTokens();
      return response;
    } catch (error: any) {
      errorMessage.value = error.message || "创建 Token 失败";
      return null;
    } finally {
      loading.value = false;
    }
  };

  const revokeToken = async (id: number): Promise<boolean> => {
    loading.value = true;
    errorMessage.value = "";
    successMessage.value = "";

    try {
      await UserTokensAPI.revokeToken(id);
      successMessage.value = "Token 已撤销";
      await fetchTokens();
      return true;
    } catch (error: any) {
      errorMessage.value = error.message || "撤销 Token 失败";
      return false;
    } finally {
      loading.value = false;
    }
  };

  const clearMessages = () => {
    errorMessage.value = "";
    successMessage.value = "";
  };

  return {
    tokens,
    loading,
    errorMessage,
    successMessage,
    fetchTokens,
    createToken,
    revokeToken,
    clearMessages,
  };
}
