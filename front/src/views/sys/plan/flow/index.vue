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
      <VueFlow v-model:nodes="nodes" v-model:edges="edges" fit-view-on-init :node-types="nodeTypes" class="flow-canvas">
        <Background pattern-color="#dbe4f0" :gap="24" /><Controls position="bottom-left" /><MiniMap position="bottom-right" />
      </VueFlow>
    </div>
  </fs-page>
</template>
<script setup lang="ts">
import { ref } from "vue";
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
const { fitView } = useVueFlow();
const nodeTypes = { flow: PlanFlowNode };

function makeGraph() {
  const graphNodes: Node[] = [];
  const graphEdges: Edge[] = [];
  const centerPositionY = 330;
  const rowSpacing = 144;
  const getRowPositionY = (rowIndex: number, rowCount: number) => centerPositionY - ((rowCount - 1) * rowSpacing) / 2 + rowIndex * rowSpacing;
  const sourceY = new Map<number, number>();
  const repositoryY = new Map<number, number>();
  plans.value.forEach((plan, planIndex) => {
    const positionY = getRowPositionY(planIndex, plans.value.length);
    if (!sourceY.has(plan.sourceId)) {
      sourceY.set(plan.sourceId, positionY);
    }
    if (!repositoryY.has(plan.repositoryId)) {
      repositoryY.set(plan.repositoryId, positionY);
    }
  });
  const unusedSourceY = sources.value.filter((source) => !sourceY.has(Number(source.id)));
  const unusedRepositoryY = repositories.value.filter((repository) => !repositoryY.has(Number(repository.id)));
  unusedSourceY.forEach((source, sourceIndex) => {
    sourceY.set(Number(source.id), getRowPositionY(sourceIndex, unusedSourceY.length));
  });
  unusedRepositoryY.forEach((repository, repositoryIndex) => {
    repositoryY.set(Number(repository.id), getRowPositionY(repositoryIndex, unusedRepositoryY.length));
  });
  sources.value.forEach((source) => {
    graphNodes.push({
      id: `source-${source.id}`,
      type: "flow",
      position: { x: 48, y: sourceY.get(Number(source.id)) ?? centerPositionY },
      data: { title: source.name, meta: source.pluginName, kind: "source" },
      class: "flow-source"
    });
  });
  repositories.value.forEach((repository) => {
    graphNodes.push({
      id: `repository-${repository.id}`,
      type: "flow",
      position: { x: 790, y: repositoryY.get(Number(repository.id)) ?? centerPositionY },
      data: { title: repository.name, meta: repository.pluginName, kind: "repository" },
      class: "flow-repository"
    });
  });
  plans.value.forEach((plan, planIndex) => {
    const planPositionY = getRowPositionY(planIndex, plans.value.length);
    const sourceNodeId = `source-${plan.sourceId}`;
    const planNodeId = `plan-${plan.id}`;
    const repositoryNodeId = `repository-${plan.repositoryId}`;
    let planStatus = "已停用";
    let planClass = "flow-plan";
    let edgeClass = "edge-paused";
    if (plan.enabled) {
      planStatus = "运行中";
      planClass = "flow-plan active";
      edgeClass = "edge-active";
    }
    graphNodes.push({
      id: planNodeId,
      type: "flow",
      position: { x: 410, y: planPositionY },
      data: { title: plan.name, meta: `${plan.schedule} · ${plan.repoSubPath || "仓库根目录"}`, status: `${planStatus} · ${plan.lastStatus || "等待执行"}`, kind: "plan" },
      class: planClass
    });
    graphEdges.push(
      { id: `source-plan-${plan.id}`, source: sourceNodeId, sourceHandle: "right", target: planNodeId, targetHandle: "left", type: "bezier", animated: plan.enabled, class: edgeClass, markerEnd: "arrowclosed" },
      { id: `plan-repository-${plan.id}`, source: planNodeId, sourceHandle: "right", target: repositoryNodeId, targetHandle: "left", type: "bezier", animated: plan.enabled, class: edgeClass, markerEnd: "arrowclosed" }
    );
  });
  nodes.value = graphNodes.filter((node, nodeIndex, allNodes) => allNodes.findIndex((candidateNode) => candidateNode.id === node.id) === nodeIndex);
  edges.value = graphEdges;
  requestAnimationFrame(() => fitView({ padding: 0.3, minZoom: 0.65, maxZoom: 1.15 }));
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
    .flow-canvas {
      min-height: 0;
      flex: 1;
      background: #fff;
    }
  }
  .zone-bar {
    z-index: 4;
    padding: 10px;
    display: grid;
    grid-template-columns: 1fr 1fr 1fr;
    pointer-events: none;
    color: #364152;
    align-items: center;
    background: rgba(255, 255, 255, 0.94);
    border-bottom: 1px solid rgba(220, 229, 240, 0.78);
  }
  .zone-bar > div {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 15px;
  }
  .zone-bar > div:nth-child(2) {
    justify-content: center;
  }
  .zone-bar > div:nth-child(3) {
    justify-content: flex-end;
  }
  .zone-bar b {
    font-weight: 650;
  }
  .zone-bar :deep(.ant-btn) {
    pointer-events: auto;
    padding: 0 4px;
    color: #64748b;
  }
  :deep(.vue-flow__edge-path) {
    stroke-width: 2.5;
  }
  :deep(.edge-active .vue-flow__edge-path) {
    stroke: #8b5cf6;
    stroke-dasharray: 8 6;
    animation: dash 1s linear infinite;
  }
  :deep(.edge-paused .vue-flow__edge-path) {
    stroke: #a8b5c5;
    stroke-dasharray: 4 7;
  }
  :deep(.vue-flow__controls) {
    border: 1px solid #dce5f0;
    box-shadow: 0 5px 15px rgba(34, 55, 84, 0.1);
  }
  @keyframes dash {
    to {
      stroke-dashoffset: -28;
    }
  }
}
</style>
