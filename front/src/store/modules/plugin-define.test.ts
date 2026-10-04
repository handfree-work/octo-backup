import { beforeEach, describe, expect, it, vi } from "vitest";
import { createPinia, setActivePinia } from "pinia";
import { getPluginMetadata } from "/src/views/sys/plugin/plugin-api";
import { usePluginDefineStore } from "./plugin-define";

vi.mock("/src/views/sys/plugin/plugin-api", () => ({
  getPluginMetadata: vi.fn()
}));

describe("plugin store", () => {
  beforeEach(async () => {
    setActivePinia(createPinia());
    vi.clearAllMocks();
    vi.mocked(getPluginMetadata).mockResolvedValue([{ type: "access", name: "access.ssh", title: "SSH", description: "", version: "1", fields: [{ key: "host", title: "Host", type: "string" }] }]);
    await usePluginDefineStore().clear();
  });

  it("loads plugin defines once", async () => {
    const store = usePluginDefineStore();
    const define = await store.getPluginDefine("access.ssh");
    expect(await store.getPluginDefine("access.ssh")).toEqual(define);
    expect(getPluginMetadata).toHaveBeenCalledTimes(1);
  });

  it("reloads the metadata cache", async () => {
    const store = usePluginDefineStore();
    await store.init();
    await store.reload();
    expect(getPluginMetadata).toHaveBeenCalledTimes(2);
  });
});
