<template>
  <section class="about-page">
    <div class="about-header">
      <div class="brand-mark"><img :src="envRef.LOGO_PATH" alt="OctoBackup logo" /></div>
      <div>
        <div class="eyebrow">ABOUT OctoBackup</div>
        <h1>把 Restic 变成团队可用的备份服务</h1>
        <p>OctoBackup 是一个轻量、可自托管的远程备份管理系统，帮助你用清晰的界面管理备份基础设施。</p>
      </div>
    </div>
    <div class="about-grid">
      <article>
        <h2>项目定位</h2>
        <p>系统围绕 Restic 的快照与仓库能力构建，提供存储仓库管理、主机管理、备份任务编排与数据恢复入口。单实例运行，数据留在你的环境中。</p>
      </article>
      <article>
        <h2>设计原则</h2>
        <ul>
          <li>可控：配置透明，执行结果可追踪。</li>
          <li>可靠：依赖成熟的 Restic 备份格式。</li>
          <li>轻量：SQLite 与 Go 服务，部署简单。</li>
        </ul>
      </article>
      <article>
        <h2>开始使用</h2>
        <p>登录后从“存储仓库”开始，添加一个可用的 Local、SFTP、S3 或 MinIO 仓库，再按需配置备份任务。</p>
        <a-button type="primary" @click="router.push('/sys/repository')">前往存储仓库</a-button>
      </article>
    </div>
    <footer>OctoBackup · 远程备份管理 · v{{ version }}</footer>
  </section>
</template>
<script lang="ts" setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import { env } from "/@/utils/util.env";
const router = useRouter();
const envRef = ref(env);
const version = import.meta.env.VITE_APP_VERSION;
</script>
<style lang="less" scoped>
.about-page {
  min-height: 100%;
  padding: 54px clamp(24px, 7vw, 100px);
  background: #f5f8fc;
  color: #18283e;
}
.about-header {
  display: flex;
  align-items: center;
  gap: 28px;
  max-width: 960px;
  margin: 0 auto 48px;
}
.brand-mark {
  display: grid;
  place-items: center;
  flex: 0 0 94px;
  width: 94px;
  height: 94px;
  border-radius: 20px;
  background: #fff;
  box-shadow: 0 10px 30px #b8cde540;
}
.brand-mark img {
  width: 62px;
  height: 62px;
}
.eyebrow {
  color: #3178d4;
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.14em;
}
h1 {
  margin: 9px 0 12px;
  font-size: clamp(28px, 4vw, 44px);
  line-height: 1.2;
}
.about-header p {
  max-width: 650px;
  margin: 0;
  color: #65758a;
  font-size: 16px;
  line-height: 1.8;
}
.about-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  max-width: 960px;
  margin: 0 auto;
}
.about-grid article {
  padding: 25px;
  border: 1px solid #e0e8f1;
  border-radius: 10px;
  background: #fff;
}
.about-grid h2 {
  margin: 0 0 12px;
  font-size: 18px;
}
.about-grid p,
li {
  color: #65758a;
  font-size: 14px;
  line-height: 1.8;
}
.about-grid ul {
  padding-left: 19px;
  margin: 0;
}
.about-grid :deep(.ant-btn) {
  margin-top: 8px;
}
footer {
  max-width: 960px;
  margin: 42px auto 0;
  color: #9aa8b8;
  font-size: 12px;
}
@media (max-width: 760px) {
  .about-header {
    align-items: flex-start;
    flex-direction: column;
  }
  .about-grid {
    grid-template-columns: 1fr;
  }
}
</style>
