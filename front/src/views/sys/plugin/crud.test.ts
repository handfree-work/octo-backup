import { describe, expect, it } from "vitest";
import { collectPluginConfig } from "./config";

const metadata = [
  {
    type: "access",
    name: "access.tencent",
    title: "腾讯云授权",
    description: "",
    version: "1.0.0",
    fields: [
      { key: "secretId", title: "SecretId", type: "string", required: true },
      { key: "secretKey", title: "SecretKey", type: "password", encrypt: true },
      { key: "region", title: "区域", type: "string" }
    ]
  }
];

describe("collectPluginConfig", () => {
  it("collects plugin fields from the nested config form model", () => {
    expect(collectPluginConfig({ pluginName: "access.tencent", config: { secretId: "122", secretKey: "333", region: "ap-shanghai" } }, metadata)).toEqual({
      secretId: "122",
      secretKey: "333",
      region: "ap-shanghai"
    });
  });

  it("supports flat fields and ignores the configured response marker", () => {
    expect(collectPluginConfig({ pluginName: "access.tencent", secretId: "122", secretKey: "", config: { configured: true } }, metadata)).toEqual({ secretId: "122" });
  });

  it("does not submit an unchanged masked encrypted field", () => {
    expect(collectPluginConfig({ pluginName: "access.tencent", config: { secretKey: "to****et", secretKeyConfigured: true, secretKeyOriginal: "to****et" } }, metadata)).toEqual({});
  });

  it("submits an encrypted field when the masked value is changed", () => {
    expect(collectPluginConfig({ pluginName: "access.tencent", config: { secretKey: "to****at", secretKeyConfigured: true, secretKeyOriginal: "to****et" } }, metadata)).toEqual({ secretKey: "to****at" });
  });
});
