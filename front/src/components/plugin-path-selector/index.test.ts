import { describe, expect, it } from "vitest";
import { formatPathInput, normalizePathNodes, parsePathInput } from "/src/components/plugin-path-selector/helpers";

describe("plugin path selector helpers", () => {
  it("parses pasted paths into a unique array", () => {
    expect(parsePathInput("/etc\n /var/www\n/etc\n")).toEqual(["/etc", "/var/www"]);
  });

  it("formats array values for editing", () => {
    expect(formatPathInput(["/etc", "/var/www"])).toBe("/etc\n/var/www");
  });

  it("normalizes action directory results into tree nodes", () => {
    expect(
      normalizePathNodes({
        paths: [
          { path: "/etc", isLeaf: true },
          { name: "/var", children: [{ path: "/var/log" }] }
        ]
      })
    ).toEqual([
      { title: "/etc", key: "/etc", isLeaf: true },
      { title: "/var", key: "/var", isLeaf: false, children: [{ title: "/var/log", key: "/var/log", isLeaf: true }] }
    ]);
  });
});
