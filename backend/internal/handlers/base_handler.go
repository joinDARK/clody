package base

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func InitBaseRoutes(r *gin.Engine, db *gorm.DB, logger zerolog.Logger) {
	base := r.Group("/bases")
	base.POST("/", CreateBase(db, logger))
	base.GET("/", GetBases(db, logger))
	base.PATCH("/:id", UpdateBase(db, logger))
	base.DELETE("/:id", DeleteBase(db, logger))
}