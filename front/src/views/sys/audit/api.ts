import { request } from "/src/api/service";
export interface AuditLog {
  id: number;
  userId: number;
  username: string;
  operation: string;
  method: string;
  path: string;
  status: number;
  ip: string;
  duration: number;
  createdAt: number;
}
export const getAuditLogPage = (data: Record<string, unknown> = {}) => request({ url: "/audit/page", method: "post", data });
