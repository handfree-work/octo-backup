<template>
  <plugin-selector :model-value="modelValue" plugin-type="access" :plugin-name="pluginName" placeholder="请选择授权" dialog-title="选择授权提供者" @update:model-value="emit('update:modelValue', $event)" />
</template>

<script setup lang="ts">
defineOptions({
  name: "AccessSelector"
});

defineProps<{ modelValue?: string | number | null; pluginName?: string }>();
const emit = defineEmits<{ (event: "update:modelValue", value: number | null): void }>();
</script>
