import { request } from "/src/api/service";

export interface PluginField {
  key: string;
  title: string;
  type: string;
  required?: boolean;
  encrypt?: boolean;
  default?: unknown;
  placeholder?: string;
  component?: Record<string, unknown>;
  options?: Array<{ label: string; value: string }>;
}
export interface PluginAction {
  name: string;
  title?: string;
  permission?: string;
}
export interface PluginMetadata {
  type: string;
  name: string;
  title: string;
  description: string;
  version: string;
  icon?: string;
  group?: string;
  fields?: PluginField[];
  actions?: PluginAction[];
}
export interface PluginInstance {
  id: string | number;
  name: string;
  pluginType: string;
  pluginName: string;
  description?: string;
  config?: Record<string, unknown> & { configured?: boolean };
}
export const getPluginMetadata = (data: { type?: string; name?: string } = {}): Promise<PluginMetadata[]> => request({ url: "/plugin/metadata", method: "post", data });
export const getPluginPage = (data: { offset?: number; limit?: number; pluginType?: string; name?: string } = {}) => request({ url: "/plugin/page", method: "post", data });
export const createPlugin = (data: { name: string; pluginType: string; pluginName: string; config: Record<string, unknown>; description?: string }) => request({ url: "/plugin/create", method: "post", data });
export const updatePlugin = (id: string | number, data: Partial<{ name: string; pluginType: string; pluginName: string; config: Record<string, unknown>; description: string }>) =>
  request({ url: `/plugin/update?id=${id}`, method: "post", data });
export const deletePlugin = (id: string | number) => request({ url: `/plugin/delete?id=${id}`, method: "post", data: {} });
export const executePluginAction = (id: string | number, action: string, params: Record<string, unknown> = {}) => request({ url: "/plugin/action", method: "post", data: { id, action, params } });
