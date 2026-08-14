<script lang="ts" setup>
import { onMounted } from 'vue';

import { useFsAsync, useFsRef } from '@fast-crud/fast-crud';

import createCrudOptions from './crud';

const { crudRef, crudBinding } = useFsRef();

onMounted(async () => {
  const { crudExpose } = await useFsAsync({
    crudBinding,
    crudRef,
    createCrudOptions,
  });
  await crudExpose.doRefresh();
});
</script>

<template>
  <Page
    auto-content-height
    content-class="fast-crud-demo-page h-full bg-card p-0"
    title="文本组件"
  >
    <fs-crud ref="crudRef" v-bind="crudBinding" />
  </Page>
</template>

<style scoped>
:deep(.fast-crud-demo-page) {
  background: hsl(var(--card));
}

:deep(.fast-crud-demo-page .fs-crud-container) {
  min-height: 100%;
  background: hsl(var(--card));
}
</style>
