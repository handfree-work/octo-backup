import { beforeEach, describe, expect, it, vi } from 'vitest';

const { post } = vi.hoisted(() => ({ post: vi.fn() }));

vi.mock('#/api/request', () => ({
  baseRequestClient: { post },
  requestClient: { post },
}));

import { loginApi, registerApi } from './auth';

describe('认证 API', () => {
  beforeEach(() => {
    post.mockReset();
  });

  it('解包登录响应并返回 JWT 与用户信息', async () => {
    const result = {
      expiresAt: 1_786_285_124,
      token: 'jwt-token',
      user: { id: 1, nickName: '管理员', role: 'admin', username: 'admin' },
    };
    post.mockResolvedValue({ data: { data: result } });

    await expect(loginApi({ password: 'secret', username: 'admin' })).resolves.toEqual(result);
    expect(post).toHaveBeenCalledWith('/auth/login', {
      password: 'secret',
      username: 'admin',
    });
  });

  it('提交注册信息并解包新用户', async () => {
    const user = { id: 2, role: 'read', username: 'reader' };
    post.mockResolvedValue({ data: { data: user } });

    await expect(
      registerApi({
        nickName: '只读用户',
        password: 'secret',
        username: 'reader',
      }),
    ).resolves.toEqual(user);
    expect(post).toHaveBeenCalledWith('/auth/register', {
      nickName: '只读用户',
      password: 'secret',
      username: 'reader',
    });
  });
});
