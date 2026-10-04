<template>
  <fs-page class="plan-flow-page">
    <template #header>
      <div class="title">
        备份计划
        <span class="sub">编排来源、计划与存储仓库之间的数据流</span>
      </div>
    </template>
    <div class="canvas-shell">
      <div class="zone-bar">
        <div><b>来源</b><a-button type="link" size="small" @click="goTo('source')">＋ 添加</a-button></div>
        <div><b>计划</b><a-button type="link" size="small" @click="openPlanAdd">＋ 添加</a-button></div>
        <div><b>仓库</b><a-button type="link" size="small" @click="goTo('repository')">＋ 添加</a-button></div>
      </div>
      <VueFlow v-model:nodes="nodes" v-model:edges="edges" :node-types="nodeTypes" class="flow-canvas">
        <Background pattern-color="#dbe4f0" :gap="24" /><Controls position="bottom-left" /><MiniMap position="bottom-right" />
      </VueFlow>
    </div>
  </fs-page>
</template>
<script setup lang="ts">
defineOptions({
  name: "BackupPlanFlow"
});

import { nextTick, ref } from "vue";
import { VueFlow, useVueFlow, type Edge, type Node } from "@vue-flow/core";
import { Background } from "@vue-flow/background";
import { Controls } from "@vue-flow/controls";
import { MiniMap } from "@vue-flow/minimap";
import PlanFlowNode from "./PlanFlowNode.vue";
import "@vue-flow/core/dist/style.css";
import "@vue-flow/core/dist/theme-default.css";
import { getBackupPlanPage, type BackupPlan } from "../api";
import { getPluginInstancePage, type PluginInstanceSimple } from "../../plugin/plugin-api";
const plans = ref<BackupPlan[]>([]);
const sources = ref<PluginInstanceSimple[]>([]);
const repositories = ref<PluginInstanceSimple[]>([]);
const nodes = ref<Node[]>([]);
const edges = ref<Edge[]>([]);
const { setCenter } = useVueFlow();
const nodeTypes = { flow: PlanFlowNode };

