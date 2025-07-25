package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/rs/zerolog"
)

func InitTableRoutes(r *gin.Engine, db *gorm.DB, logger zerolog.Logger) {
	table := r.Group("/tables")
	table.POST("/", CreateTable(db, logger))
	table.GET("/", GetTables(db, logger))
	table.PATCH("/:id", UpdateTable(db, logger))
	table.DELETE("/:id", DeleteTable(db, logger))
}