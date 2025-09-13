package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(cfg *Config) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	// TODO: attach middleware

	api := r.Group("/api")
	{
		api.POST("/shorten", shortenHandler)
		api.GET("/stats/:code", statsHandler)
	}
	r.GET("/r/:code", redirectHandler)
	return r
}

func shortenHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

func statsHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

func redirectHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
