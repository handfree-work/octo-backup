import { defineStore } from "pinia";
import { getPluginMetadata, type PluginMetadata } from "/src/views/sys/plugin/plugin-api";

interface PluginState {
  defines: PluginMetadata[] | null;
}

let definesRequest: Promise<PluginMetadata[]> | undefined;
let storeEpoch = 0;

export const usePluginDefineStore = defineStore("app.pluginDefine", {
  state: (): PluginState => ({ defines: null }),
  actions: {
    async init() {
      if (this.defines) return this.defines;
      const requestEpoch = storeEpoch;
      const request = (definesRequest ||= getPluginMetadata());
      try {
        const defines = await request;
        if (requestEpoch === storeEpoch) this.defines = defines;
        return defines;
      } catch (error) {
        if (definesRequest === request) definesRequest = undefined;
        throw error;
      }
    },
    async getPluginDefine(name: string) {
      const defines = await this.init();
      return defines.find((define) => define.name === name);
    },
    async clear() {
      storeEpoch++;
      this.defines = null;
      definesRequest = undefined;
    },
    async reload() {
      await this.clear();
      await this.init();
    }
  }
});
