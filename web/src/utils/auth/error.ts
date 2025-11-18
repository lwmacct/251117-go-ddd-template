/**
 * 错误处理工具
 */

/**
 * 格式化错误信息
 * @param error - 错误对象
 * @returns 格式化后的错误消息
 */

export function formatAuthError(error: any): string {
  // 优先使用服务器返回的错误信息
  if (error.response?.data?.error) {
    return error.response.data.error;
  }

  // 使用错误消息
  if (error.message) {
    return error.message;
  }

  // HTTP 状态码错误
  if (error.response?.status) {
    const status = error.response.status;
    switch (status) {
      case 400:
        return "请求参数错误";
      case 401:
        return "认证失败，请检查您的用户名和密码";
      case 403:
        return "没有权限访问";
      case 404:
        return "请求的资源不存在";
      case 500:
        return "服务器内部错误";
      default:
        return `请求失败 (${status})`;
    }
  }

  // 网络错误
  if (error.code === "ECONNABORTED" || error.message === "Network Error") {
    return "网络连接失败，请检查您的网络";
  }

  // 默认错误
  return "操作失败，请重试";
}
