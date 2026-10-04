import { beforeEach, describe, expect, it, vi } from "vitest";
import { createPinia, setActivePinia } from "pinia";
import { createPluginInstance, getPluginInstanceInfo, getPluginInstancePage, getPluginMetadata, updatePluginInstance } from "/src/views/sys/plugin/plugin-api";
import { usePluginDefineStore } from "/src/store/modules/plugin-define";
import createPluginCrudOptions, { setPluginAccessType } from "/src/views/sys/plugin/crud";

vi.mock("/src/views/sys/plugin/plugin-api", () => ({
  getPluginMetadata: vi.fn(),
  getPluginInstanceInfo: vi.fn(),
  getPluginInstancePage: vi.fn(),
  createPluginInstance: vi.fn(),
  updatePluginInstance: vi.fn(),
  deletePluginInstance: vi.fn()
}));

describe("plugin metadata loading", () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    vi.mocked(getPluginMetadata).mockResolvedValue([]);
  });

  it("loads the selected plugin define once", async () => {
    const store = usePluginDefineStore();
    const define = await store.getPluginDefine("access.tencent");
    expect(await store.getPluginDefine("access.tencent")).toEqual(define);
    expect(getPluginMetadata).toHaveBeenCalledTimes(1);
  });
});

describe("plugin selector constraints", () => {
  it("copies the selected repository define access type to the form", () => {
    const form: Record<string, unknown> = {};

    setPluginAccessType(form, "access.ssh");

    expect(form.config).toEqual({ accessType: "access.ssh" });
  });

  it("clears the previous authorization when the selected access type changes", () => {
    const form = { config: { accessType: "access.s3", accessId: 42 } };

    setPluginAccessType(form, "access.ssh", true);

    expect(form.config).toEqual({ accessType: "access.ssh", accessId: undefined });
  });
});

describe("plugin config requests", () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    vi.clearAllMocks();
  });

  it("keeps accessType in config sent by create and update requests", async () => {
    const { crudOptions } = createPluginCrudOptions({
      context: { pluginType: "repository", metadata: [] },
      crudExpose: { getFormWrapperRef: () => undefined }
    } as any);
    const config = { accessId: 7, accessType: "access.ssh", path: "/backup" };
    const form = { name: "backup", pluginName: "repository.sftp", config };

    await (crudOptions.request as any).addRequest({ form });
    await (crudOptions.request as any).editRequest({ form, row: { id: 12, pluginName: "repository.sftp" } });

    expect(createPluginInstance).toHaveBeenCalledWith(expect.objectContaining({ config }));
    expect(updatePluginInstance).toHaveBeenCalledWith(12, expect.objectContaining({ config }));
  });

  it("loads a full plugin record by id before editing", async () => {
    const { crudOptions } = createPluginCrudOptions({
      context: { pluginType: "repository", metadata: [] },
      crudExpose: { getFormWrapperRef: () => undefined }
    } as any);
    const row = { id: 12, name: "backup", pluginName: "repository.sftp", config: { path: "/backup" } };
    vi.mocked(getPluginInstanceInfo).mockResolvedValue(row as any);

    await expect((crudOptions.request as any).infoRequest({ row })).resolves.toEqual(row);
    expect(getPluginInstanceInfo).toHaveBeenCalledWith(12);
  });

  it("filters instances through the instance page API", async () => {
    const { crudOptions } = createPluginCrudOptions({ context: { pluginType: "access", pluginName: "access.ssh", metadata: [] }, crudExpose: { getFormWrapperRef: () => undefined } } as any);
    await (crudOptions.request as any).pageRequest({ currentPage: 1, pageSize: 20, query: {} });
    expect(getPluginInstancePage).toHaveBeenCalledWith({ offset: undefined, limit: 20, pluginType: "access", pluginName: "access.ssh", name: undefined });
  });

  it("loads plugin name options from the plugin define store", async () => {
    await usePluginDefineStore().clear();
    vi.mocked(getPluginMetadata).mockResolvedValue([
      { name: "access.ssh", title: "SSH", type: "access", fields: [] },
      { name: "repository.local", title: "本地目录", type: "repository", fields: [] }
    ] as any);
    const { crudOptions } = createPluginCrudOptions({
      context: { pluginType: "access", metadata: [] },
      crudExpose: { getFormWrapperRef: () => undefined }
    } as any);

    await expect((crudOptions.columns as any).pluginName.dict.getData()).resolves.toEqual([{ value: "access.ssh", label: "SSH" }]);
    expect(getPluginMetadata).toHaveBeenCalledTimes(1);
  });
});
