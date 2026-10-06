<template>
  <fs-page class="repository-files-page">
    <template #header
      ><div class="page-heading"><b>已备份文件</b><span>浏览 Restic 仓库中的快照和文件</span></div></template
    >
    <div class="empty-state">
      <Database :size="42" /><b>文件浏览接口尚未接入</b>
      <p>选择仓库后，这里将展示快照、目录和文件，并支持恢复操作。</p>
      <a-button type="primary" @click="goRepository">管理存储仓库</a-button>
    </div>
  </fs-page>
</template>
<script setup lang="ts">
import { Database } from "lucide-vue-next";
function goRepository() {
  window.location.hash = "#/sys/repository";
}
</script>
<style lang="less">
.repository-files-page {
  .page-heading {
    display: flex;
    gap: 12px;
    align-items: baseline;
    span {
      color: #8a94a6;
      font-size: 12px;
    }
  }
  .empty-state {
    min-height: 360px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 12px;
    color: #667085;
    border: 1px dashed #d0d5dd;
    border-radius: 10px;
    background: #fafbfc;
    p {
      margin: 0;
      color: #98a2b3;
    }
  }
}
</style>
