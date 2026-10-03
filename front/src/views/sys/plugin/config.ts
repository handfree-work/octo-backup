import type { PluginMetadata } from "./api";

export function collectPluginConfig(form: Record<string, any>, metadata: PluginMetadata[]): Record<string, unknown> {
  const config: Record<string, unknown> = {};
  const fields = metadata.find((item) => item.name === form.pluginName)?.fields || [];

  fields.forEach((field) => {
    const value = Object.prototype.hasOwnProperty.call(form, field.key) ? form[field.key] : form.config?.[field.key];
    if (value === undefined || (field.encrypt && (value === "" || isUnchangedMask(value, form, field.key)))) return;
    config[field.key] = value;
  });

  return config;
}

function isUnchangedMask(value: unknown, form: Record<string, any>, key: string): boolean {
  if (typeof value !== "string" || !value.includes("*")) return false;
  return form.config?.[`${key}Configured`] === true && value === form.config?.[`${key}Original`];
}
