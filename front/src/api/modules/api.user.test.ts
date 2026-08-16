import { beforeEach, describe, expect, it, vi } from "vitest";

const request = vi.hoisted(() => vi.fn());

vi.mock("../service", () => ({ request }));

import { login, register } from "./api.user";

describe("认证 API", () => {
  beforeEach(() => {
    request.mockReset();
  });

  it("通过 POST /auth/login 登录", async () => {
    request.mockResolvedValue({ token: "jwt" });

    await login({ username: "admin", password: "123456" });

    expect(request).toHaveBeenCalledWith({
      url: "/auth/login",
      method: "post",
      data: { username: "admin", password: "123456" }
    });
  });

  it("通过 POST /auth/register 注册", async () => {
    request.mockResolvedValue({ id: 1, username: "admin", role: "read" });

    await register({ username: "admin", password: "123456" });

    expect(request).toHaveBeenCalledWith({
      url: "/auth/register",
      method: "post",
      data: { username: "admin", password: "123456" }
    });
  });
});
