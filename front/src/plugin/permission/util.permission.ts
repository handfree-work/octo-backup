import { usePermissionStore } from "./store.permission";
import { NoPermissionError } from "./errors";
import { message } from "ant-design-vue";
import { useUserStore } from "/@/store/modules/user";
const util = {
  hasPermissions: (value: string | string[]): boolean => {
    let need: string[] = [];
    if (typeof value === "string") {
      need.push(value);
    } else if (value && value instanceof Array && value.length > 0) {
      need = need.concat(value);
    }
    if (need.length === 0) {
      throw new Error('need permissions! Like "sys:user:view" ');
    }
    const userRole = useUserStore().getUserInfo?.role;
    // 角色权限按 admin > write > read 逐级继承，路由只需声明最低角色。
    const roleLevel: Record<string, number> = { read: 1, write: 2, admin: 3 };
    const currentLevel = typeof userRole === "string" ? roleLevel[userRole] || 0 : 0;
    return currentLevel > 0 && need.some((required) => currentLevel >= (roleLevel[required] || 0));
  },
  requirePermissions: (value: any) => {
    if (!util.hasPermissions(value)) {
      message.error("对不起，您没有权限执行此操作");
      throw new NoPermissionError();
    }
  }
};

export default util;
