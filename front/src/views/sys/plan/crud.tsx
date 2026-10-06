import { AddReq, CreateCrudOptionsRet, DelReq, EditReq, dict } from "@fast-crud/fast-crud";
import { getPluginInstancePage } from "/src/views/sys/plugin/plugin-api";
import { createBackupPlan, deleteBackupPlan, getBackupPlanInfo, getBackupPlanPage, updateBackupPlan, type BackupPlan } from "./api";
const instanceDict = (pluginType: string) =>
  dict({
    value: "id",
    label: "name",
    getData: async () => {
      const page = await getPluginInstancePage({ pluginType, offset: 0, limit: 1000 });
      return page.records;
    }
  });
export default function (): CreateCrudOptionsRet<BackupPlan> {
  return {
    crudOptions: {
      request: {
        pageRequest: ({ query }: any) => getBackupPlanPage({ offset: query.offset, limit: query.limit, name: query.query?.name }),
        infoRequest: ({ row, mode }: any) => {
          if (mode === "add") {
            return {};
          }
          return getBackupPlanInfo(row.id);
        },
        addRequest: ({ form }: AddReq<BackupPlan>) => createBackupPlan(form),
        editRequest: ({ form, row }: EditReq<BackupPlan>) => updateBackupPlan(row.id, form),
        delRequest: ({ row }: DelReq<BackupPlan>) => deleteBackupPlan(row.id)
      },
      form: {
        wrapper: {
          width: "900px"
        }
      },
      columns: {
        id: { title: "ID", type: "text", form: { show: false } },
        name: { title: "计划名称", type: "text", search: { show: true }, form: { helper: "备份计划名称", rules: [{ required: true, message: "请输入名称" }] } },
        sourceId: { title: "备份来源", type: "dict-select", dict: instanceDict("source"), form: { helper: "备份数据来源", rules: [{ required: true, message: "请选择备份来源" }] } },
        repositoryId: { title: "存储仓库", type: "dict-select", dict: instanceDict("repository"), form: { helper: "备份数据存储到哪个仓库", rules: [{ required: true, message: "请选择存储仓库" }] } },
        repoTag: { title: "仓库标签", type: "text", form: { component: { placeholder: "可选，用于区分同一仓库中的备份" }, helper: "使用 Restic 标签区分备份计划" } },
        schedule: {
          title: "Cron 调度",
          type: "text",
          form: {
            component: { name: "cron-editor", vModel: "modelValue", allowEveryMin: true },
            rules: [
              { required: true, message: "请选择 Cron 调度" },
              { pattern: /^\S+(\s+\S+){4}$/, message: "请输入 5 段式 Cron 表达式" }
            ]
          }
        },
        enabled: {
          title: "启用",
          type: "dict-switch",
          dict: dict({
            data: [
              { value: true, label: "启用" },
              { value: false, label: "禁用" }
            ]
          })
        },
        lastStatus: {
          title: "最近状态",
          type: "dict-select",
          dict: dict({
            data: [
              { value: "pending", label: "待执行" },
              { value: "running", label: "执行中" },
              { value: "success", label: "成功" },
              { value: "failed", label: "失败" }
            ]
          }),
          form: { show: false }
        },
        createdAt: { title: "创建时间", type: "datetime", form: { show: false } }
      } as any,
      rowHandle: { fixed: "right" }
    }
  };
}

