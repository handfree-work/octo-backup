export const aboutResource = [
  {
    title: "关于",
    name: "about",
    path: "/about",
    redirect: "/about/index",
    meta: { icon: "lucide:circle-help", order: 9999 },
    children: [{
      title: "项目介绍",
      name: "aboutIndex",
      path: "/about/index",
      component: "/framework/about/index.vue",
      meta: { icon: "lucide:info" }
    }]
  }
];
export default aboutResource;
