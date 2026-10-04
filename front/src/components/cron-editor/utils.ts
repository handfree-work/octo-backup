import parser from "cron-parser";
import dayjs from "dayjs";

export function isFivePartCron(cron?: string): boolean {
  if (typeof cron !== "string" || cron.trim() === "") {
    return false;
  }

  return cron.trim().split(/\s+/).length === 5;
}

export function getCronNextTimes(cron: string, count: number = 1) {
  if (!isFivePartCron(cron)) {
    return [];
  }
  const nextTimes = [];
  const interval = parser.parseExpression(cron);
  for (let i = 0; i < count; i++) {
    const next = interval.next().getTime();
    nextTimes.push(dayjs(next).format("YYYY-MM-DD HH:mm:ss"));
  }
  return nextTimes;
}
