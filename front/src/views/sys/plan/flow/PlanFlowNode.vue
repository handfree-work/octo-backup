<template>
  <div :class="['plan-flow-node', `flow-${data.kind}`]">
    <Handle v-if="data.kind === 'plan' || data.kind === 'repository'" id="left" type="target" :position="Position.Left" />
    <div class="flow-card-title"><Server v-if="data.kind === 'source'" /><CalendarClock v-if="data.kind === 'plan'" /><Database v-if="data.kind === 'repository'" />{{ data.title }}</div>
    <div class="flow-card-meta">{{ data.meta }}</div>
    <div v-if="data.kind === 'plan'" class="flow-card-status" :class="{ 'is-enabled': data.status === '已启用', 'is-disabled': data.status === '已禁用' }">{{ data.status }}</div>
    <Handle v-if="data.kind === 'source' || data.kind === 'plan'" id="right" type="source" :position="Position.Right" />
  </div>
</template>
<script setup lang="ts">
defineOptions({
  name: "PlanFlowNode"
});

import { Handle, Position } from "@vue-flow/core";
import { CalendarClock, Database, Server } from "lucide-vue-next";
defineProps<{ data: { title: string; meta: string; status?: string; kind: "source" | "plan" | "repository" } }>();
</script>
<style lang="less">
.plan-flow-node {
  width: 188px;
  padding: 7px 10px;
  min-height: 50px;
  border-radius: 12px;
  color: #1c2b40;
  font-size: 12px;
  font-weight: 650;
  line-height: 1.65;
  white-space: pre-line;
  box-shadow: 0 10px 24px rgba(34, 55, 84, 0.12);
  transition-property: box-shadow, border-color;
  transition-duration: 160ms;
  transition-timing-function: ease-out;
  &:hover {
    box-shadow: 0 12px 26px rgba(34, 55, 84, 0.15);
  }
  .flow-card-title,
  .flow-card-meta,
  .flow-card-status {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  .flow-card-title {
    display: flex;
    align-items: center;
    gap: 7px;
    padding-right: 48px;
    font-weight: 650;
    svg {
      width: 15px;
      height: 15px;
      flex: 0 0 auto;
    }
  }
  &.flow-source .flow-card-title svg {
    color: #2563eb;
  }
  &.flow-plan .flow-card-title svg {
    color: #7c3aed;
  }
  &.flow-repository .flow-card-title svg {
    color: #059669;
  }
  .flow-card-meta {
    margin-top: 3px;
    color: #718096;
    font-size: 11px;
  }
  .flow-card-status {
    position: absolute;
    top: 7px;
    right: 9px;
    display: inline-flex;
    padding: 1px 7px;
    border-radius: 999px;
    color: #6d28d9;
    background: rgba(139, 92, 246, 0.1);
    font-size: 11px;
    line-height: 1.5;
    pointer-events: none;
  }
  .flow-card-status.is-enabled {
    color: #15803d;
    background: rgba(34, 197, 94, 0.12);
  }
  .flow-card-status.is-disabled {
    color: #64748b;
    background: rgba(100, 116, 139, 0.12);
  }
  .vue-flow__handle {
    width: 10px;
    height: 10px;
    border: 2px solid #fff;
    box-shadow: 0 0 0 1px currentColor;
  }
  &.flow-source {
    border: 1px solid #b8d7ff;
    border-left: 4px solid #3b82f6;
    background: linear-gradient(135deg, #eaf3ff, #fff);
  }
  &.flow-plan {
    border: 1px solid #d5c6ff;
    border-left: 4px solid #8b5cf6;
    background: linear-gradient(135deg, #f2edff, #fff);
  }
  &.flow-repository {
    border: 1px solid #a9e8d2;
    border-left: 4px solid #10b981;
    background: linear-gradient(135deg, #e9fbf5, #fff);
  }
  &.flow-source .vue-flow__handle-right {
    color: #3b82f6;
    background: #3b82f6;
  }
  &.flow-plan .vue-flow__handle {
    color: #8b5cf6;
    background: #8b5cf6;
  }
  &.flow-repository .vue-flow__handle-left {
    color: #10b981;
    background: #10b981;
  }
}
</style>
