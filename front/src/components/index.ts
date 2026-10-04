import { defineAsyncComponent } from "vue";

const AsyncHighLight = defineAsyncComponent(() => import("./highlight/index.vue"));
const AccessSelector = defineAsyncComponent(() => import("./access-selector/index.vue"));
const PluginSelector = defineAsyncComponent(() => import("./plugin-selector/index.vue"));
const PluginFileInput = defineAsyncComponent(() => import("./plugin-file-input/index.vue"));
const PluginPathSelector = defineAsyncComponent(() => import("./plugin-path-selector/index.vue"));
export default {
  install(app: any) {
    app.component("FsHighlight", AsyncHighLight);
    app.component("AccessSelector", AccessSelector);
    app.component("PluginSelector", PluginSelector);
    app.component("PluginFileInput", PluginFileInput);
    app.component("PluginPathSelector", PluginPathSelector);
  }
};
