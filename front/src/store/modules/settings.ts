import { defineStore } from "pinia";

export const useSettingStore = defineStore("settings", {
  state: () => ({ themeConfig: { colorPrimary: "#1677ff", mode: "light" as string } }),
  actions: {
    setPrimaryColor(color: string) {
      this.themeConfig.colorPrimary = color;
    },
    setDarkMode(mode: string) {
      this.themeConfig.mode = mode;
    }
  }
});
