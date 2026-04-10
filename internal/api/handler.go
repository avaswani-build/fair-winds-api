package api

import (
	"net/http"

	"github.com/avaswani-build/fair-winds-api/internal/domain"
	"github.com/avaswani-build/fair-winds-api/internal/service"
	"github.com/gin-gonic/gin"
)

type SummaryResponse struct {
	Forecast       domain.Forecast       `json:"forecast"`
	Recommendation domain.Recommendation `json:"recommendation"`
}

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func Summary(c *gin.Context) {
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
