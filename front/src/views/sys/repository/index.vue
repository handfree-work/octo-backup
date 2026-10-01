<template>
  <fs-page>
    <template #header><div class="title">存储仓库管理</div></template>
    <fs-crud ref="crudRef" v-bind="crudBinding" />
  </fs-page>
</template>
<script lang="ts">
import { defineComponent, onMounted } from "vue";
import { useFs } from "@fast-crud/fast-crud";
import createCrudOptions from "./crud";
export default defineComponent({ name: "StorageRepository", setup() { const { crudBinding, crudRef, crudExpose } = useFs({ createCrudOptions }); onMounted(() => crudExpose.doRefresh()); return { crudBinding, crudRef }; } });
</script>
