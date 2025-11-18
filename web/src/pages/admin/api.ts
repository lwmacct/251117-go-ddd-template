/**
 * 快手 API 服务 - 极简版
 */
import axios from "axios";
import type { ApiResponse, KuaishouTextResponse } from "./types";

const API_BASE_URL = "/api/kuaishou";

/** 默认配置 */
const DEFAULT_DB_ID = "bendy93687365";
const DEFAULT_PROVIDER = "bendy";

/** axios 实例 */
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
});

/**
 * 通用代理请求（核心方法）
 * @param path 完整路径，如：/provider/bendy93687365/20251108/pressNodeList/pressNodeList.txt
 */
export const proxyRequest = async (path: string) => {
  const { data } = await apiClient.get<ApiResponse<KuaishouTextResponse>>(
    "/proxy",
    { params: { path } }
  );
  return data;
};

/**
 * 以下是常用接口的快捷函数（可选使用）
 * 所有函数都是对 proxyRequest 的简单封装
 */

/** 压测节点列表 */
export const getPressNodeList = (params: { day: string; db_id?: string }) => {
  const dbId = params.db_id || DEFAULT_DB_ID;
  return proxyRequest(
    `/provider/${dbId}/${params.day}/pressNodeList/pressNodeList.txt`
  );
};

/** 黑名单列表 */
export const getLimitNodeList = (params: { day: string; db_id?: string }) => {
  const dbId = params.db_id || DEFAULT_DB_ID;
  return proxyRequest(
    `/provider/${dbId}/${params.day}/limitNodeList/limitNodeList.txt`
  );
};

/** 问题资源 */
export const getProblemResource = (params: { day: string; db_id?: string }) => {
  const dbId = params.db_id || DEFAULT_DB_ID;
  return proxyRequest(
    `/provider/${dbId}/${params.day}/queryProblemResource/queryProblemResource.txt`
  );
};

/** SLA 统计数据 */
export const getSLA = (params: {
  day: string;
  provider?: string;
  db_id?: string;
}) => {
  const dbId = params.db_id || DEFAULT_DB_ID;
  const provider = params.provider || DEFAULT_PROVIDER;
  return proxyRequest(
    `/provider/${dbId}/${params.day}/sla/${provider}_${params.day}_sla.txt`
  );
};

/** 机器列表状态 */
export const getNodeList = (params: {
  day: string;
  provider_name?: string;
  db_id?: string;
}) => {
  const dbId = params.db_id || DEFAULT_DB_ID;
  const provider = params.provider_name || DEFAULT_PROVIDER;
  return proxyRequest(
    `/provider/${dbId}/${params.day}/allNodeList/${provider}_${params.day}_nodeList.txt`
  );
};

/** 质量指标 */
export const getQualityIndex = (params: {
  datetime: string;
  time_slot: string;
  file_id: string;
  db_id?: string;
}) => {
  const dbId = params.db_id || DEFAULT_DB_ID;
  return proxyRequest(
    `/provider/${dbId}/${params.datetime}/${params.time_slot}/${params.file_id}`
  );
};

/** 需求缺口表 */
export const getResGap = (params: { type: "vod" | "live"; db_id?: string }) => {
  const dbId = params.db_id || DEFAULT_DB_ID;
  const fileName = params.type === "vod" ? "resgap.txt" : "live_resgap.txt";
  return proxyRequest(`/provider/${dbId}/public/${fileName}`);
};

/** 利用率规则说明 */
export const getUtilizationRules = () => {
  return proxyRequest("/provider/ks2022/public/rules.txt");
};
