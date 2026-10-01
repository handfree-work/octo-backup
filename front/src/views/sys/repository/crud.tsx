import { AddReq, CreateCrudOptionsProps, CreateCrudOptionsRet, DelReq, EditReq } from "@fast-crud/fast-crud";
import { createRepository, deleteRepository, getRepositoryPage, updateRepository } from "./api";

export default function (_props: CreateCrudOptionsProps): CreateCrudOptionsRet {
  return { crudOptions: { request: { pageRequest: getRepositoryPage, addRequest: ({ form }: AddReq) => createRepository(form), editRequest: ({ form, row }: EditReq) => updateRepository(row.id, form), delRequest: ({ row }: DelReq) => deleteRepository(row.id) }, rowHandle: { fixed: "right" }, columns: {
    id: { title: "ID", type: "text", form: { show: false }, column: { width: 80 } },
    name: { title: "仓库名称", type: "text", search: { show: true }, form: { rules: [{ required: true, message: "请输入仓库名称" }] } },
    repository: { title: "仓库地址", type: "text", form: { rules: [{ required: true, message: "请输入仓库地址" }] } },
    password: { title: "仓库密码", type: "text", column: { show: false }, form: { component: { showPassword: true }, helper: "编辑时留空表示不修改" } },
    description: { title: "备注", type: "text" },
    createdAt: { title: "创建时间", type: "datetime", form: { show: false } }
  } } };
}
