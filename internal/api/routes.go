package api

import (
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.Default()
	registerRoutes(router)
	return router
}

func registerRoutes(router *gin.Engine) {
	router.GET("/health", Health)
	router.GET("/summary", Summary)
}
