package api

import (
	"github.com/avaswani-build/fair-winds-api/internal/weather"
	"github.com/gin-gonic/gin"
)

func NewRouter(client weather.Client) *gin.Engine {
	router := gin.Default()
	handler := &Handler{
		WeatherClient: client,
	}
	registerRoutes(router, handler)
	return router
}

func registerRoutes(router *gin.Engine, h *Handler) {
	router.GET("/health", Health)
	router.GET("/summary-mock", SummaryMock)
	router.GET("/summary", h.Summary)
	router.GET("/timeline", h.GetTimeline)
}
