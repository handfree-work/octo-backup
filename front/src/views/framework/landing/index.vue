<template>
  <div class="landing-page">
    <nav class="landing-nav">
      <div class="nav-inner">
        <div class="brand"><img :src="envRef.LOGO_PATH" alt="OctoBackup" /><span>OctoBackup</span></div>
        <div class="nav-actions"><a href="#features">核心能力</a><router-link class="btn btn-outline" to="/login">登录</router-link><router-link class="btn btn-primary" to="/register">开始使用</router-link></div>
      </div>
    </nav>
    <main>
      <section class="hero">
        <div class="hero-content">
          <div class="eyebrow">OPEN SOURCE · SELF HOSTED</div>
          <h1>让备份真正成为<br /><strong>可以依赖的基础设施</strong></h1>
          <p>OctoBackup 是一个基于 Restic 的远程备份管理系统。用一个清晰、轻量的控制台，管理仓库、主机、任务和恢复。</p>
          <div class="hero-actions">
            <router-link class="btn btn-large btn-primary" to="/login">进入控制台</router-link><a class="btn btn-large btn-quiet" href="#features">了解能力 <span>↓</span></a>
          </div>
          <div class="trust-row">
            <span><b>✓</b> 数据自托管</span><span><b>✓</b> 支持多种仓库</span><span><b>✓</b> 基于 Restic</span>
          </div>
        </div>
        <div class="hero-visual">
          <div class="visual-grid"></div>
          <div class="backup-panel">
            <div class="panel-top"><span class="panel-dot"></span><span>BACKUP OVERVIEW</span><span class="panel-live">LIVE</span></div>
            <div class="panel-number">99.9<small>%</small></div>
            <div class="panel-caption">最近 30 天备份成功率</div>
            <div class="panel-bars"><i></i><i></i><i></i><i></i><i></i><i></i><i></i><i></i><i></i><i></i><i></i><i></i></div>
            <div class="panel-footer"><span>12 个快照</span><span class="green">运行正常</span></div>
          </div>
        </div>
      </section>
      <section id="features" class="features">
        <div class="section-heading">
          <div class="eyebrow">BUILT FOR RECOVERY</div>
          <h2>从备份到恢复，<span>一处管理</span></h2>
          <p>把重复、分散的备份工作收拢起来，让团队专注于真正重要的数据。</p>
        </div>
        <div class="feature-grid">
          <article v-for="feature in features" :key="feature.title">
            <div class="feature-icon"><fs-icon :icon="feature.icon" /></div>
            <h3>{{ feature.title }}</h3>
            <p>{{ feature.description }}</p>
          </article>
        </div>
      </section>
    </main>
    <footer>
      <span>OctoBackup</span><span>远程备份管理 · {{ envRef.COPYRIGHT }}</span>
    </footer>
  </div>
