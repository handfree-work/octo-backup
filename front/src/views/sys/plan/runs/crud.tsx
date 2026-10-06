import { CreateCrudOptionsRet, dict } from "@fast-crud/fast-crud";
import { getBackupRunPage, type BackupLog } from "../api";
export default function (): CreateCrudOptionsRet<BackupLog> {
  return {
    crudOptions: {
      request: { pageRequest: ({ query }: any) => getBackupRunPage({ offset: query.offset, limit: query.limit, planId: query.query?.planId }) },
      rowHandle: { show: false },
      columns: {
        id: { title: "ID", type: "text" },
        planName: { title: "备份计划", type: "text", search: { show: true } },
        status: {
          title: "状态",
          type: "dict-select",
          dict: dict({
            data: [
              { value: "queued", label: "排队中", color: "default" },
              { value: "running", label: "执行中", color: "processing" },
              { value: "success", label: "成功", color: "success" },
              { value: "failed", label: "失败", color: "error" }
            ]
          })
        },
        progress: { title: "进度", type: "number", column: { component: { name: "a-progress", vModel: "percent", size: "small" } } },
        stage: { title: "阶段", type: "text" },
        startedAt: { title: "开始时间", type: "datetime" },
        finishedAt: { title: "结束时间", type: "datetime" },
        error: { title: "错误信息", type: "text", column: { ellipsis: true, tooltip: true, width: 240 } },
        result: { title: "执行结果", type: "text" }
      } as any
    }
  };
}
