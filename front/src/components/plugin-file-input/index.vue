<template>
  <div class="plugin-file-input component-plugin-file-input">
    <a-button type="primary" @click="openFile">选择文件</a-button>
    <a-textarea class="plugin-file-input__textarea" :value="modelValue" :placeholder="placeholder" :auto-size="{ minRows: 3, maxRows: 8 }" @update:value="onInput" />
    <input ref="fileInput" class="plugin-file-input__native" type="file" @change="onFileChange" />
  </div>
</template>

<script setup lang="ts">
defineOptions({
  name: "PluginFileInput"
});

import { ref } from "vue";

const props = withDefaults(defineProps<{ modelValue?: string; placeholder?: string }>(), { placeholder: "选择文件或直接粘贴" });
const emit = defineEmits<{ (event: "update:modelValue", value: string): void }>();
const fileInput = ref<HTMLInputElement>();

function openFile() {
  fileInput.value?.click();
}
function onInput(value: string) {
  emit("update:modelValue", value);
}
function onFileChange(event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0];
  if (!file) return;
  const reader = new FileReader();
  reader.onload = () => emit("update:modelValue", String(reader.result || ""));
  reader.readAsText(file);
}
</script>

<style lang="less">
.plugin-file-input {
  display: flex;
  flex-direction: column;
  align-items: flex-start;

  .plugin-file-input__textarea {
    width: 100%;
    margin-top: 8px;
  }

  .plugin-file-input__native {
    position: absolute;
    width: 1px;
    height: 1px;
    overflow: hidden;
    clip: rect(0 0 0 0);
    clip-path: inset(50%);
    white-space: nowrap;
  }
}
</style>
