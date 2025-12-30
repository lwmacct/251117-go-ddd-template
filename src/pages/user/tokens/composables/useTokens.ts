/**
 * Personal Access Token 管理 Composable
 */
import { ref } from "vue";
import { userTokensApi, extractData, type PatTokenDTO, type PatCreateDTO, type PatCreateResultDTO } from "@/api";
import { useSnackbar } from "@/composables";

export function useTokens() {
  const tokens = ref<PatTokenDTO[]>([]);
  const loading = ref(false);

  // 消息提示
  const { success, error } = useSnackbar();

  const fetchTokens = async () => {
    loading.value = true;

    try {
      const response = await userTokensApi.apiUserTokensGet();
      tokens.value = response.data.data ?? [];
    } catch (err) {
      error((err as Error).message || "获取 Token 列表失败");
    } finally {
      loading.value = false;
    }
  };

  const createToken = async (data: PatCreateDTO): Promise<PatCreateResultDTO | null> => {
    loading.value = true;

    try {
      const response = await userTokensApi.apiUserTokensPost(data);
      const result = extractData<PatCreateResultDTO>(response.data);
      success("Token 创建成功");
      await fetchTokens();
      return result ?? null;
    } catch (err) {
      error((err as Error).message || "创建 Token 失败");
      return null;
    } finally {
      loading.value = false;
    }
  };

  const deleteToken = async (id: number): Promise<boolean> => {
    loading.value = true;

    try {
      await userTokensApi.apiUserTokensIdDelete(id);
      success("Token 已删除");
      await fetchTokens();
      return true;
    } catch (err) {
      error((err as Error).message || "删除 Token 失败");
      return false;
    } finally {
      loading.value = false;
    }
  };

  const disableToken = async (id: number): Promise<boolean> => {
    loading.value = true;

    try {
      await userTokensApi.apiUserTokensIdDisablePatch(id);
      success("Token 已禁用");
      await fetchTokens();
      return true;
    } catch (err) {
      error((err as Error).message || "禁用 Token 失败");
      return false;
    } finally {
      loading.value = false;
    }
  };

  const enableToken = async (id: number): Promise<boolean> => {
    loading.value = true;

    try {
      await userTokensApi.apiUserTokensIdEnablePatch(id);
      success("Token 已启用");
      await fetchTokens();
      return true;
    } catch (err) {
      error((err as Error).message || "启用 Token 失败");
      return false;
    } finally {
      loading.value = false;
    }
  };

  return {
    tokens,
    loading,
    fetchTokens,
    createToken,
    deleteToken,
    disableToken,
    enableToken,
  };
}
