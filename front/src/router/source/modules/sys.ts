import LayoutPass from "/@/layout/layout-pass.vue";
import BasicLayout from "/@/vben/layouts/basic/layout.vue";

export const sysResources = [
  {
    title: "系统管理",
    name: "sys",
    path: "/sys",
    redirect: "/sys/authority",
    meta: {
      icon: "ion:settings-outline",
      permission: "sys"
    },
    children: [
      {
        title: "用户管理",
        name: "user",
        meta: {
          icon: "ion:person-outline",
          // 用户管理属于管理员功能，管理员角色可访问所有更低级别权限。
          permission: "admin"
        },
        path: "/sys/authority/user",
        component: "/sys/authority/user/index.vue"
      },
      {
        title: "授权管理",
        name: "access",
        meta: { icon: "ion:key-outline", permission: "read" },
        path: "/sys/authority/access",
        component: "/sys/authority/access/index.vue"
      },
      {
        title: "存储仓库",
        name: "repository",
        meta: { icon: "ion:server-outline", permission: "read" },
        path: "/sys/repository",
        component: "/sys/repository/index.vue"
      },
      {
        title: "备份来源",
        name: "source",
        meta: { icon: "ion:folder-open-outline", permission: "read" },
        path: "/sys/source",
        component: "/sys/source/index.vue"
      },
      {
        title: "备份计划",
        name: "plan",
        meta: { icon: "ion:calendar-outline", permission: "read" },
        path: "/sys/plan",
        component: "/sys/plan/index.vue"
      },
      {
        title: "备份数据流",
        name: "plan-flow",
        meta: { icon: "ion:git-network-outline", permission: "read" },
        path: "/sys/plan/flow",
        component: "/sys/plan/flow/index.vue"
      }
    ]
  }
];

export default sysResources;
