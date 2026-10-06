import { CreateCrudOptionsRet } from "@fast-crud/fast-crud";
import { getAuditLogPage, type AuditLog } from "./api";
export function createCrudOptions(): CreateCrudOptionsRet<AuditLog> {
  return {
    crudOptions: {
      request: { pageRequest: ({ query }: any) => getAuditLogPage({ offset: query.offset, limit: query.limit, username: query.query?.username, path: query.query?.path }) },
      rowHandle: { show: false },
      columns: {
        id: { title: "ID", type: "text" },
        operation: { title: "操作", type: "text" },
        username: { title: "用户", type: "text", search: { show: true } },
        method: { title: "方法", type: "text" },
        path: { title: "路径", type: "text", search: { show: true } },
        status: { title: "状态码", type: "number" },
        ip: { title: "IP 地址", type: "text" },
        duration: { title: "耗时(ms)", type: "number" },
        createdAt: { title: "时间", type: "datetime" }
      } as any
    }
  };
}
