<template>
  <fs-page>
    <template #header>
      <div class="title">备份计划 <span class="sub">管理备份来源、调度和存储仓库</span></div>
    </template>
    <fs-crud ref="crudRef" v-bind="crudBinding" />
  </fs-page>
</template>
<script setup lang="ts">
import { onMounted } from "vue";
import { useFs } from "@fast-crud/fast-crud";
import createCrudOptions from "./crud";
const { crudBinding, crudRef, crudExpose } = useFs({ createCrudOptions });
function goToFlow() {
  window.location.hash = "#/sys/plan/flow";
}
onMounted(() => crudExpose.doRefresh());
</script>
<style scoped>
.plan-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.sub {
  margin-left: 12px;
  color: #8490a1;
  font-size: 13px;
  font-weight: 400;
}
</style>
