import { request } from "/src/api/service";
import SHA256 from "crypto-js/sha256";
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
  createdAt?: number;
  updatedAt?: number;
}
export interface LoginRes {
  token: string;
  expiresAt: number;
  user: UserInfoRes;
}
export interface RegisterReq {
  username: string;
  password: string;
}
export interface RegisterRes {
  id: string | number;
  username: string;
  role: UserInfoRes["role"];
}
export const hashPassword = (password: string): string => SHA256(password).toString();
export const login = (data: LoginReq): Promise<LoginRes> => request({ url: "/auth/login", method: "post", data: { ...data, password: hashPassword(data.password) } });
export const register = (data: RegisterReq): Promise<RegisterRes> => request({ url: "/auth/register", method: "post", data: { ...data, password: hashPassword(data.password) } });
