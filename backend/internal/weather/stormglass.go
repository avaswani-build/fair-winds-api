package weather

type StormglassResponse struct {
	Hours []Hour `json:"hours"`
}

type Hour struct {
	Time          string      `json:"time"`
	WindSpeed     SourceValue `json:"windSpeed"`
	Gust          SourceValue `json:"gust"`
	WaveHeight    SourceValue `json:"waveHeight"`
	WindDirection SourceValue `json:"windDirection"`
}

type SourceValue struct {
	SG float64 `json:"sg"`
}
