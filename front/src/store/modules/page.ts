import { defineStore } from "pinia";

export const usePageStore = defineStore("page", {
  state: () => ({ keepAlive: [] as string[] }),
  actions: {
    close() {},
    closeLeft() {},
    closeRight() {},
    closeOther() {},
    closeAll() {},
    openedSort() {},
    getOpened() {
      return [] as any[];
    }
  }
});
