import { describe, expect, it, vi } from "vitest";
import { getPluginInstancePage } from "/src/views/sys/plugin/plugin-api";
import createPlanCrudOptions from "./crud";

vi.mock("/src/views/sys/plugin/plugin-api", () => ({
  getPluginInstancePage: vi.fn()
}));

vi.mock("./api", () => ({
  createBackupPlan: vi.fn(),
  deleteBackupPlan: vi.fn(),
  getBackupPlanInfo: vi.fn(),
  getBackupPlanPage: vi.fn(),
  updateBackupPlan: vi.fn()
}));

describe("backup plan cron field", () => {
  it("uses the cron editor and validates five-part expressions", () => {
    vi.mocked(getPluginInstancePage).mockResolvedValue({ records: [] } as any);
    const { crudOptions } = createPlanCrudOptions();
    const schedule = (crudOptions.columns as any).schedule;
    const rules = schedule.form.rules;

    expect(schedule.form.component).toMatchObject({ name: "cron-editor", allowEveryMin: true });
    expect(rules[1].pattern.test("0 2 * * *")).toBe(true);
    expect(rules[1].pattern.test("0 0 2 * * *")).toBe(false);
  });
});
