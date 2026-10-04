<template>
  <PluginCrud plugin-type="access" title="授权配置" description="集中管理主机、云厂商和对象存储的连接授权。" />
</template>

<script setup lang="ts">
defineOptions({
  name: "AccessManager"
});

import PluginCrud from "/@/views/sys/plugin/index.vue";
</script>
