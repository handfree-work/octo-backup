<template>
  <div :class="['plan-flow-node', `flow-${data.kind}`, 'pointer']">
    <Handle v-if="data.kind === 'plan'" id="left" type="target" :position="Position.Left" />
    <div class="flow-card-title"><Folder v-if="data.kind === 'source'" /><CalendarClock v-if="data.kind === 'plan'" /><Database v-if="data.kind === 'repository'" />{{ data.title }}</div>
    <div class="flow-card-meta" :title="data.cron">{{ data.meta }}</div>
    <div v-if="data.kind === 'plan'" class="flow-card-status" :class="{ 'is-enabled': data.status === '已启用', 'is-disabled': data.status === '已禁用' }">{{ data.status }}</div>
    <button v-if="data.kind === 'plan'" type="button" class="flow-card-run" title="立即执行" @click.stop="data.onRun?.()">
      <Play :size="12" />
    </button>
    <Handle v-if="data.kind === 'source' || data.kind === 'plan'" id="right" type="source" :position="Position.Right" />
    <div v-if="data.kind === 'repository'" class="flow-card-paths">
      <div v-for="subPath in data.subPaths" :key="subPath.id" class="flow-card-path">
        <Handle :id="`path-${subPath.id}`" type="target" :position="Position.Left" />
        <span>{{ subPath.path }}</span>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
defineOptions({
  name: "PlanFlowNode"
});

import { Handle, Position } from "@vue-flow/core";
import { CalendarClock, Database, Folder } from "lucide-vue-next";
import { Play } from "lucide-vue-next";
defineProps<{ data: { title: string; meta: string; cron?: string; status?: string; kind: "source" | "plan" | "repository"; subPaths?: Array<{ id: number; path: string }>; onRun?: () => void } }>();
</script>
<style lang="less">
.plan-flow-node {
  width: 240px;
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

  .flow-card-meta {
    font-weight: 450;
    padding-right: 28px;
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
  .flow-card-paths {
    margin-top: 5px;
    border-top: 1px solid rgba(16, 185, 129, 0.18);
    padding-top: 4px;
  }
  .flow-card-path {
    position: relative;
    display: flex;
    align-items: center;
    min-height: 18px;
    padding-left: 8px;
    color: #047857;
    font-size: 10px;
    font-weight: 600;
  }
  .flow-card-path .vue-flow__handle {
    left: -16px;
    width: 7px;
    height: 7px;
    border-width: 1px;
    color: #10b981;
    background: #10b981;
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
  .flow-card-run {
    position: absolute;
    right: 9px;
    bottom: 7px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 22px;
    height: 22px;
    padding: 0;
    border: 0;
    border-radius: 6px;
    color: #7c3aed;
    background: rgba(139, 92, 246, 0.1);
    cursor: pointer;
  }
  .flow-card-run:hover {
    color: #fff;
    background: #7c3aed;
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
