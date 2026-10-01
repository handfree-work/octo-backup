import { request } from "/src/api/service";

export interface Repository { id: string | number; name: string; repository: string; description?: string; createdAt?: number; updatedAt?: number }
export interface RepositoryInput { name: string; repository: string; password?: string; description?: string }
export interface RepositoryPage { offset: number; limit: number; records: Repository[]; total: number }
export const getRepositoryPage = (data: { offset?: number; limit?: number } = {}): Promise<RepositoryPage> => request({ url: "/repository/page", method: "post", data });
export const createRepository = (data: RepositoryInput): Promise<Repository> => request({ url: "/repository/create", method: "post", data });
export const updateRepository = (id: string | number, data: Partial<RepositoryInput>): Promise<Repository> => {
  const payload = { ...data };
  if (typeof payload.password === "string" && payload.password.trim() === "") delete payload.password;
  return request({ url: `/repository/update?id=${id}`, method: "post", data: payload });
};
export const deleteRepository = (id: string | number) => request({ url: `/repository/delete?id=${id}`, method: "post", data: {} });
