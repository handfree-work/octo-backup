<script lang="ts" setup>
import { onMounted } from 'vue';

import { useFsAsync, useFsRef } from '@fast-crud/fast-crud';

import createCrudOptions from './crud';

const { crudRef, crudBinding, context } = useFsRef();

onMounted(async () => {
  const { crudExpose } = await useFsAsync({
    context,
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
    title="选择组件"
  >
    <fs-crud ref="crudRef" v-bind="crudBinding">
      <template #actionbar-right>
        <a-button @click="context.dynamicUpdateDictOptions">
          动态增加选项
        </a-button>
      </template>
    </fs-crud>
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
