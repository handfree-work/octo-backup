import path from 'node:path';
import process from 'node:process';

import tailwindcss from '@tailwindcss/vite';
import vue from '@vitejs/plugin-vue';
import vueJsx from '@vitejs/plugin-vue-jsx';
import { defineConfig, loadEnv } from 'vite';

import { createLocalMockApiPlugin } from './vite.mock';

const root = process.cwd();
const vbenPackagesRoot = path.resolve(root, 'src/vben-packages');

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, root);
  let devAlias = [];
  let devServerFs = {};

  if (mode.startsWith('debug')) {
    devAlias = [
      {
        find: /@fast-crud\/fast-crud\/dist/,
        replacement: path.resolve('../../fast-crud/src'),
      },
      {
        find: /@fast-crud\/ui-antdv-next\/dist/,
        replacement: path.resolve('../../ui/ui-antdv-next/src'),
      },
    ];
    devServerFs = {
      allow: ['../../../'],
    };
    console.log('devAlias', devAlias);
  }

  return {
    plugins: [
      vue(),
      vueJsx(),
      tailwindcss(),
      ...(env.VITE_NITRO_MOCK === 'true' ? [createLocalMockApiPlugin()] : []),
    ],
    resolve: {
      alias: [
        ...devAlias,
        { find: '#', replacement: path.resolve(root, 'src') },
        {
          find: /^@vben\/common-ui\/es\/(.+)$/,
          replacement: `${vbenPackagesRoot}/@vben/common-ui/components/$1`,
        },
        {
          find: '@vben/styles/antdv-next',
          replacement: `${vbenPackagesRoot}/@vben/styles/antdv-next/index.css`,
        },
        {
          find: '@vben/tailwind-config/theme',
          replacement: `${vbenPackagesRoot}/@vben/tailwind-config/theme.css`,
        },
        {
          find: '@vben-core/design/bem',
          replacement: `${vbenPackagesRoot}/@vben-core/design/less-bem/bem.less`,
        },
        {
          find: '@vben-core/design/theme',
          replacement: `${vbenPackagesRoot}/@vben-core/design/css/global.css`,
        },
        {
          find: /^@vben\/plugins\/(.+)$/,
          replacement: `${vbenPackagesRoot}/@vben/plugins/$1/index.ts`,
        },
        {
          find: /^@vben\/styles\/(.+)$/,
          replacement: `${vbenPackagesRoot}/@vben/styles/$1`,
        },
        {
          find: /^@vben\/tailwind-config\/(.+)$/,
          replacement: `${vbenPackagesRoot}/@vben/tailwind-config/$1`,
        },
        {
          find: /^@vben-core\/design\/(.+)$/,
          replacement: `${vbenPackagesRoot}/@vben-core/design/$1`,
        },
        {
          find: /^@vben-core\/([^/]+)\/(.+)$/,
          replacement: `${vbenPackagesRoot}/@vben-core/$1/$2`,
        },
        {
          find: /^@vben\/([^/]+)\/(.+)$/,
          replacement: `${vbenPackagesRoot}/@vben/$1/$2`,
        },
        { find: '@vben/access', replacement: `${vbenPackagesRoot}/@vben/access/index.ts` },
        { find: '@vben/common-ui', replacement: `${vbenPackagesRoot}/@vben/common-ui/index.ts` },
        { find: '@vben/constants', replacement: `${vbenPackagesRoot}/@vben/constants/index.ts` },
        { find: '@vben/hooks', replacement: `${vbenPackagesRoot}/@vben/hooks/index.ts` },
        { find: '@vben/icons', replacement: `${vbenPackagesRoot}/@vben/icons/index.ts` },
        { find: '@vben/layouts', replacement: `${vbenPackagesRoot}/@vben/layouts/index.ts` },
        { find: '@vben/locales', replacement: `${vbenPackagesRoot}/@vben/locales/index.ts` },
        { find: '@vben/plugins', replacement: `${vbenPackagesRoot}/@vben/plugins/index.ts` },
        { find: '@vben/preferences', replacement: `${vbenPackagesRoot}/@vben/preferences/index.ts` },
        { find: '@vben/request', replacement: `${vbenPackagesRoot}/@vben/request/index.ts` },
        { find: '@vben/stores', replacement: `${vbenPackagesRoot}/@vben/stores/index.ts` },
        { find: '@vben/styles', replacement: `${vbenPackagesRoot}/@vben/styles/index.ts` },
        {
          find: '@vben/tailwind-config',
          replacement: `${vbenPackagesRoot}/@vben/tailwind-config/index.ts`,
        },
        { find: '@vben/types', replacement: `${vbenPackagesRoot}/@vben/types/index.ts` },
        { find: '@vben/utils', replacement: `${vbenPackagesRoot}/@vben/utils/index.ts` },
        { find: '@vben-core/design', replacement: `${vbenPackagesRoot}/@vben-core/design/index.ts` },
        { find: '@vben-core/icons', replacement: `${vbenPackagesRoot}/@vben-core/icons/index.ts` },
        { find: '@vben-core/shared', replacement: `${vbenPackagesRoot}/@vben-core/shared/index.ts` },
        { find: '@vben-core/typings', replacement: `${vbenPackagesRoot}/@vben-core/typings/index.ts` },
        {
          find: '@vben-core/composables',
          replacement: `${vbenPackagesRoot}/@vben-core/composables/index.ts`,
        },
        {
          find: '@vben-core/preferences',
          replacement: `${vbenPackagesRoot}/@vben-core/preferences/index.ts`,
        },
        { find: '@vben-core/form-ui', replacement: `${vbenPackagesRoot}/@vben-core/form-ui/index.ts` },
        {
          find: '@vben-core/layout-ui',
          replacement: `${vbenPackagesRoot}/@vben-core/layout-ui/index.ts`,
        },
        { find: '@vben-core/menu-ui', replacement: `${vbenPackagesRoot}/@vben-core/menu-ui/index.ts` },
        { find: '@vben-core/popup-ui', replacement: `${vbenPackagesRoot}/@vben-core/popup-ui/index.ts` },
        {
          find: '@vben-core/shadcn-ui',
          replacement: `${vbenPackagesRoot}/@vben-core/shadcn-ui/index.ts`,
        },
        { find: '@vben-core/tabs-ui', replacement: `${vbenPackagesRoot}/@vben-core/tabs-ui/index.ts` },
      ],
    },
    define: {
      __APP_ENV__: JSON.stringify(env.VITE_APP_ENV),
    },
    server: {
      host: true,
      port: Number(env.VITE_PORT) || 5320,
      warmup: {
        clientFiles: [
          './index.html',
          './src/main.ts',
          './src/bootstrap.ts',
          './src/{views,layouts,router,store,api,adapter}/*',
        ],
      },
      fs: devServerFs,
      ...(env.VITE_API_PROXY_TARGET && env.VITE_NITRO_MOCK !== 'true'
        ? {
            proxy: {
              '/api': {
                changeOrigin: true,
                rewrite: (p) => p.replace(/^\/api/, ''),
                target: env.VITE_API_PROXY_TARGET,
                ws: true,
              },
            },
          }
        : {}),
    },
    build: {
      chunkSizeWarningLimit: 2000,
      reportCompressedSize: false,
      sourcemap: false,
      target: 'es2020',
    },
  };
});
