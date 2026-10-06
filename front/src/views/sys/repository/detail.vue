<template>
  <fs-page class="repository-detail-page">
    <a-page-header class="repository-detail-header" :title="repositoryName || '仓库详情'" @back="goBack">
      <template #subTitle>Restic 仓库快照</template>
      <template #extra><a-button :loading="loading" @click="loadSnapshots">刷新</a-button></template>
    </a-page-header>
    <a-alert v-if="errorMessage" type="error" show-icon :message="errorMessage" />
    <a-spin :spinning="loading" class="repository-detail-loading">
      <div class="repository-detail-content">
        <section class="snapshot-list-panel">
          <div class="panel-heading">
            <div>
              <span class="panel-kicker">REPOSITORY</span>
              <h2>Snapshots</h2>
            </div>
            <a-tag color="blue">{{ snapshots.length }} 个</a-tag>
          </div>
          <a-empty v-if="!snapshots.length && !loading" description="仓库中暂无快照" />
          <div v-else class="snapshot-list">
            <button v-for="snapshot in snapshots" :key="snapshot.id" class="snapshot-item" :class="{ active: selectedSnapshot?.id === snapshot.id }" type="button" @click="selectedSnapshot = snapshot">
              <span class="snapshot-status"><DatabaseOutlined /></span
              ><span class="snapshot-item-body"
                ><strong>{{ formatDate(snapshot.time) }}</strong
                ><span>{{ snapshot.hostname || "未知主机" }}</span
                ><span class="snapshot-id">{{ snapshot.id }}</span></span
              ><RightOutlined />
            </button>
          </div>
        </section>
        <section class="snapshot-detail-panel">
          <template v-if="selectedSnapshot">
            <div class="detail-heading">
              <div>
                <span class="panel-kicker">SNAPSHOT DETAIL</span>
                <h1>{{ formatDate(selectedSnapshot.time) }}</h1>
              </div>
              <a-tag color="green">已保存</a-tag>
            </div>
            <a-divider />
            <div class="detail-grid">
              <div>
                <span>Snapshot ID</span><strong>{{ selectedSnapshot.id }}</strong>
              </div>
              <div>
                <span>主机</span><strong>{{ selectedSnapshot.hostname || "未知" }}</strong>
              </div>
              <div>
                <span>用户</span><strong>{{ selectedSnapshot.username || "未知" }}</strong>
              </div>
              <div>
                <span>标签</span><strong>{{ formatTags(selectedSnapshot.tags) }}</strong>
              </div>
              <div>
                <span>文件数量</span><strong>{{ selectedSnapshot.summary?.total_files ?? "-" }}</strong>
              </div>
              <div>
                <span>处理字节</span><strong>{{ formatBytes(selectedSnapshot.summary?.total_bytes_processed) }}</strong>
              </div>
            </div>
            <div class="paths-section">
              <div class="section-title">备份路径</div>
              <a-list bordered size="small" :data-source="selectedSnapshot.paths || []"
                ><template #renderItem="{ item }"
                  ><a-list-item><FolderOutlined /> {{ item }}</a-list-item></template
                ></a-list
              >
            </div>
          </template>
          <a-empty v-else description="选择一个快照查看详情" />
        </section>
      </div>
    </a-spin>
  </fs-page>
</template>
<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import { DatabaseOutlined, FolderOutlined, RightOutlined } from "@ant-design/icons-vue";
import { getPluginInstanceInfo, getPluginInstanceSnapshots } from "/@/views/sys/plugin/plugin-api";
interface Snapshot {
  id: string;
  time?: string;
  hostname?: string;
  username?: string;
  tags?: string[];
  paths?: string[];
  summary?: { total_files?: number; total_bytes_processed?: number };
}
const route = useRoute();
const loading = ref(false);
const errorMessage = ref("");
const repositoryName = ref("");
const snapshots = ref<Snapshot[]>([]);
const selectedSnapshot = ref<Snapshot>();
function goBack() {
  window.location.hash = "#/sys/repository";
}
function formatDate(value?: string) {
  if (!value) return "未知时间";
  return new Date(value).toLocaleString("zh-CN", { hour12: false });
}
function formatTags(tags?: string[]) {
  return tags?.length ? tags.join(", ") : "无标签";
}
function formatBytes(value?: number) {
  if (!value) return "-";
  if (value < 1024 * 1024) return `${(value / 1024).toFixed(1)} KiB`;
  return `${(value / 1024 / 1024).toFixed(2)} MiB`;
}
async function loadSnapshots() {
  loading.value = true;
  errorMessage.value = "";
  try {
    const [info, result] = await Promise.all([getPluginInstanceInfo(String(route.params.id)), getPluginInstanceSnapshots(String(route.params.id))]);
    repositoryName.value = info.name;
    snapshots.value = Array.isArray(result) ? result : [];
    selectedSnapshot.value = snapshots.value[0];
  } catch (error: any) {
    errorMessage.value = error?.message || "加载仓库快照失败";
  } finally {
    loading.value = false;
  }
}
onMounted(loadSnapshots);
</script>
<style lang="less">
.repository-detail-page {
  background: #f5f7fa;
}
.repository-detail-header {
  background: #fff;
  border-bottom: 1px solid #e5e7eb;
}
.repository-detail-loading {
  display: block;
  min-height: 520px;
}
.repository-detail-content {
  display: grid;
  grid-template-columns: minmax(280px, 34%) 1fr;
  min-height: 560px;
  background: #fff;
}
.snapshot-list-panel {
  border-right: 1px solid #e5e7eb;
  padding: 24px 0;
}
.snapshot-detail-panel {
  padding: 32px 40px;
}
.panel-heading,
.detail-heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
}
.detail-heading {
  padding: 0;
}
.panel-kicker {
  color: #8b95a5;
  font-size: 11px;
  letter-spacing: 0.12em;
}
h1,
h2 {
  margin: 4px 0 0;
  color: #17212f;
}
h1 {
  font-size: 26px;
}
h2 {
  font-size: 18px;
}
.snapshot-list {
  margin-top: 18px;
}
.snapshot-item {
  display: flex;
  align-items: center;
  width: 100%;
  gap: 12px;
  padding: 15px 24px;
  border: 0;
  border-left: 3px solid transparent;
  background: transparent;
  color: #4b5563;
  text-align: left;
  cursor: pointer;
}
.snapshot-item:hover,
.snapshot-item.active {
  background: #eef6ff;
  border-left-color: #1677ff;
}
.snapshot-status {
  color: #1677ff;
  font-size: 17px;
}
.snapshot-item-body {
  display: grid;
  flex: 1;
  gap: 3px;
  min-width: 0;
}
.snapshot-item-body strong {
  color: #17212f;
}
.snapshot-item-body span {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.snapshot-id {
  color: #9aa4b2;
  font-family: monospace;
  font-size: 12px;
}
.detail-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
}
.detail-grid div {
  display: grid;
  gap: 6px;
}
.detail-grid span {
  color: #8b95a5;
  font-size: 13px;
}
.detail-grid strong {
  color: #263343;
  font-size: 15px;
  word-break: break-word;
}
.paths-section {
  margin-top: 36px;
}
.section-title {
  margin-bottom: 10px;
  color: #263343;
  font-weight: 600;
}
@media (max-width: 800px) {
  .repository-detail-content {
    grid-template-columns: 1fr;
  }
  .snapshot-list-panel {
    border-right: 0;
    border-bottom: 1px solid #e5e7eb;
  }
  .snapshot-detail-panel {
    padding: 24px;
  }
  .detail-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
