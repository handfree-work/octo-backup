import { beforeEach, describe, expect, it, vi } from "vitest";

const { request } = vi.hoisted(() => ({ request: vi.fn() }));
vi.mock("/src/api/service", () => ({ request }));

import { login, register } from "./api";

describe("auth api password hashing", () => {
  beforeEach(() => request.mockReset().mockResolvedValue({}));

  it("hashes login password before sending", async () => {
    await login({ username: "alice", password: "secret" });
    expect(request).toHaveBeenCalledWith(expect.objectContaining({ data: { username: "alice", password: "5ebe2294ecd0e0f08eab7690d2a6ee69" } }));
  });

  it("hashes register password before sending", async () => {
    await register({ username: "alice", password: "secret" });
    expect(request).toHaveBeenCalledWith(expect.objectContaining({ data: { username: "alice", password: "5ebe2294ecd0e0f08eab7690d2a6ee69" } }));
  });
});
