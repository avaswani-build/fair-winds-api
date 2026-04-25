export type WindLevel = "light" | "medium" | "heavy";

export type TimelinePoint = {
  time: string;
  level: WindLevel;
};