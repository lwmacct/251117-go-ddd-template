/**
 * Personal Access Token 管理 Composable
 */
import { ref } from "vue";
import { userTokensApi, extractData, type PatTokenDTO, type PatCreateTokenDTO, type PatCreateTokenResultDTO } from "@/api";

export function useTokens() {
  const tokens = ref<PatTokenDTO[]>([]);
  const loading = ref(false);
  const errorMessage = ref("");
  const successMessage = ref("");

  const fetchTokens = async () => {
    loading.value = true;
    errorMessage.value = "";

    try {
      const response = await userTokensApi.apiUserTokensGet();
      tokens.value = response.data.data ?? [];
    } catch (error) {
      errorMessage.value = (error as Error).message || "获取 Token 列表失败";
    } finally {
      loading.value = false;
    }
  };

  const createToken = async (data: PatCreateTokenDTO): Promise<PatCreateTokenResultDTO | null> => {
    loading.value = true;
    errorMessage.value = "";
    successMessage.value = "";

    try {
      const response = await userTokensApi.apiUserTokensPost(data);
      const result = extractData<PatCreateTokenResultDTO>(response.data);
      successMessage.value = "Token 创建成功";
      await fetchTokens();
      return result ?? null;
    } catch (error) {
      errorMessage.value = (error as Error).message || "创建 Token 失败";
      return null;
    } finally {
      loading.value = false;
    }
  };

  const deleteToken = async (id: number): Promise<boolean> => {
    loading.value = true;
    errorMessage.value = "";
    successMessage.value = "";

    try {
      await userTokensApi.apiUserTokensIdDelete(id);
      successMessage.value = "Token 已删除";
      await fetchTokens();
      return true;
    } catch (error) {
      errorMessage.value = (error as Error).message || "删除 Token 失败";
      return false;
    } finally {
      loading.value = false;
    }
  };

  const disableToken = async (id: number): Promise<boolean> => {
    loading.value = true;
    errorMessage.value = "";
    successMessage.value = "";

    try {
      await userTokensApi.apiUserTokensIdDisablePatch(id);
      successMessage.value = "Token 已禁用";
      await fetchTokens();
      return true;
    } catch (error) {
      errorMessage.value = (error as Error).message || "禁用 Token 失败";
      return false;
    } finally {
      loading.value = false;
    }
  };

  const enableToken = async (id: number): Promise<boolean> => {
    loading.value = true;
    errorMessage.value = "";
    successMessage.value = "";

    try {
      await userTokensApi.apiUserTokensIdEnablePatch(id);
      successMessage.value = "Token 已启用";
      await fetchTokens();
      return true;
    } catch (error) {
      errorMessage.value = (error as Error).message || "启用 Token 失败";
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
    deleteToken,
    disableToken,
    enableToken,
    clearMessages,
  };
}
