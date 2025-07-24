package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/rs/zerolog"
)

func InitFieldRoutes(r *gin.Engine, db *gorm.DB, logger zerolog.Logger) {
	field := r.Group("/fields")
	field.POST("/", CreateField(db, logger))
	field.GET("/", GetFields(db, logger))
	field.PATCH("/:id", UpdateField(db, logger))
	field.DELETE("/:id", DeleteField(db, logger))
}