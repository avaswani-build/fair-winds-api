export type SummaryResponse = {
  location: string
  windAvgKts: number
  gustKts: number
  windDir: string
  recommendation: {
    status: string;
    reason: string;
  }
}