import { request } from "/src/api/service";
import { hashPassword } from "/@/views/framework/auth/api";

export interface UserInfoRes {
  id: string | number;
  username: string;
  nickName: string;
  avatar?: string;
  role: "admin" | "write" | "read";
  createdAt?: number;
  updatedAt?: number;
}
export interface CreateUserReq {
  username: string;
  password: string;
  nickName?: string;
  role: UserInfoRes["role"];
}
export interface UpdateUserReq {
  username?: string;
  password?: string;
  nickName?: string;
  role?: UserInfoRes["role"];
}
export interface UserListQuery {
  offset?: number;
  limit?: number;
}
export interface UserListRes {
  offset: number;
  limit: number;
  records: UserInfoRes[];
  total: number;
}

export const getUserList = (data: UserListQuery = {}): Promise<UserListRes> => request({ url: "/user/page", method: "post", data });
export const createUser = (data: CreateUserReq): Promise<UserInfoRes> =>
  request({ url: "/user/create", method: "post", data: { ...data, password: hashPassword(data.password) } });
export const updateUser = (id: string | number, data: UpdateUserReq): Promise<UserInfoRes> => {
  const payload = { ...data };
  if (typeof payload.password === "string") {
    if (payload.password.trim() === "") delete payload.password;
    else payload.password = hashPassword(payload.password);
  }
  return request({ url: `/user/update?id=${id}`, method: "post", data: payload });
};
export const deleteUser = (id: string | number) => request({ url: `/user/delete?id=${id}`, method: "post", data: {} });
