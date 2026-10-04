import { describe, expect, it } from "vitest";
import { isFivePartCron } from "./utils";

describe("cron editor utils", () => {
  it("accepts a standard five-part cron expression", () => {
    expect(isFivePartCron("0 2 * * *")).toBe(true);
    expect(isFivePartCron("*/15 9-18 * * 1-5")).toBe(true);
  });

  it("rejects expressions with seconds or missing fields", () => {
    expect(isFivePartCron("0 0 2 * * *")).toBe(false);
    expect(isFivePartCron("0 2 * *")).toBe(false);
    expect(isFivePartCron("")).toBe(false);
    expect(isFivePartCron(undefined)).toBe(false);
  });
});
