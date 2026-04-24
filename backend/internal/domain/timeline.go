package domain

type WindLevel string

const (
	WindLight  WindLevel = "light"
	WindMedium WindLevel = "medium"
	WindHeavy  WindLevel = "heavy"
)

type TimelinePoint struct {
	Time  string    `json:"time"`
	Level WindLevel `json:"level"`
}
