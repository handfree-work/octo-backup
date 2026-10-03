<template>
  <div v-if="showSourceLink" class="fs-source-link-group">
    <div class="fs-source-link" @click="goSource('https://github.com/handfree-work/octo-backup')">OctoBackup 源码</div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
export default defineComponent({
  name: "FsSourceLink",
  setup() {
    const router = useRouter();
    const showSourceLink = ref(false);
    watch(
      () => {
        return router.currentRoute.value.fullPath;
      },
      (value) => {
        showSourceLink.value = false;
      },
      { immediate: true }
    );
    const middle = "/";
    function goSource(prefix: any) {
      const path = router.currentRoute.value.fullPath;
      window.open(prefix + middle);
    }
    return {
      goSource,
      showSourceLink
    };
  }
});
</script>

<style lang="less">
.fs-source-link-group {
  position: fixed;
  right: 3px;
  bottom: 20px;
  z-index: 1000;
  .fs-source-link {
    text-align: left;
    cursor: pointer;
    font-size: 12px;
    border-radius: 5px 0 0 5px;
    padding: 5px;
    background: #666;
    color: #fff;
    margin-bottom: 5px;
  }
}
</style>
