import LayoutOutside from "/src/layout/layout-outside.vue";
import Error404 from "/src/views/framework/error/404.vue";
const errorPage = [{ path: "/:pathMatch(.*)*", name: "not-found", component: Error404 }];
export const outsideResource = [
  {
    title: "Landing",
    name: "landing",
    path: "/",
    component: "/framework/landing/index.vue",
    meta: { title: "OctoBackup - 远程备份管理" }
  },
  {
    title: "outside",
    name: "outside",
    path: "/outside",
    component: LayoutOutside,
    children: [
      {
        meta: {
          title: "登录"
        },
        name: "login",
        path: "/login",
        component: "/framework/login/index.vue"
      },
      {
        meta: {
          title: "注册"
        },
        name: "register",
        path: "/register",
        component: "/framework/register/index.vue"
      }
    ]
  },
  ...errorPage
];
