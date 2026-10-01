import { beforeEach, describe, expect, it, vi } from "vitest";

const { request } = vi.hoisted(() => ({ request: vi.fn() }));
vi.mock("/src/api/service", () => ({ request }));

import { createUser, updateUser } from "./api";

describe("user api password hashing", () => {
  beforeEach(() => request.mockReset().mockResolvedValue({}));

  it("hashes password when creating a user", async () => {
    await createUser({ username: "alice", password: "secret", role: "read" });
    expect(request).toHaveBeenCalledWith(expect.objectContaining({ data: expect.objectContaining({ password: "5ebe2294ecd0e0f08eab7690d2a6ee69" }) }));
  });

  it("omits blank password when updating a user", async () => {
    await updateUser(1, { password: "  ", nickName: "Alice" });
    expect(request).toHaveBeenCalledWith(expect.objectContaining({ data: { nickName: "Alice" } }));
  });

  it("hashes non-blank password when updating a user", async () => {
    await updateUser(1, { password: "secret" });
    expect(request).toHaveBeenCalledWith(expect.objectContaining({ data: { password: "5ebe2294ecd0e0f08eab7690d2a6ee69" } }));
  });
});