function makeGraph() {
  const graphNodes: Node[] = [];
  const graphEdges: Edge[] = [];
  const centerPositionY = 330;
  const rowSpacing = 100;
  const getRowPositionY = (rowIndex: number, rowCount: number) => centerPositionY - ((rowCount - 1) * rowSpacing) / 2 + rowIndex * rowSpacing;
  sources.value.forEach((source) => {
    const sourceIndex = sources.value.indexOf(source);
    graphNodes.push({
      id: `source-${source.id}`,
      type: "flow",
      position: { x: 48, y: getRowPositionY(sourceIndex, sources.value.length) },
      data: { title: source.name, meta: source.pluginName, kind: "source" },
      class: "flow-source"
    });
  });
  repositories.value.forEach((repository) => {
    const repositoryIndex = repositories.value.indexOf(repository);
    graphNodes.push({
      id: `repository-${repository.id}`,
      type: "flow",
      position: { x: 790, y: getRowPositionY(repositoryIndex, repositories.value.length) },
      data: { title: repository.name, meta: repository.pluginName, kind: "repository" },
      class: "flow-repository"
    });
  });
  plans.value.forEach((plan, planIndex) => {
    const planPositionY = getRowPositionY(planIndex, plans.value.length);
    const sourceNodeId = `source-${plan.sourceId}`;
    const planNodeId = `plan-${plan.id}`;
    const repositoryNodeId = `repository-${plan.repositoryId}`;
    let planClass = "flow-plan";
    let edgeClass = "edge-paused";
    if (plan.enabled) {
      planClass = "flow-plan active";
      edgeClass = "edge-active";
    }
    graphNodes.push({
      id: planNodeId,
      type: "flow",
      position: { x: 410, y: planPositionY },
      data: {
        title: plan.name,
        meta: `${plan.schedule} · ${plan.repoSubPath || "仓库根目录"}`,
        status: plan.enabled ? "已启用" : "已禁用",
        kind: "plan"
      },
      class: planClass
    });
    graphEdges.push(
      { id: `source-plan-${plan.id}`, source: sourceNodeId, sourceHandle: "right", target: planNodeId, targetHandle: "left", type: "bezier", animated: plan.enabled, class: edgeClass, markerEnd: "arrowclosed" },
      { id: `plan-repository-${plan.id}`, source: planNodeId, sourceHandle: "right", target: repositoryNodeId, targetHandle: "left", type: "bezier", animated: plan.enabled, class: edgeClass, markerEnd: "arrowclosed" }
    );
  });
  nodes.value = graphNodes.filter((node, nodeIndex, allNodes) => allNodes.findIndex((candidateNode) => candidateNode.id === node.id) === nodeIndex);
  edges.value = graphEdges;
  nextTick(() => {
    requestAnimationFrame(() => {
      const nodeWidth = 188;
      const nodeHeight = 50;
      const minimumX = Math.min(...nodes.value.map((node) => node.position.x));
      const maximumX = Math.max(...nodes.value.map((node) => node.position.x + nodeWidth));
      const minimumY = Math.min(...nodes.value.map((node) => node.position.y));
      const maximumY = Math.max(...nodes.value.map((node) => node.position.y + nodeHeight));
      const graphCenterX = (minimumX + maximumX) / 2;
      const graphCenterY = (minimumY + maximumY) / 2;
      setCenter(graphCenterX, graphCenterY, { zoom: 1.1 });
    });
  });
}
function goTo(name: string) {
  window.location.hash = `#/sys/${name}`;
}
function openPlanAdd() {
  window.location.hash = "#/sys/plan";
}
async function loadFlow() {
  const [planPage, sourcePage, repositoryPage] = await Promise.all([
    getBackupPlanPage({ offset: 0, limit: 1000 }),
    getPluginInstancePage({ pluginType: "source", offset: 0, limit: 1000 }),
    getPluginInstancePage({ pluginType: "repository", offset: 0, limit: 1000 })
  ]);
  plans.value = planPage.records || [];
  sources.value = sourcePage.records || [];
  repositories.value = repositoryPage.records || [];
  makeGraph();
}
loadFlow();
</script>
<style lang="less">
.plan-flow-page {
  .canvas-shell {
    height: 100%;
    border: 1px solid #dce5f0;
    border-radius: 0;
    background: #fff;
    display: flex;
    flex-direction: column;
    position: relative;
    overflow: hidden;
    .flow-canvas {
      min-height: 0;
      flex: 1;
      background: #fff;
      position: relative;
      z-index: 1;
    }
  }
  .canvas-shell::before {
    content: "";
    position: absolute;
    inset: 52px 0 0;
    pointer-events: none;
    background: linear-gradient(90deg, rgba(59, 130, 246, 0.025) 0 33.33%, rgba(139, 92, 246, 0.025) 33.33% 66.66%, rgba(16, 185, 129, 0.025) 66.66% 100%);
    z-index: 0;
  }
  .zone-bar {
    z-index: 4;
    padding: 10px;
    display: flex;
    justify-content: space-evenly;
    pointer-events: none;
    color: #364152;
    align-items: center;
    background: rgba(255, 255, 255, 0.94);
    border-bottom: 1px solid rgba(220, 229, 240, 0.78);
    box-shadow: 0 1px 0 rgba(255, 255, 255, 0.8);
  }
  .zone-bar > div {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 15px;
    min-height: 28px;
    padding: 0 14px;
    border-radius: 8px;
  }
  .zone-bar b {
    font-weight: 650;
    letter-spacing: 0.02em;
  }
  .zone-bar :deep(.ant-btn) {
    pointer-events: auto;
    padding: 0 4px;
    color: #64748b;
  }
  :deep(.vue-flow__edge-path) {
    stroke-width: 2;
    stroke-linecap: round;
  }
  :deep(.edge-active .vue-flow__edge-path) {
    stroke: #8b5cf6;
    stroke-dasharray: 8 6;
    animation: dash 1s linear infinite;
  }
  :deep(.edge-paused .vue-flow__edge-path) {
    stroke: #b7c0cd;
    stroke-dasharray: 4 7;
  }
  :deep(.vue-flow__controls) {
    border: 1px solid #dce5f0;
    box-shadow: 0 5px 15px rgba(34, 55, 84, 0.1);
    border-radius: 10px;
    overflow: hidden;
  }
  @keyframes dash {
    to {
      stroke-dashoffset: -28;
    }
  }
}
</style>
