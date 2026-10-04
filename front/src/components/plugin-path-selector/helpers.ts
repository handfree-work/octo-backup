export interface PluginPathNode {
  title: string;
  key: string;
  children?: PluginPathNode[];
  isLeaf?: boolean;
}

export function parsePathInput(value: string | string[] | null | undefined): string[] {
  const values = Array.isArray(value) ? value : String(value || "").split(/\r?\n/);
  return [...new Set(values.map((item) => String(item).trim()).filter(Boolean))];
}

export function formatPathInput(value: string | string[] | null | undefined): string {
  return parsePathInput(value).join("\n");
}

export function normalizePathNodes(value: unknown): PluginPathNode[] {
  const source = Array.isArray(value) ? value : (value as any)?.paths || (value as any)?.records || (value as any)?.data || [];
  if (!Array.isArray(source)) return [];
  return source
    .map((item: any) => {
      const path = String(item?.path ?? item?.key ?? item?.name ?? item ?? "");
      const children = normalizePathNodes(item?.children);
      return { title: String(item?.title ?? item?.name ?? path), key: path, isLeaf: item?.isLeaf ?? children.length === 0, ...(children.length ? { children } : {}) };
    })
    .filter((item) => item.key);
}
