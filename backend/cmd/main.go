package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	
	"dater/backend/internal/logger"
)

func main() {
	logger.Init()
	
	r := gin.New()
	r.Use(cors.Default())
	r.Use(logger.GinLogger())
	
	r.GET("/ping", func(c *gin.Context) {
		logger.Log.Info().Msg("Ping returned a response, the server is running and connection is successful")
		c.JSON(200, gin.H{"message": "connection is successful"})
	})
	
	logger.Log.Info().Msg("Server started successfully: http://localhost:8080")
	r.Run(":8080")
}