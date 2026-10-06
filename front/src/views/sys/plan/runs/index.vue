<template>
  <fs-page class="backup-log-page">
    <template #header><div class="title">备份日志</div></template>
    <fs-crud ref="crudRef" v-bind="crudBinding" />
  </fs-page>
</template>
<script setup lang="ts">
import { onMounted } from "vue";
import { useFs } from "@fast-crud/fast-crud";
import createCrudOptions from "./crud";
const { crudBinding, crudRef, crudExpose } = useFs({ createCrudOptions });
onMounted(() => crudExpose.doRefresh());
defineOptions({ name: "BackupLog" });
</script>
<style lang="less">
.backup-log-page {
}
</style>
