<template>
  <div class="plugin-path-selector component-plugin-path-selector">
    <div class="plugin-path-selector__actions">
      <a-button :disabled="!instanceId" @click="openPicker">选择目录</a-button>
      <span v-if="!instanceId" class="plugin-path-selector__hint">保存来源后可浏览 SSH 目录</span>
    </div>
    <a-textarea :value="formatPathInput(value)" :placeholder="placeholder" :auto-size="{ minRows: 3, maxRows: 8 }" @update:value="onInput" />
    <a-modal v-model:open="pickerOpen" title="选择备份目录" :confirm-loading="loading" @ok="confirmSelection">
      <a-alert v-if="errorMessage" type="error" :message="errorMessage" show-icon />
      <a-spin :spinning="loading">
        <a-tree v-model:checked-keys="checkedKeys" checkable :tree-data="treeData" :field-names="treeFieldNames" />
      </a-spin>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
defineOptions({
  name: "PluginPathSelector"
});

import { ref } from "vue";
import { executePluginInstanceAction } from "/src/views/sys/plugin/plugin-api";
import { formatPathInput, normalizePathNodes, parsePathInput, type PluginPathNode } from "./helpers";

const props = withDefaults(defineProps<{ value?: string | string[]; placeholder?: string; instanceId?: string | number; action?: string; accessId?: string | number }>(), {
  placeholder: "选择目录或粘贴路径，每行一个",
  action: "onListPaths"
});
const emit = defineEmits<{ (event: "update:value", value: string[]): void }>();
const pickerOpen = ref(false);
const loading = ref(false);
const errorMessage = ref("");
const treeData = ref<PluginPathNode[]>([]);
const checkedKeys = ref<string[]>([]);
const treeFieldNames = { title: "title", key: "key", children: "children" };

function onInput(value: string) {
  emit("update:value", parsePathInput(value));
}
async function openPicker() {
  if (!props.instanceId) return;
  pickerOpen.value = true;
  loading.value = true;
  errorMessage.value = "";
  try {
    const result = await executePluginInstanceAction(props.instanceId, props.action, { accessId: props.accessId });
    treeData.value = normalizePathNodes(result);
    checkedKeys.value = parsePathInput(props.value);
  } catch (error: any) {
    errorMessage.value = error?.message || "读取 SSH 目录失败";
  } finally {
    loading.value = false;
  }
}
function confirmSelection() {
  emit("update:value", parsePathInput(checkedKeys.value));
  pickerOpen.value = false;
}
</script>

<style lang="less">
.plugin-path-selector {
  display: grid;
  gap: 8px;
  .plugin-path-selector__actions {
    display: flex;
    align-items: center;
    gap: 8px;
  }
  .plugin-path-selector__hint {
    color: #8c8c8c;
    font-size: 12px;
  }
}
</style>
