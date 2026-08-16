import { describe, expect, it } from "vitest";

import { unpackResponseData } from "./response";

describe("响应数据解包", () => {
  it("解包 Go 后端的 data 响应", () => {
    expect(unpackResponseData({ data: { token: "jwt" } })).toEqual({ token: "jwt" });
  });

  it("保持 mock 接口的 code: 0 响应兼容", () => {
    expect(unpackResponseData({ code: 0, data: { token: "jwt" } })).toEqual({ token: "jwt" });
  });

  it("在请求方显式要求时保留完整 mock 响应", () => {
    const response = { code: 0, data: { token: "jwt" } };
    expect(unpackResponseData(response, false)).toBe(response);
  });
});
