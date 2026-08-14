import type { RouteRecordRaw } from 'vue-router';

import { $t } from '#/locales';

const routes: RouteRecordRaw[] = [
  {
    meta: {
      icon: 'ic:baseline-view-in-ar',
      keepAlive: true,
      order: 1000,
      title: $t('demos.title'),
    },
    name: 'Demos',
    path: '/demos',
    children: [
      {
        meta: {
          title: $t('demos.antd'),
        },
        name: 'AntDesignDemos',
        path: '/demos/ant-design-next',
        component: () => import('#/views/demos/antd/index.vue'),
      },
      {
        meta: {
          icon: 'lucide:table-2',
          title: 'Fast CRUD',
        },
        name: 'FastCrud',
        path: '/demos/fast-crud',
        children: [
          {
            meta: {
              title: 'Fast CRUD HelloWorld',
            },
            name: 'FastCrudHelloWorld',
            path: '/demos/fast-crud/helloworld',
            component: () =>
              import('#/views/demos/fast-crud/helloworld/index.vue'),
          },
          {
            meta: {
              icon: 'lucide:component',
              title: '组件示例',
            },
            name: 'FastCrudComponent',
            path: '/demos/fast-crud/component',
            children: [
              {
                meta: {
                  title: '文本输入(input)',
                },
                name: 'ComponentText',
                path: '/demos/fast-crud/component/text',
                component: () =>
                  import('#/views/demos/fast-crud/component/text/index.vue'),
              },
              {
                meta: {
                  title: '选择(select)',
                },
                name: 'ComponentSelect',
                path: '/demos/fast-crud/component/select',
                component: () =>
                  import('#/views/demos/fast-crud/component/select/index.vue'),
              },
            ],
          },
        ],
      },
    ],
  },
];

export default routes;
