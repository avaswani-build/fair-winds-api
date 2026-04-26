import type { SummaryResponse } from "../types/summary";

export async function getSummary(
  lat: number,
  long: number,
): Promise<SummaryResponse> {
  const response = await fetch(
    `/api/summary?lat=${encodeURIComponent(lat)}&lng=${encodeURIComponent(long)}`,
  );

  if (!response.ok) {
    throw new Error(`Failed to fetch summary: ${response.status}`);
  }

  const raw = await response.json();
  return {
    location: raw.forecast.location,
    windAvgKts: raw.forecast.wind_avg_kts,
    gustKts: raw.forecast.gust_kts,
    windDir: raw.forecast.wind_dir,
    recommendation: raw.recommendation,
  };
}