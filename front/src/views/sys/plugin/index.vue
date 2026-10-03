<template>
  <fs-page>
    <template #header>
      <div class="plugin-page-header">
        <div>
          <div class="title">{{ title }}</div>
          <div class="subtitle">{{ description }}</div>
        </div>
      </div>
    </template>
    <a-alert v-if="errorMessage" type="error" :message="errorMessage" show-icon closable @close="errorMessage = ''" />
    <PluginCrud v-if="metadata.length" :plugin-type="pluginType" :metadata="metadata" />
    <a-spin v-else :spinning="loading" />
  </fs-page>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import PluginCrud from "./components/PluginCrud.vue";
import { getPluginMetadata, type PluginMetadata } from "./api";

const props = withDefaults(defineProps<{ pluginType: string; title?: string; description?: string }>(), {
  title: "插件管理",
  description: "管理插件实例配置。"
});
const metadata = ref<PluginMetadata[]>([]);
const loading = ref(true);
const errorMessage = ref("");

onMounted(async () => {
  try {
    metadata.value = await getPluginMetadata({ type: props.pluginType });
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
