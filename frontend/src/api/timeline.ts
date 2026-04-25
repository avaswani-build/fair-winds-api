import type { TimelinePoint } from "../types/timeline";

export async function fetchTimeline(
  lat: number,
  lng: number
): Promise<TimelinePoint[]> {
  const res = await fetch(`/timeline?lat=${lat}&lng=${lng}`);

  if (!res.ok) {
    throw new Error("Failed to fetch wind timeline");
  }

  return res.json();
}