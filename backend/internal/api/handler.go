package api

import (
	"net/http"
	"strconv"

	"github.com/avaswani-build/fair-winds-api/internal/domain"
	"github.com/avaswani-build/fair-winds-api/internal/service"
	"github.com/avaswani-build/fair-winds-api/internal/weather"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	WeatherClient weather.Client
}

type SummaryResponse struct {
	Forecast       domain.Forecast       `json:"forecast"`
	Recommendation domain.Recommendation `json:"recommendation"`
}

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func SummaryMock(c *gin.Context) {
	forecast := domain.Forecast{
		Location:   "NY Harbor",
		WindAvgKts: 12,
		GustKts:    18,
		WindDir:    "SW",
	}

	recommendation := service.Recommend(forecast)

	c.JSON(http.StatusOK, SummaryResponse{
		Forecast:       forecast,
		Recommendation: recommendation,
	})
}

func (h *Handler) Summary(c *gin.Context) {
	latStr := c.Query("lat")
	lngStr := c.Query("lng")

	if latStr == "" || lngStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing required query params: lat and lng",
		})
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid lat value",
		})
		return
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid lng value",
		})
		return
	}

	if lng > 180 || lng < -180 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "longitude out of range",
		})
		return
	}

	if lat > 90 || lat < -90 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "latitude out of range",
		})
		return
	}

	forecast, err := h.WeatherClient.GetForecast(lat, lng)
	if err != nil {
		switch err {
		case weather.ErrPaymentRequired:
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": "Weather service usage limit reached. Please try again later.",
			})
		case weather.ErrForbidden:
			c.JSON(http.StatusBadGateway, gin.H{
				"error": "Weather service authentication failed.",
			})
		case weather.ErrUnprocessableContent:
			c.JSON(http.StatusBadGateway, gin.H{
				"error": "Weather service could not process the requested location.",
			})
		case weather.ErrServiceUnavailable:
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": "Weather service is temporarily unavailable. Please try again later.",
			})
		default:
			c.JSON(http.StatusBadGateway, gin.H{
				"error": "Unable to retrieve weather data right now.",
			})
		}
		return
	}

	recommendation := service.Recommend(forecast)

	c.JSON(http.StatusOK, SummaryResponse{
		Forecast:       forecast,
		Recommendation: recommendation,
	})
}
