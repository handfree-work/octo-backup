import { createIconifyIcon } from "./create-icon";

export * from "./create-icon";

export * from "./lucide";

export type { IconifyIcon as IconifyIconStructure } from "@iconify/vue";
export { addCollection, addIcon, Icon as IconifyIcon, listIcons } from "@iconify/vue";

export const MdiKeyboardEsc = createIconifyIcon("mdi:keyboard-esc");
export const EmptyIcon = createIconifyIcon("mdi:minus");
export const MdiGithub = createIconifyIcon("mdi:github");
export const MdiGoogle = createIconifyIcon("mdi:google");
export const MdiQqchat = createIconifyIcon("mdi:qqchat");
export const MdiWechat = createIconifyIcon("mdi:wechat");
