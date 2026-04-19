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

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusPaymentRequired:
		return domain.Forecast{}, ErrPaymentRequired
	case http.StatusForbidden:
		return domain.Forecast{}, ErrForbidden
	case http.StatusUnprocessableEntity:
		return domain.Forecast{}, ErrUnprocessableContent
	case http.StatusServiceUnavailable:
		return domain.Forecast{}, ErrServiceUnavailable
	default:
		return domain.Forecast{}, ErrUpstream
	}

	var sgResp StormglassResponse
	if err := json.NewDecoder(resp.Body).Decode(&sgResp); err != nil {
		return domain.Forecast{}, err
	}

	if len(sgResp.Hours) == 0 {
		return domain.Forecast{}, fmt.Errorf("no forecast data returned")
	}

	hour := sgResp.Hours[0]

	forecast := domain.Forecast{
		Location:   fmt.Sprintf("%.4f, %.4f", lat, lng),
		WindAvgKts: hour.WindSpeed.SG,
		GustKts:    hour.Gust.SG,
		WindDir:    degreesToCardinal(hour.WindDirection.SG),
	}

	return forecast, nil
}
