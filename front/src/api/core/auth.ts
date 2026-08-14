import { baseRequestClient } from '#/api/request';

export namespace AuthApi {
  export interface LoginParams {
    password: string;
    username: string;
  }

  export interface RegisterParams extends LoginParams {
    nickName?: string;
  }

  export interface User {
    avatar?: string;
    id: number;
    nickName: string;
    role: 'admin' | 'read' | 'write';
    username: string;
  }

  export interface LoginResult {
    expiresAt: number;
    token: string;
    user: User;
  }

  export interface RawResponse<T> {
    data: {
      data: T;
    };
  }
}

export async function loginApi(data: AuthApi.LoginParams) {
  const response = await baseRequestClient.post<AuthApi.RawResponse<AuthApi.LoginResult>>('/auth/login', data);
  return response.data.data;
}

export async function registerApi(data: AuthApi.RegisterParams) {
  const response = await baseRequestClient.post<AuthApi.RawResponse<AuthApi.User>>('/auth/register', data);
  return response.data.data;
}
