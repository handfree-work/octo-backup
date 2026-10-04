import { request } from "/src/api/service";
export interface BackupPlan {
  id: number;
  name: string;
  sourceId: number;
  repositoryId: number;
  repoSubPath?: string;
  schedule: string;
  enabled: boolean;
  lastStatus?: string;
  lastError?: string;
  lastRunAt?: number;
}
export const getBackupPlanPage = (data: Record<string, unknown> = {}) => request({ url: "/plan/page", method: "post", data });
export const getBackupPlanInfo = (id: number) => request({ url: `/plan/info?id=${id}`, method: "post", data: {} });
export const createBackupPlan = (data: Partial<BackupPlan>) => request({ url: "/plan/create", method: "post", data });
export const updateBackupPlan = (id: number, data: Partial<BackupPlan>) => request({ url: `/plan/update?id=${id}`, method: "post", data });
export const deleteBackupPlan = (id: number) => request({ url: `/plan/delete?id=${id}`, method: "post", data: {} });
