<template>
  <div class="cron-editor">
    <div class="flex-o">
      <cron-light :disabled="disabled || readonly" :period="period" class="flex-o cron-ant" locale="zh-CN" format="crontab" :model-value="modelValue" @update:model-value="onUpdate" @error="onError" />
    </div>
    <div class="mt-2 flex">
      <a-input :disabled="true" :readonly="readonly" :value="modelValue" @change="onChange"></a-input>
      <fs-icon icon="ion:close-circle" class="pointer fs-16 ml-5 color-gray" title="清除" @click="onClear"></fs-icon>
    </div>
    <div class="helper">下次触发时间：{{ nextTime }}</div>
    <div class="fs-helper">{{ errorMessage }}</div>
  </div>
</template>

<script lang="ts" setup>
import { computed, ref } from "vue";
import { useI18n } from "vue-i18n";
import { CronLight } from "@vue-js-cron/light";
import { getCronNextTimes } from "./utils";

const { t } = useI18n();
defineOptions({
  name: "CronEditor"
});
const props = defineProps<{
  modelValue?: string;
  disabled?: boolean;
  readonly?: boolean;
}>();

const period = ref<string>("");
if (props.modelValue == null || props.modelValue.endsWith("* * *")) {
  period.value = "day";
} else if (props.modelValue.endsWith("* *")) {
  period.value = "month";
} else if (props.modelValue.endsWith("*")) {
  period.value = "year";
}
const emit = defineEmits<{
  "update:modelValue": any;
  change: any;
}>();

const errorMessage = ref<string | null>(null);

const onUpdate = (value: string) => {
  if (value === props.modelValue) {
    return;
  }
  const arr: string[] = value.split(" ");

  value = arr.join(" ");

  emit("update:modelValue", value);
  errorMessage.value = undefined;
};

const onPeriod = (value: string) => {
  period.value = value;
};

const onChange = (e: any) => {
  const value = e.target.value;
  onUpdate(value);
};
const onError = (error: any) => {
  errorMessage.value = error;
};

const onClear = () => {
  if (props.disabled) {
    return;
  }
  onUpdate("");
};

const nextTime = computed(() => {
  if (props.modelValue == null) {
    return "请输入表达式";
  }

  try {
    const nextTimes = getCronNextTimes(props.modelValue, 2);
    return nextTimes.join("，");
  } catch (e) {
    console.log(e);
    return "表达式错误";
  }
});
</script>
<style lang="less">
.cron-editor {
  .cron-ant {
    flex-wrap: wrap;

    &* > {
      margin-bottom: 2px;
      display: flex;
      align-items: center;
    }

    .vcron-select-list {
      min-width: 56px;
    }

    .vcron-select-input {
      min-height: 22px;
      background-color: #fff;
    }

    .vcron-select-container {
      display: flex;
      align-items: center;
    }
  }
}
</style>
