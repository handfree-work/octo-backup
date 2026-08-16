import { errorCreate } from "./tools";

/** 解包项目后端和旧 mock 共同使用的响应 envelope。 */
export function unpackResponseData(dataAxios: any, unpack = true): any {
  if (!unpack) {
    return dataAxios;
  }

  if (dataAxios && typeof dataAxios === "object" && "code" in dataAxios) {
    if (dataAxios.code === 0) {
      return dataAxios.data;
    }
    errorCreate(`${dataAxios.msg || "请求失败"}`);
    return dataAxios;
  }

  if (dataAxios && typeof dataAxios === "object" && "data" in dataAxios) {
    return dataAxios.data;
  }

  errorCreate(`非标准返回：${dataAxios}`);
  return dataAxios;
}
