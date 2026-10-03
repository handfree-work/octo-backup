export const crudResources = [
  {
    title: "CRUD示例",
    name: "crud",
    path: "/crud",
    redirect: "/crud/basis",
    meta: {
      icon: "ion:apps-sharp",
      isMenu: false
    },
    children: [
     
      {
        title: "基本特性",
        name: "basis",
        path: "/crud/basis",
        redirect: "/crud/basis/i18n",
        meta: {
          icon: "ion:disc-outline"
        },
        children: [
          {
            title: "FirstDemo",
            name: "FsCrudFirst",
            path: "/crud/basis/first",
            component: "/crud/basis/first/index.vue"
          },
        ]
      }
    ]
  }
];

export default crudResources;
