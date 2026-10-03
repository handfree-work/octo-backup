<template>
  <fs-crud ref="crudRef" v-bind="crudBinding" />
</template>

<script setup lang="ts">
import { onMounted } from "vue";
import { useFs } from "@fast-crud/fast-crud";
import createCrudOptions from "../crud";
import type { PluginInstance, PluginMetadata } from "../api";

const props = defineProps<{ pluginType: string; metadata: PluginMetadata[] }>();
const context = { pluginType: props.pluginType, metadata: props.metadata };
const { crudBinding, crudRef, crudExpose } = useFs<PluginInstance, typeof context>({ createCrudOptions, context });

onMounted(() => crudExpose.doRefresh());
</script>
