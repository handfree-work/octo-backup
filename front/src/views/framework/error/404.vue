<script setup lang="ts">
import { useRouter } from "vue-router";

defineOptions({ name: "SiteException404" });

const router = useRouter();
function toHome() {
  router.push({ path: "/" });
}
</script>
<template>
  <a-result status="404" title="404" sub-title="Sorry, the page you visited does not exist.">
    <template #extra>
      <a-button type="primary" @click="toHome">返回首页（ Back Home ） </a-button>
    </template>
  </a-result>
</template>
