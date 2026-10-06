import axios from "axios";
import { get } from "lodash-es";
import Adapter from "axios-mock-adapter";
import { errorLog } from "./tools";
import { env } from "/src/utils/util.env";
import { useUserStore } from "../store/modules/user";
import { unpackResponseData } from "./response";

export { unpackResponseData } from "./response";
/**
 * @description 创建请求实例
 */
function createService() {
  // 创建一个 axios 实例
  const service = axios.create();
  // 请求拦截
  service.interceptors.request.use(
    (config) => config,
    (error) => {
      // 发送失败
      console.log(error);
      return Promise.reject(error);
    }
  );
  // 响应拦截
  service.interceptors.response.use(
    (response) => {
      if (response.config.responseType === "blob") {
        return response;
      }
      // 同时兼容旧 mock 的 { code: 0, data } 和 Go 后端的 { data }。
      return unpackResponseData(response.data, (response.config as any).unpack !== false);
    },
    (error) => {
      const status = get(error, "response.status");
      const backendMessage = get(error, "response.data.error");
      switch (status) {
        case 400:
          error.message = backendMessage || "请求错误";
          break;
        case 401:
          error.message = backendMessage || "未授权，请登录";
          break;
        case 403:
          error.message = backendMessage || "拒绝访问";
          break;
        case 404:
          error.message = `请求地址出错: ${error.response.config.url}`;
          break;
        case 408:
          error.message = "请求超时";
          break;
        case 500:
          error.message = "服务器内部错误";
          break;
        case 501:
          error.message = "服务未实现";
          break;
        case 502:
          error.message = "网关错误";
          break;
        case 503:
          error.message = "服务不可用";
          break;
        case 504:
          error.message = "网关超时";
          break;
        case 505:
          error.message = "HTTP版本不受支持";
          break;
        default:
          break;
      }
      errorLog(error);
      if (status === 401) {
        const userStore = useUserStore();
        userStore.logout();
      }
      return Promise.reject(error);
    }
  );
  return service;
}

/**
 * @description 创建请求方法
 * @param {Object} service axios 实例
 */
function createRequestFunction(service: any) {
  return function (config: any) {
    const configDefault = {
      headers: {
        "Content-Type": get(config, "headers.Content-Type", "application/json")
      },
      timeout: 30000,
      baseURL: env.API,
      data: {}
    };
    const userStore = useUserStore();
    const token = userStore.getToken;
    if (token != null) {
      // @ts-ignore
      configDefault.headers.Authorization = `Bearer ${token}`;
    }
    const requestConfig = Object.assign({}, configDefault, config, {
      headers: Object.assign({}, configDefault.headers, config.headers || {})
    });
    return service(requestConfig);
  };
}

// 用于真实网络请求的实例和请求方法
export const service = createService();
export const request = createRequestFunction(service);

// 用于模拟网络请求的实例和请求方法
export const serviceForMock = createService();
export const requestForMock = createRequestFunction(serviceForMock);

// 网络请求数据模拟工具
export const mock = new Adapter(serviceForMock, { delayResponse: 200 });
