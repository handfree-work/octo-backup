<template>
  <plugin-selector :value="value" plugin-type="access" :plugin-name="pluginName" placeholder="请选择授权" dialog-title="选择授权提供者" @update:value="emit('update:value', $event)" />
</template>

<script setup lang="ts">
defineOptions({
  name: "AccessSelector"
});

defineProps<{ value?: string | number | null; pluginName?: string }>();
const emit = defineEmits<{ (event: "update:value", value: number | null): void }>();
</script>
