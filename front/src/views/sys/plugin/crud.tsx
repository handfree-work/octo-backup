import { AddReq, CreateCrudOptionsProps, CreateCrudOptionsRet, DelReq, dict, EditReq } from "@fast-crud/fast-crud";
import { createPlugin, deletePlugin, getPluginPage, type PluginField, type PluginInstance, type PluginMetadata, updatePlugin } from "./api";
import { collectPluginConfig } from "./config";
import set from "lodash-es/set";

export type PluginCrudContext = { pluginType: string; metadata: PluginMetadata[] };

function getComponent(field: PluginField): Record<string, unknown> {
  const configured = { ...(field.component || {}) } as Record<string, unknown>;
  return {
    vModel: "value",
    ...configured,
    name: configured.name || "a-input"
  };
}

export default function ({ context, crudExpose }: CreateCrudOptionsProps<PluginInstance, PluginCrudContext>): CreateCrudOptionsRet<PluginInstance> {
  const { metadata, pluginType } = context;
  const dynamicFieldKeys = new Set<string>();
  const columns: Record<string, any> = {
    id: { title: "ID", type: "text", form: { show: false }, column: { width: 80 } },
    name: { title: "名称", type: "text", order: -11, search: { show: true }, form: { order: -11, rules: [{ required: true, message: "请输入名称" }] } },
    pluginName: {
      title: "插件类型",
      type: "dict-select",
      order: -10,
      dict: dict({ data: metadata.map((item) => ({ value: item.name, label: item.title })) }),
      editForm: { component: { disabled: true } },
      form: {
        order: -10,
        rules: [{ required: true, message: "请选择插件类型" }],
        component: { showSearch: true },
        valueChange: {
          immediate: true,
          handle: ({ value, form, mode, immediate }: { value: string; form: Record<string, any>; mode: string; immediate: boolean }) => {
            buildDynamicFields(value, form, mode);
          }
        }
      }
    },
    createdAt: { title: "创建时间", type: "datetime", form: { show: false } }
  };

  function buildDynamicFields(pluginName: string, form: Record<string, any>, mode: string) {
    const formWrapper = crudExpose.getFormWrapperRef();
    const formColumns = formWrapper?.formOptions?.columns;
    if (!formColumns) {
      return;
    }
    console.log('crudBinding.value[mode + "Form"].columns', formColumns);
    const define = metadata.find((item) => item.name === pluginName);
    if (!define) {
      return;
    }
    (define.fields || []).forEach((field, index) => {
      const key = "config." + field.key;
      const component = getComponent(field);
      const fieldColumn: Record<string, any> = {
        title: field.title,
        type: field.type || "text",
        key,
        order: 10,
        ...field,
        component
      };
      if (field.options?.length && component.name === "a-select") {
        fieldColumn.form.component = { ...component, options: field.options };
      }
      formColumns[key] = fieldColumn;
      dynamicFieldKeys.add(key);
    });
    console.log('crudBinding.value[mode + "Form"].columns', formColumns);
  }

  return {
    crudOptions: {
      request: {
        pageRequest: (query: any) => {
          const pageSize = query.pageSize ?? query.limit ?? 20;
          const offset = query.offset ?? Math.max(0, ((query.currentPage ?? 1) - 1) * pageSize);
          return getPluginPage({ offset, limit: pageSize, pluginType, name: query.query?.name });
        },
        addRequest: async ({ form }: AddReq) => createPlugin({ name: form.name, description: form.description || "", pluginType, pluginName: form.pluginName, config: collectPluginConfig(form, metadata) }),
        editRequest: async ({ form, row }: EditReq) => updatePlugin(row.id, { name: form.name, description: form.description || "", config: collectPluginConfig(form, metadata) }),
        delRequest: async ({ row }: DelReq) => deletePlugin(row.id)
      },
      rowHandle: { fixed: "right" },
      columns: columns as any
    }
  };
}
