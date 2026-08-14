<script lang="ts" setup>
import { onMounted } from 'vue';

import { dict, useFsAsync, useFsRef } from '@fast-crud/fast-crud';
import type {
  AddReq,
  CreateCrudOptionsProps,
  CreateCrudOptionsRet,
  DelReq,
  EditReq,
  UserPageQuery,
  UserPageRes,
} from '@fast-crud/fast-crud';
import { cloneDeep, find, maxBy, merge, remove } from 'lodash-es';

import { Page } from '@vben/common-ui';

const records = [{ id: 1, name: 'Hello World', type: 1 }];

const createCrudOptions = async function (
  _props: CreateCrudOptionsProps,
): Promise<CreateCrudOptionsRet> {
  const pageRequest = async (_query: UserPageQuery): Promise<UserPageRes> => {
    return {
      records: cloneDeep(records),
      offset: 0,
      limit: 20,
      total: records.length,
    };
  };

  const editRequest = async ({ form, row }: EditReq) => {
    const target = find(records, (item) => item.id === row.id);
    merge(target, form);
    return target;
  };

  const delRequest = async ({ row }: DelReq) => {
    remove(records, (item) => item.id === row.id);
  };

  const addRequest = async ({ form }: AddReq) => {
    const maxRecord = maxBy(records, (item) => item.id);
    form.id = (maxRecord?.id || 0) + 1;
    records.push(form as (typeof records)[number]);
    return form;
  };

  return {
    crudOptions: {
      request: {
        pageRequest,
        addRequest,
        editRequest,
        delRequest,
      },
      columns: {
        name: {
          title: '姓名',
          type: 'text',
          search: { show: true },
          column: {
            resizable: true,
            width: 200,
          },
        },
        type: {
          title: '类型',
          type: 'dict-select',
          dict: dict({
            data: [
              { value: 1, label: '开始' },
              { value: 0, label: '停止' },
            ],
          }),
        },
      },
    },
  };
};

const { crudRef, crudBinding } = useFsRef();

onMounted(async () => {
  const { crudExpose } = await useFsAsync({
    crudRef,
    crudBinding,
    createCrudOptions,
  });
  await crudExpose.doRefresh();
});
</script>

<template>
  <Page
    auto-content-height
    content-class="fast-crud-demo-page h-full bg-card p-0"
    title="Fast CRUD HelloWorld"
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
