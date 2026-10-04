declare module "vue-tippy" {
  import type { Plugin } from "vue";
  export const Tippy: any;
  export const useTippy: any;
  export const setDefaultProps: (props: any) => void;
  const plugin: Plugin;
  export default plugin;
}
declare module "tippy.js" {
  export type DefaultProps = any;
  export type Props = any;
  export const createTippy: any;
}
declare module "@vueuse/integrations/useQRCode" {
  export const useQRCode: any;
}
declare module "sortablejs" {
  const Sortable: any;
  export default Sortable;
}
declare module "sortablejs/modular/sortable.complete.esm.js" {
  const Sortable: any;
  export default Sortable;
}
declare module "@vben/preferences" {
  export const preferences: any;
  export const updatePreferences: (value: any) => void;
}
declare module "@vben/stores" {
  export const useAccessStore: any;
  export const useUserStore: any;
}
declare module "/@/vben/types/global" {
  export type AppConfig = any;
}
declare global {
  interface Window {
    _VBEN_ADMIN_PRO_APP_CONF_?: any;
  }
}
