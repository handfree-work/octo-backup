import { request } from "../service";
/**
 * @description: Login interface parameters
 */
export interface LoginReq {
  username: string;
  password: string;
}

export interface UserInfoRes {
  id: string | number;
  username: string;
  nickName: string;
  avatar?: string;
  role: "admin" | "write" | "read";
}

export interface LoginRes {
  token: string;
  expiresAt: number;
  user: UserInfoRes;
}

export async function login(data: LoginReq): Promise<LoginRes> {
  return await request({
    url: "/auth/login",
    method: "post",
    data
  });
}

export interface RegisterReq {
  username: string;
  password: string;
}

export interface RegisterRes {
  id: string | number;
  username: string;
  role: "admin" | "write" | "read";
}

export async function register(data: RegisterReq): Promise<RegisterRes> {
  return await request({
    url: "/auth/register",
    method: "post",
    data
  });
}
