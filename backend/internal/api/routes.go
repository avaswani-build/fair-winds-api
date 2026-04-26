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
	api := router.Group("/api")
	api.GET("/health", Health)
	api.GET("/summary-mock", SummaryMock)
	api.GET("/summary", h.Summary)
	api.GET("/timeline", h.GetTimeline)
}
