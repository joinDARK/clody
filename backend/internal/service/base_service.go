package service

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type CreateBaseInput struct {
	Name        string  `json:"name" binding:"required,min=3,max=100"`
	Description *string `json:"description" binding:"omitempty"`
}

type UpdateBaseInput struct {
	Name        string `json:"name" binding:"omitempty,min=3,max=100"`
	Description string `json:"description" binding:"omitempty"`
}

func GetBases(db *gorm.DB, logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Debug().Msg("Getting base...")

		var base []Base
		if err := db.Find(&base).Error; err != nil {
			logger.Error().Err(err).Msg("Failed to get base")
			c.JSON(500, gin.H{"error": "Failed to get base"})
			return
		}

		logger.Debug().Msg("Bases retrieved")
		logger.Info().Msg("Bases retrieved")
		c.JSON(200, gin.H{"data": base})
	}
}

func CreateBase(db *gorm.DB, logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Debug().Msg("Creating base...")

		var input CreateBaseInput
		if err := c.ShouldBindJSON(&input); err.Error() == "EOF" || err != nil {
			logger.Error().Err(err).Msg("Failed to bind base")
			c.JSON(400, gin.H{"error": "Failed to bind base"})
			return
		}

		if input.Name == "" {
			input.Name = "Default"
		}

		defaultBase := &Base{
			Name:        input.Name,
			Description: input.Description,
		}

		if err := db.Create(defaultBase).Error; err != nil {
			logger.Error().Err(err).Msg("Failed to create base")
			c.JSON(500, gin.H{"error": "Failed to create base"})
			return
		}

		logger.Debug().Msg("Base created")
		logger.Info().Msg("Base created")
		c.JSON(201, gin.H{"data": defaultBase})
	}
}

func UpdateBase(db *gorm.DB, logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем ID пользователя из пути
		id := c.Param("id")
		if id == "" {
			logger.Error().Msg("Empty ID in URI")
			c.JSON(404, gin.H{"error": "Empty ID in URI"})
			return
		}

		logger.Debug().Msg("Updating base...")

		var input UpdateBaseInput
		if err := c.ShouldBindJSON(&input); err != nil && err.Error() != "EOF" {
			logger.Error().Err(err).Msg("Failed to bind base")
			c.JSON(400, gin.H{"error": "Failed to bind base"})
			return
		}

		i, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to parse ID")
			c.JSON(400, gin.H{"error": "Failed to parse ID"})
			return
		}

		updates := map[string]interface{}{}
		if input.Name != "" {
			updates["name"] = input.Name
		}
		if input.Description != "" {
			updates["description"] = input.Description
		}

		if len(updates) == 0 {
			logger.Error().Msg("No fields to update")
			c.JSON(400, gin.H{"error": "No fields to update"})
			return
		}

		if err := db.Model(&Base{}).Where("id = ?", i).Updates(updates).Error; err != nil {
			logger.Error().Err(err).Msg("Failed to update base")
			c.JSON(500, gin.H{"error": "Failed to update base"})
			return
		}
		
		var updateBase Base
		if err := db.First(&updateBase, i).Error; err != nil {
			logger.Error().Err(err).Msg("Failed to get updated base")
			c.JSON(500, gin.H{"error": "Failed to get updated base"})
			return
		}

		logger.Debug().Msg("Base updated")
		logger.Info().Msg("Base updated")
		c.JSON(205, gin.H{"data": updateBase})
	}
}

func DeleteBase(db *gorm.DB, logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем ID базы из пути
		id := c.Param("id")
		if id == "" {
			logger.Error().Msg("Empty ID in URI")
			c.JSON(404, gin.H{"error": "Empty ID in URI"})
			return
		}

		logger.Debug().Msg("Deleting base...")

		i, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to parse ID")
			c.JSON(400, gin.H{"error": "Failed to parse ID"})
			return
		}

		if err := db.Delete(&Base{}, i).Error; err != nil {
			logger.Error().Err(err).Msg("Failed to delete base")
			c.JSON(500, gin.H{"error": "Failed to delete base"})
			return
		}

		logger.Debug().Msg("Base deleted")
		logger.Info().Msg("Base deleted")
		c.JSON(204, gin.H{})
	}
}
