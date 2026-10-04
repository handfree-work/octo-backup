<script setup lang="ts">
import { useRouter } from "vue-router";

defineOptions({ name: "SiteException403" });

const router = useRouter();
function toHome() {
  router.push({ path: "/" });
}
</script>

<template>
  <a-result status="403" title="403" sub-title="Sorry,  拒绝访问 .">
    <template #extra>
      <a-button type="primary" @click="toHome">返回首页（ Back Home ） </a-button>
    </template>
  </a-result>
</template>