</template>
<script lang="ts" setup>
import { ref } from "vue";
import { env } from "/@/utils/util.env";
const envRef = ref(env);
const features = [
  { title: "统一管理仓库", description: "Local、SFTP、S3 和 MinIO 仓库集中配置，密码只写入不回显。", icon: "lucide:database" },
  { title: "可追踪的备份任务", description: "按主机和计划组织任务，快速掌握每一次备份的执行状态。", icon: "lucide:calendar-check" },
  { title: "可靠的数据恢复", description: "从快照中找回需要的数据，恢复路径清晰，操作简单。", icon: "lucide:rotate-ccw" },
  { title: "轻量自托管", description: "Go、SQLite 和 Docker Compose，部署在自己的环境中。", icon: "lucide:server" },
  { title: "权限与审计", description: "按角色控制访问范围，让备份管理更符合团队协作。", icon: "lucide:shield-check" },
  { title: "基于 Restic", description: "沿用成熟的去重、加密和快照能力，避免重复造轮子。", icon: "lucide:box" }
];
</script>
<style lang="less" scoped>
.landing-page {
  min-height: 100vh;
  color: #142238;
  background: #f8fafc;
}
.landing-nav {
  position: sticky;
  top: 0;
  z-index: 10;
  border-bottom: 1px solid #e7edf4;
  background: rgba(248, 250, 252, 0.9);
  backdrop-filter: blur(14px);
}
.nav-inner {
  display: flex;
  height: 72px;
  max-width: 1180px;
  align-items: center;
  justify-content: space-between;
  margin: auto;
  padding: 0 24px;
}
.brand {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 20px;
  font-weight: 700;
}
.brand img {
  width: 34px;
  height: 34px;
}
.nav-actions {
  display: flex;
  align-items: center;
  gap: 24px;
}
.nav-actions > a:not(.btn) {
  color: #64748b;
  font-size: 14px;
  text-decoration: none;
}
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 7px;
  padding: 10px 20px;
  font-weight: 600;
  text-decoration: none;
  transition: 0.2s;
}
.btn-primary {
  color: #fff;
  background: #3178d4;
  box-shadow: 0 6px 16px #3178d433;
}
.btn-primary:hover {
  background: #2568be;
  transform: translateY(-1px);
}
.btn-outline {
  border: 1px solid #d6e0eb;
  color: #334155;
  background: #fff;
}
.btn-large {
  padding: 14px 25px;
}
.btn-quiet {
  color: #516174;
}
.hero {
  display: grid;
  grid-template-columns: 1.05fr 0.95fr;
  gap: 60px;
  max-width: 1180px;
  min-height: 600px;
  align-items: center;
  margin: auto;
  padding: 80px 24px;
}
.eyebrow {
  color: #3178d4;
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.16em;
}
.hero h1 {
  margin: 18px 0;
  font-size: clamp(38px, 5vw, 64px);
  line-height: 1.12;
  font-weight: 400;
}
.hero h1 strong {
  color: #3178d4;
  font-weight: 700;
}
.hero-content > p {
  max-width: 570px;
  color: #64748b;
  font-size: 17px;
  line-height: 1.8;
}
.hero-actions {
  display: flex;
  gap: 12px;
  margin-top: 30px;
}
.trust-row {
  display: flex;
  flex-wrap: wrap;
  gap: 22px;
  margin-top: 34px;
  color: #748397;
  font-size: 13px;
}
.trust-row b,
.green {
  color: #20a574;
}
.hero-visual {
  position: relative;
  min-height: 380px;
}
.visual-grid {
  position: absolute;
  inset: 24px 0 24px 24px;
  background-image: linear-gradient(#dce8f4 1px, transparent 1px), linear-gradient(90deg, #dce8f4 1px, transparent 1px);
  background-size: 32px 32px;
  opacity: 0.65;
}
.backup-panel {
  position: absolute;
  top: 78px;
  left: 42px;
  width: min(430px, calc(100% - 42px));
  padding: 28px;
  border: 1px solid #dfe8f2;
  border-radius: 12px;
  background: #fff;
  box-shadow: 0 24px 60px #31547b1a;
}
.panel-top,
.panel-footer {
  display: flex;
  justify-content: space-between;
  color: #8291a3;
  font-size: 10px;
  letter-spacing: 0.08em;
}
.panel-dot {
  width: 8px;
  height: 8px;
  margin-right: auto;
  border-radius: 50%;
  background: #25b47e;
}
.panel-live {
  color: #20a574;
}
.panel-number {
  margin-top: 36px;
  color: #142238;
  font-size: 65px;
  font-weight: 700;
  line-height: 1;
}
.panel-number small {
  color: #3178d4;
  font-size: 28px;
}
.panel-caption {
  margin-top: 10px;
  color: #8291a3;
  font-size: 13px;
}
.panel-bars {
  display: flex;
  align-items: end;
  gap: 7px;
  height: 85px;
  margin: 24px 0;
  border-bottom: 1px solid #edf1f5;
}
.panel-bars i {
  flex: 1;
  height: 45%;
  border-radius: 3px 3px 0 0;
  background: #9fc6ed;
}
.panel-bars i:nth-child(3n) {
  height: 80%;
  background: #3178d4;
}
.panel-bars i:nth-child(4n) {
  height: 62%;
}
.panel-footer {
  letter-spacing: 0;
  font-size: 12px;
}
.features {
  padding: 95px 24px 110px;
  background: #fff;
}
.section-heading {
  max-width: 1180px;
  margin: auto;
  text-align: center;
}
.section-heading h2 {
  margin: 14px 0 10px;
  font-size: 38px;
  font-weight: 400;
}
.section-heading h2 span {
  color: #3178d4;
  font-weight: 700;
}
.section-heading p {
  color: #748397;
}
.feature-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 18px;
  max-width: 1180px;
  margin: 48px auto 0;
}
.feature-grid article {
  padding: 26px;
  border: 1px solid #e6edf4;
  border-radius: 9px;
}
.feature-icon {
  display: grid;
  place-items: center;
  width: 42px;
  height: 42px;
  margin-bottom: 20px;
  border-radius: 8px;
  color: #3178d4;
  background: #edf5ff;
  font-size: 21px;
}
.feature-grid h3 {
  margin: 0 0 8px;
  font-size: 17px;
}
.feature-grid p {
  margin: 0;
  color: #748397;
  font-size: 14px;
  line-height: 1.7;
}
footer {
  display: flex;
  justify-content: space-between;
  max-width: 1180px;
  margin: auto;
  padding: 26px 24px;
  color: #8795a5;
  font-size: 12px;
}
footer span:first-child {
  color: #31445d;
  font-weight: 700;
}
@media (max-width: 800px) {
  .nav-actions > a:not(.btn) {
    display: none;
  }
  .hero {
    grid-template-columns: 1fr;
    padding-top: 60px;
  }
  .hero-visual {
    min-height: 330px;
  }
  .feature-grid {
    grid-template-columns: 1fr;
  }
  .hero-actions {
    flex-wrap: wrap;
  }
  footer {
    flex-direction: column;
    gap: 8px;
  }
}
</style>
