<template>
  <div class="plan-flow-node">
    <Handle v-if="data.kind === 'plan' || data.kind === 'repository'" id="left" type="target" :position="Position.Left" />
    <div class="flow-card-title">{{ data.title }}</div>
    <div class="flow-card-meta">{{ data.meta }}</div>
    <div v-if="data.status" class="flow-card-status">{{ data.status }}</div>
    <Handle v-if="data.kind === 'source' || data.kind === 'plan'" id="right" type="source" :position="Position.Right" />
  </div>
</template>
<script setup lang="ts">
import { Handle, Position } from "@vue-flow/core";
defineProps<{ data: { title: string; meta: string; status?: string; kind: "source" | "plan" | "repository" } }>();
</script>
<style lang="less">
.plan-flow-node {
  width: 220px;
  padding: 10px 14px;
  min-height: 64px;
  border-radius: 12px;
  color: #1c2b40;
  font-size: 13px;
  font-weight: 650;
  line-height: 1.65;
  white-space: pre-line;
  box-shadow: 0 10px 24px rgba(34, 55, 84, 0.12);
  .flow-card-title,
  .flow-card-meta,
  .flow-card-status {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  .flow-card-title {
    font-weight: 650;
  }
  .flow-card-meta {
    margin-top: 5px;
    color: #718096;
    font-size: 12px;
  }
  .flow-card-status {
    margin-top: 4px;
    color: #7c3aed;
    font-size: 11px;
  }
  :deep(.vue-flow__handle) {
    width: 10px;
    height: 10px;
    border: 2px solid #fff;
    box-shadow: 0 0 0 1px currentColor;
  }
  :deep(.flow-source) {
    border: 1px solid #b8d7ff;
    border-left: 4px solid #3b82f6;
    background: linear-gradient(135deg, #eaf3ff, #fff);
  }
  :deep(.flow-plan) {
    border: 1px solid #d5c6ff;
    border-left: 4px solid #8b5cf6;
    background: linear-gradient(135deg, #f2edff, #fff);
  }
  :deep(.flow-repository) {
    border: 1px solid #a9e8d2;
    border-left: 4px solid #10b981;
    background: linear-gradient(135deg, #e9fbf5, #fff);
  }
  :deep(.flow-source .vue-flow__handle-right) {
    color: #3b82f6;
    background: #3b82f6;
  }
  :deep(.flow-plan .vue-flow__handle) {
    color: #8b5cf6;
    background: #8b5cf6;
  }
  :deep(.flow-repository .vue-flow__handle-left) {
    color: #10b981;
    background: #10b981;
  }
}
</style>
