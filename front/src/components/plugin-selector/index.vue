<template>
  <div class="plugin-selector component-plugin-selector">
    <fs-table-select
      :model-value="modelValue"
      :dict="pluginDict"
      :create-crud-options="pluginCrudOptions"
      :dialog="{ title: dialogTitle }"
      :show-current="false"
      height="50vh"
      @update:model-value="onUpdateValue"
      @selected-change="onSelectedChange"
    >
      <template #default="{ open }">
        <a-tag v-if="selectedRecord" color="green">{{ selectedRecord.name || modelValue }}</a-tag>
        <span v-else class="plugin-selector__placeholder">{{ placeholder }}</span>
        <a-button v-if="selectedRecord" type="text" size="small" aria-label="清除授权" @click="clear">
          <CloseCircleOutlined />
        </a-button>
        <a-button class="plugin-selector__choose" @click="open">选择</a-button>
      </template>
    </fs-table-select>
  </div>
</template>

<script setup lang="ts">
defineOptions({
  name: "PluginSelector"
});

import { onMounted, ref, watch } from "vue";
import { CloseCircleOutlined } from "@ant-design/icons-vue";
import { dict, type CreateCrudOptionsProps, type CreateCrudOptionsRet } from "@fast-crud/fast-crud";
import { getPluginInstanceSimpleByIds } from "/src/views/sys/plugin/plugin-api";
import type { PluginInstance, PluginInstanceSimple, PluginMetadata } from "/src/views/sys/plugin/plugin-api";
import createPluginCrudOptions, { type PluginCrudContext } from "/src/views/sys/plugin/crud";
import { usePluginDefineStore } from "/src/store/modules/plugin-define";

const props = withDefaults(defineProps<{ modelValue?: string | number | null; pluginType: string; pluginName?: string; placeholder?: string; dialogTitle?: string }>(), {
  placeholder: "请选择",
  dialogTitle: "选择插件"
});
const emit = defineEmits<{ (event: "update:modelValue", value: number | null): void }>();
const selectedRecord = ref<PluginInstanceSimple>();
const pluginDefineStore = usePluginDefineStore();
const pluginDefine = ref<PluginMetadata>();
const pluginDict = dict({
  value: "id",
  label: "name",
  getNodesByValues: async (values: (string | number)[]) => {
    const ids = values.map(Number).filter((id) => Number.isSafeInteger(id) && id > 0);
    return ids.length ? getPluginInstanceSimpleByIds(ids) : [];
  }
});

onMounted(async () => {
  const value = props.modelValue;
  if (value == null || value === "") return;
  const nodes = await getPluginInstanceSimpleByIds([Number(value)]);
  if (props.modelValue === value) selectedRecord.value = nodes[0];
});

watch(
  () => props.modelValue,
  async (value) => {
    if (value == null || value === "") {
      selectedRecord.value = undefined;
      return;
    }
    const nodes = await getPluginInstanceSimpleByIds([Number(value)]);
    if (props.modelValue === value) selectedRecord.value = nodes[0];
  }
);

async function pluginCrudOptions(args: CreateCrudOptionsProps<PluginInstance, PluginCrudContext>): Promise<CreateCrudOptionsRet<PluginInstance>> {
  const metadata = (await pluginDefineStore.init()).filter((define) => define.type === props.pluginType && (!props.pluginName || define.name === props.pluginName));
  return createPluginCrudOptions({ ...args, context: { ...args.context, pluginType: props.pluginType, pluginName: props.pluginName, metadata } });
}

async function onSelectedChange(rows: PluginInstance[] = []) {
  selectedRecord.value = rows[0];
  pluginDefine.value = rows[0] ? await pluginDefineStore.getPluginDefine(rows[0].pluginName) : undefined;
}
function onUpdateValue(value: number | null) {
  emit("update:modelValue", value);
  if (value == null) selectedRecord.value = undefined;
}
function clear() {
  selectedRecord.value = undefined;
  pluginDefine.value = undefined;
  emit("update:modelValue", null);
}
</script>

<style lang="less">
.plugin-selector {
  .plugin-selector__placeholder {
    margin-right: 8px;
    color: #8c8c8c;
  }

  .plugin-selector__choose {
    margin-left: 8px;
  }

  .plugin-selector__helper {
    margin-top: 4px;
    color: #687386;
    font-size: 13px;
  }
}
</style>
