import type { IncomingMessage, ServerResponse } from 'node:http';
import type { Plugin } from 'vite';

import { randomUUID } from 'node:crypto';

import qs from 'qs';

import componentSelectMocks from './src/views/demos/fast-crud/component/select/mock';
import componentTextMocks from './src/views/demos/fast-crud/component/text/mock';

interface MockUser {
  id: number;
  password: string;
  realName: string;
  roles: string[];
  username: string;
  homePath?: string;
}

const users: MockUser[] = [
  {
    id: 0,
    password: '123456',
    realName: 'Vben',
    roles: ['super'],
    username: 'vben',
  },
  {
    id: 1,
    password: '123456',
    realName: 'Admin',
    roles: ['admin'],
    username: 'admin',
    homePath: '/workspace',
  },
  {
    id: 2,
    password: '123456',
    realName: 'Jack',
    roles: ['user'],
    username: 'jack',
    homePath: '/analytics',
  },
];

const accessTokens = new Map<string, string>();
const refreshTokens = new Map<string, string>();

const dashboardMenus = [
  {
    meta: {
      order: -1,
      title: 'page.dashboard.title',
    },
    name: 'Dashboard',
    path: '/dashboard',
    redirect: '/analytics',
    children: [
      {
        name: 'Analytics',
        path: '/analytics',
        component: '/dashboard/analytics/index',
        meta: {
          affixTab: true,
          title: 'page.dashboard.analytics',
        },
      },
      {
        name: 'Workspace',
        path: '/workspace',
        component: '/dashboard/workspace/index',
        meta: {
          title: 'page.dashboard.workspace',
        },
      },
    ],
  },
];

const demoMenus = [
  {
    meta: {
      icon: 'ic:baseline-view-in-ar',
      order: 1000,
      title: 'demos.title',
    },
    name: 'Demos',
    path: '/demos',
    redirect: '/demos/antd',
    children: [
      {
        name: 'AntdDemo',
        path: '/demos/antd',
        component: '/demos/antd/index',
        meta: {
          title: 'demos.antd.title',
        },
      },
      {
        meta: {
          icon: 'lucide:table-2',
          title: 'Fast CRUD',
        },
        name: 'FastCrud',
        path: '/demos/fast-crud',
        redirect: '/demos/fast-crud/helloworld',
        children: [
          {
            meta: {
              title: 'Fast CRUD HelloWorld',
            },
            name: 'FastCrudHelloWorld',
            path: '/demos/fast-crud/helloworld',
            component: '/demos/fast-crud/helloworld/index',
          },
          {
            meta: {
              icon: 'lucide:component',
              title: '组件示例',
            },
            name: 'FastCrudComponent',
            path: '/demos/fast-crud/component',
            redirect: '/demos/fast-crud/component/text',
            children: [
              {
                meta: {
                  title: '文本输入(input)',
                },
                name: 'ComponentText',
                path: '/demos/fast-crud/component/text',
                component: '/demos/fast-crud/component/text/index',
              },
              {
                meta: {
                  title: '选择(select)',
                },
                name: 'ComponentSelect',
                path: '/demos/fast-crud/component/select',
                component: '/demos/fast-crud/component/select/index',
              },
            ],
          },
        ],
      },
    ],
  },
];

const viewMocks = [...componentTextMocks, ...componentSelectMocks];

function jsonResponse(
  response: ServerResponse,
  data: unknown,
  status = 200,
  headers: Record<string, string> = {},
) {
  response.statusCode = status;
  response.setHeader('Content-Type', 'application/json; charset=utf-8');
  Object.entries(headers).forEach(([key, value]) => {
    response.setHeader(key, value);
  });
  response.end(JSON.stringify(data));
}

function success(data: unknown) {
  return {
    code: 0,
    data,
    error: null,
    message: 'ok',
  };
}

function errorResponse(message: string) {
  return {
    code: -1,
    data: null,
    error: message,
    message,
  };
}

function getCookies(request: IncomingMessage) {
  return Object.fromEntries(
    (request.headers.cookie ?? '')
      .split(';')
      .filter(Boolean)
      .map((item) => {
        const [key, ...value] = item.trim().split('=');
        return [key, decodeURIComponent(value.join('='))];
      }),
  );
}

function getRequestParams(request: IncomingMessage) {
  const url = new URL(request.url ?? '/', 'http://localhost');
  return qs.parse(url.search.slice(1));
}

function getUserFromAccessToken(request: IncomingMessage) {
  const authorization = request.headers.authorization;
  const token = authorization?.match(/^Bearer\s+(.+)$/i)?.[1];
  const username = token ? accessTokens.get(token) : undefined;
  return users.find((user) => user.username === username);
}

function publicUser(user: MockUser) {
  const { password: _password, ...userinfo } = user;
  return userinfo;
}

function readJsonBody(request: IncomingMessage) {
  return new Promise<Record<string, unknown>>((resolve, reject) => {
    let body = '';
    request.setEncoding('utf8');
    request.on('data', (chunk) => {
      body += chunk;
    });
    request.on('end', () => {
      try {
        resolve(body ? JSON.parse(body) : {});
      } catch (error) {
        reject(error);
      }
    });
    request.on('error', reject);
  });
}

