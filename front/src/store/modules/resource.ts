import { defineStore } from "pinia";

export const useResourceStore = defineStore("resource", {
  state: () => ({ frameworkMenus: [] as any[], headerMenus: [] as any[], asideMenus: [] as any[] }),
  getters: {
    getFrameworkMenus: (state) => state.frameworkMenus,
    getHeaderMenus: (state) => state.headerMenus,
    getAsideMenus: (state) => state.asideMenus
  }
});
