package handlers

import (
	"dater/backend/internal/logger"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	logger.Log.Info().Msg("Ping returned a response, the server is running and connection is successful")
	c.JSON(200, gin.H{"message": "connection is successful"})
}