function getMenuList() {
  return [...dashboardMenus, ...demoMenus];
}

function getAccessCodes(username: string) {
  if (username === 'vben') {
    return ['AC_100100', 'AC_100110', 'AC_100120', 'AC_100010'];
  }
  if (username === 'admin') {
    return ['AC_100010', 'AC_100020', 'AC_100030'];
  }
  return ['AC_1000001', 'AC_1000002'];
}

export function createLocalMockApiPlugin(): Plugin {
  return {
    name: 'local-mock-api',
    configureServer(server) {
      server.middlewares.use(async (request, response, next) => {
        const method = request.method ?? 'GET';
        const pathname = new URL(
          request.url ?? '/',
          'http://localhost',
        ).pathname;
        const mockPathname = pathname.startsWith('/api/mock/')
          ? pathname.replace(/^\/api/, '')
          : pathname;

        if (!pathname.startsWith('/api/') && !pathname.startsWith('/mock/')) {
          next();
          return;
        }

        if (method === 'OPTIONS') {
          response.statusCode = 204;
          response.end();
          return;
        }

        try {
          if (pathname === '/api/auth/login' && method === 'POST') {
            const body = await readJsonBody(request);
            const username = String(body.username ?? '');
            const password = String(body.password ?? '');
            const user = users.find(
              (item) =>
                item.username === username && item.password === password,
            );

            if (!username || !password) {
              jsonResponse(
                response,
                errorResponse('Username and password are required'),
                400,
              );
              return;
            }

            if (!user) {
              jsonResponse(
                response,
                errorResponse('Username or password is incorrect.'),
                403,
              );
              return;
            }

            const accessToken = `mock-access-${randomUUID()}`;
            const refreshToken = `mock-refresh-${randomUUID()}`;
            accessTokens.set(accessToken, user.username);
            refreshTokens.set(refreshToken, user.username);

            jsonResponse(
              response,
              success({ ...publicUser(user), accessToken }),
              200,
              {
                'Set-Cookie': `refresh_token=${encodeURIComponent(refreshToken)}; Path=/; HttpOnly; SameSite=Lax`,
              },
            );
            return;
          }

          if (pathname === '/api/auth/refresh' && method === 'POST') {
            const refreshToken = getCookies(request).refresh_token;
            const username = refreshToken
              ? refreshTokens.get(refreshToken)
              : undefined;
            const user = users.find((item) => item.username === username);

            if (!user) {
              jsonResponse(
                response,
                errorResponse('Refresh token is invalid.'),
                403,
              );
              return;
            }

            const accessToken = `mock-access-${randomUUID()}`;
            accessTokens.set(accessToken, user.username);
            jsonResponse(response, accessToken);
            return;
          }

          if (pathname === '/api/auth/logout' && method === 'POST') {
            const refreshToken = getCookies(request).refresh_token;
            if (refreshToken) {
              refreshTokens.delete(refreshToken);
            }
            jsonResponse(response, success(''), 200, {
              'Set-Cookie':
                'refresh_token=; Path=/; Max-Age=0; HttpOnly; SameSite=Lax',
            });
            return;
          }

          if (
            (pathname === '/api/auth/codes' ||
              pathname === '/api/user/info' ||
              pathname === '/api/menu/all') &&
            method === 'GET'
          ) {
            const user = getUserFromAccessToken(request);
            if (!user) {
              jsonResponse(
                response,
                errorResponse('Unauthorized Exception'),
                401,
              );
              return;
            }

            if (pathname === '/api/auth/codes') {
              jsonResponse(response, success(getAccessCodes(user.username)));
              return;
            }

            if (pathname === '/api/user/info') {
              jsonResponse(response, success(publicUser(user)));
              return;
            }

            jsonResponse(response, success(getMenuList()));
            return;
          }

          if (
            mockPathname === '/mock/dicts/OpenStatusEnum' &&
            method === 'GET'
          ) {
            jsonResponse(
              response,
              success([
                {
                  value: '1',
                  label: '打开',
                  color: 'success',
                },
                {
                  value: '2',
                  label: '停止',
                  color: 'cyan',
                },
                {
                  value: '0',
                  label: '关闭',
                  color: 'red',
                },
              ]),
            );
            return;
          }

          const mock = viewMocks.find(
            (item: any) =>
              item.path === mockPathname &&
              item.method.toUpperCase() === method.toUpperCase(),
          );
          if (mock) {
            const params = getRequestParams(request);
            const body = method === 'GET' ? params : await readJsonBody(request);
            const result = await mock.handle({
              body,
              params,
            });
            if (result.code === 0) {
              jsonResponse(response, success(result.data));
            } else {
              jsonResponse(
                response,
                errorResponse(result.msg ?? 'Mock request failed.'),
                400,
              );
            }
            return;
          }

          next();
        } catch {
          jsonResponse(response, errorResponse('Invalid mock request.'), 400);
        }
      });
    },
  };
}
