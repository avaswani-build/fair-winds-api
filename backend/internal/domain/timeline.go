package domain

type WindLevel string

const (
	WindLight  WindLevel = "light"
	WindMedium WindLevel = "medium"
	WindHeavy  WindLevel = "heavy"
)

type TimelinePoint struct {
	Time      string    `json:"time"`
	WindSpeed float64   `json:"windSpeed"`
	Gust      float64   `json:"gustSpeed"`
	Level     WindLevel `json:"level"`
}
