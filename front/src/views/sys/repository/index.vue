<template>
  <PluginCrud plugin-type="repository" title="存储仓库" description="配置 Restic 使用的本地、SFTP、S3 和 MinIO 仓库。" />
</template>

<script setup lang="ts">
import PluginCrud from "/@/views/sys/plugin/index.vue";
</script>
