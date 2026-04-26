package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/avaswani-build/fair-winds-api/internal/domain"
)

type Client interface {
	GetForecast(lat, lng float64) (domain.Forecast, error)
	GetTimeline(lat, lng float64) ([]domain.TimelinePoint, error)
}

type StormglassClient struct {
	apiKey     string
	httpClient *http.Client
}

func NewStormglassClient() *StormglassClient {
	return &StormglassClient{
		apiKey: os.Getenv("STORMGLASS_API_KEY"),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

var (
	ErrPaymentRequired      = errors.New("stormglass api rate limit exceeded")
	ErrForbidden            = errors.New("stormglass API key missing or incorrect")
	ErrUnprocessableContent = errors.New("incorrect params for requested endpoint")
	ErrServiceUnavailable   = errors.New("stormglass internal error")
	ErrUpstream             = errors.New("stormglass upstream error")
)

type MockClient struct{}

func (m MockClient) GetForecast(lat, lng float64) (domain.Forecast, error) {
	return domain.Forecast{
		Location:   "Mock Location",
		WindAvgKts: 12,
		GustKts:    18,
		WindDir:    "SW",
	}, nil
}

func (m MockClient) GetTimeline(lat, lng float64) ([]domain.TimelinePoint, error) {
	return []domain.TimelinePoint{
		{Time: "2026-04-23T14:00:00Z", Level: domain.WindLight},
		{Time: "2026-04-23T15:00:00Z", Level: domain.WindMedium},
		{Time: "2026-04-23T16:00:00Z", Level: domain.WindHeavy},
	}, nil
}

func (c *StormglassClient) GetForecast(lat, lng float64) (domain.Forecast, error) {
	if c.apiKey == "" {
		return domain.Forecast{}, fmt.Errorf("missing STORMGLASS_API_KEY")
	}

	url := fmt.Sprintf(
		"https://api.stormglass.io/v2/weather/point?lat=%f&lng=%f&params=windSpeed,gust,waveHeight,windDirection",
		lat,
		lng,
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return domain.Forecast{}, err
	}

	req.Header.Set("Authorization", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return domain.Forecast{}, err
	}
	defer resp.Body.Close()

	if err := handleStormglassStatus(resp.StatusCode); err != nil {
		return domain.Forecast{}, err
	}

	var sgResp StormglassResponse
	if err := json.NewDecoder(resp.Body).Decode(&sgResp); err != nil {
		return domain.Forecast{}, err
	}

	if len(sgResp.Hours) == 0 {
		return domain.Forecast{}, fmt.Errorf("no forecast data returned")
	}

	hour := sgResp.Hours[0]

	return domain.Forecast{
		Location:   fmt.Sprintf("%.4f, %.4f", lat, lng),
		WindAvgKts: hour.WindSpeed.SG,
		GustKts:    hour.Gust.SG,
		WindDir:    degreesToCardinal(hour.WindDirection.SG),
	}, nil
}

func (c *StormglassClient) GetTimeline(lat, lng float64) ([]domain.TimelinePoint, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("missing STORMGLASS_API_KEY")
	}

	url := fmt.Sprintf(
		"https://api.stormglass.io/v2/weather/point?lat=%f&lng=%f&params=windSpeed",
		lat,
		lng,
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := handleStormglassStatus(resp.StatusCode); err != nil {
		return nil, err
	}

	var sgResp StormglassResponse
	if err := json.NewDecoder(resp.Body).Decode(&sgResp); err != nil {
		return nil, err
	}

	points := make([]domain.TimelinePoint, 0, len(sgResp.Hours))

	for _, hour := range sgResp.Hours {
		points = append(points, domain.TimelinePoint{
			Time:      hour.Time,
			WindSpeed: hour.WindSpeed.SG,
			Gust:      hour.Gust.SG,
			Level:     classifyWind(hour.WindSpeed.SG),
		})
	}

	return points, nil
}

func classifyWind(kts float64) domain.WindLevel {
	switch {
	case kts < 6:
		return domain.WindLight
	case kts < 15:
		return domain.WindMedium
	default:
		return domain.WindHeavy
	}
}

func handleStormglassStatus(statusCode int) error {
	switch statusCode {
	case http.StatusOK:
		return nil
	case http.StatusPaymentRequired:
		return ErrPaymentRequired
	case http.StatusForbidden:
		return ErrForbidden
	case http.StatusUnprocessableEntity:
		return ErrUnprocessableContent
	case http.StatusServiceUnavailable:
		return ErrServiceUnavailable
	default:
		return ErrUpstream
	}
}
