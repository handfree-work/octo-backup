<template>
  <fs-page>
    <template #header>
      <div class="title">
        {{ title }}
        <span class="sub">{{ description }}</span>
      </div>
    </template>
    <a-alert v-if="errorMessage" type="error" :message="errorMessage" show-icon closable @close="errorMessage = ''" />
    <fs-crud v-if="metadata.length" ref="crudRef" v-bind="crudBinding" />
    <a-spin v-else :spinning="loading" />
  </fs-page>
</template>

<script setup lang="ts">
import { nextTick, onMounted, ref } from "vue";
import { useFsAsync, useFsRef } from "@fast-crud/fast-crud";
import createCrudOptions from "./crud";
import type { PluginMetadata } from "./plugin-api";
import { usePluginDefineStore } from "/src/store/modules/plugin-define";

const props = withDefaults(defineProps<{ pluginType: string; title?: string; description?: string }>(), {
  title: "插件管理",
  description: "管理插件实例配置。"
});
const metadata = ref<PluginMetadata[]>([]);
const loading = ref(true);
const errorMessage = ref("");
const { crudBinding, crudRef } = useFsRef();
const pluginDefineStore = usePluginDefineStore();

onMounted(async () => {
  try {
    const allDefines = await pluginDefineStore.init();
    metadata.value = allDefines.filter((define) => define.type === props.pluginType);
    await nextTick();
    const { crudExpose } = await useFsAsync({
      crudBinding,
      crudRef,
      createCrudOptions,
      context: { pluginType: props.pluginType, metadata: metadata.value }
    });
    await crudExpose.doRefresh();
  } catch (error: any) {
    errorMessage.value = error?.message || "加载插件定义失败";
  } finally {
    loading.value = false;
  }
});
</script>

<style scoped>
.plugin-page-header {
  display: flex;
  align-items: center;
}

.subtitle {
  margin-top: 4px;
  color: #687386;
  font-size: 13px;
}
</style>
