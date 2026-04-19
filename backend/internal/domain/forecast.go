package domain

type Forecast struct {
	Location   string  `json:"location"`
	WindAvgKts float64 `json:"wind_avg_kts"`
	GustKts    float64 `json:"gust_kts"`
	WindDir    string  `json:"wind_dir"`
}
