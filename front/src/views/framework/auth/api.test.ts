import { beforeEach, describe, expect, it, vi } from "vitest";

const { request } = vi.hoisted(() => ({ request: vi.fn() }));
vi.mock("/src/api/service", () => ({ request }));

import { login, register } from "./api";

describe("auth api password hashing", () => {
  beforeEach(() => request.mockReset().mockResolvedValue({}));

  it("hashes login password before sending", async () => {
    await login({ username: "alice", password: "secret" });
    expect(request).toHaveBeenCalledWith(expect.objectContaining({ data: { username: "alice", password: "2bb80d537b1da3e38bd30361aa855686bde0eacd7162fef6a25fe97bf527a25b" } }));
  });

  it("hashes register password before sending", async () => {
    await register({ username: "alice", password: "secret" });
    expect(request).toHaveBeenCalledWith(expect.objectContaining({ data: { username: "alice", password: "2bb80d537b1da3e38bd30361aa855686bde0eacd7162fef6a25fe97bf527a25b" } }));
  });
});
