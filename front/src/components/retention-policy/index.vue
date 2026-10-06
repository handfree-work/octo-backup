<template>
  <div class="retention-policy component-retention-policy">
    <a-radio-group :value="mode" @change="onModeChange">
      <a-radio-button value="count">按数量</a-radio-button>
      <a-radio-button value="time">按时间</a-radio-button>
      <a-radio-button value="none">不清理</a-radio-button>
    </a-radio-group>
    <div v-if="mode === 'count'" class="retention-policy__fields mt-2">
      <a-input-number v-for="field in countFields" :key="field.key" :value="policyValue[field.key] || 0" :min="0" :addon-before="field.label" @change="updateField(field.key, $event)" />
    </div>
    <div v-else-if="mode === 'time'" class="retention-policy__time mt-2">
      <div class="retention-policy__time-input">
        <a-input-number :value="timeNumber" :min="1" required placeholder="数值" @change="onNumberChange" />
        <a-select :value="timeUnit" @change="onUnitChange">
          <a-select-option value="h">小时</a-select-option>
          <a-select-option value="d">天</a-select-option>
          <a-select-option value="w">周</a-select-option>
          <a-select-option value="m">月</a-select-option>
          <a-select-option value="y">年</a-select-option>
        </a-select>
      </div>
      <div class="helper mt-2">{{ timeHelpText }}</div>
    </div>
    <div v-else class="helper mt-2">不执行快照清理</div>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
defineOptions({ name: "RetentionPolicy" });
type Policy = { type?: "count" | "time" | "none"; value?: Record<string, number | string> };
const props = withDefaults(defineProps<{ modelValue?: Policy }>(), { modelValue: () => ({}) });
const emit = defineEmits<{ (event: "update:modelValue", value: Policy): void }>();
const value = computed(() => props.modelValue || {});
const policyValue = computed(() => value.value.value || {});
const countFields = [
  { key: "keepLast", label: "最近" },
  { key: "keepHourly", label: "每小时" },
  { key: "keepDaily", label: "每日" },
  { key: "keepWeekly", label: "每周" },
  { key: "keepMonthly", label: "每月" },
  { key: "keepYearly", label: "每年" }
] as const;
const mode = computed(() => value.value.type || "none");
const timeNumber = computed(() => Number(policyValue.value.duration || 0) || undefined);
const timeUnit = computed(() => String(policyValue.value.unit || "d"));
const timeHelpText = computed(() => {
  if (!timeNumber.value) {
    return "请输入保留时长";
  }
  const unitLabels: Record<string, string> = { h: "小时", d: "天", w: "周", m: "月", y: "年" };
  return `保留 ${timeNumber.value} ${unitLabels[timeUnit.value] || "时间"}之内的数据`;
});
function onModeChange(event: { target: { value: "count" | "time" | "none" } }) {
  const next: Policy = { type: event.target.value, value: {} };
  if (event.target.value === "time") {
    next.keepWithin = "";
  } else if (event.target.value === "count") {
    for (const field of countFields) next.value![field.key] = 0;
  }
  emit("update:modelValue", next);
}
function onNumberChange(next: number | null) {
  emitTime(next || undefined, timeUnit.value);
}
function onUnitChange(next: string) {
  emitTime(timeNumber.value, next);
}
function emitTime(number: number | undefined, unit: string) {
  const keepWithin = number && number > 0 ? `${number}${unit}` : "";
  emit("update:modelValue", { type: "time", value: { duration: number || "", unit } });
}
function updateField(key: keyof Policy, nextValue: number | null) {
  emit("update:modelValue", { type: "count", value: { ...policyValue.value, [key]: nextValue || 0 } });
}
</script>

<style lang="less">
.retention-policy {
  .retention-policy__fields {
    display: grid;
    grid-template-columns: repeat(2, minmax(180px, 1fr));
    gap: 8px;
    margin-top: 8px;
  }
  .helper {
    display: block;
    margin-top: 8px;
    color: #687386;
  }
  .retention-policy__time-input {
    display: flex;
    gap: 8px;
    margin-top: 8px;
  }
}
</style>
