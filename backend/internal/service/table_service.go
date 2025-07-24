package service

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type CreateTableInput struct {
	BaseID      int64   `json:"base_id" binding:"required"`
	Name        string  `json:"name" binding:"required,min=3,max=100"`
	Description *string `json:"description" binding:"omitempty"`
}

type UpdateTableInput struct {
	BaseID      int64  `json:"base_id" binding:"omitempty"`
	Name        string `json:"name" binding:"omitempty,min=3,max=100"`
	Description string `json:"description" binding:"omitempty"`
}

func GetTables(db *gorm.DB, logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Debug().Msg("Getting tables...")

		var table []Table
		if err := db.Find(&table).Error; err != nil {
			logger.Error().Err(err).Msg("Failed to get tables")
			c.JSON(500, gin.H{"error": "Failed to get tables"})
			return
		}

		logger.Debug().Msg("Tables retrieved")
		logger.Info().Msg("Tables retrieved")
		c.JSON(200, gin.H{"data": table})
	}
}

func CreateTable(db *gorm.DB, logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info().Msg("Creating table...")

		var input CreateTableInput
		if err := c.ShouldBindJSON(&input); err.Error() == "EOF" || err != nil {
			logger.Error().Err(err).Msg("Failed to bind table")
			c.JSON(400, gin.H{"error": "Failed to bind table"})
			return
		}

		if input.Name == "" {
			input.Name = "New Table"
		}

		defaultTable := &Table{
			BaseID:      input.BaseID,
			Name:        input.Name,
			Description: input.Description,
		}

		if err := db.Create(defaultTable).Error; err != nil {
			logger.Error().Err(err).Msg("Failed to create table")
			c.JSON(500, gin.H{"error": "Failed to create table"})
			return
		}

		logger.Info().Msg("Table created")
		c.JSON(201, gin.H{"data": defaultTable})
	}
}

func UpdateTable(db *gorm.DB, logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем ID таблицы из пути
		id := c.Param("id")
		if id == "" {
			logger.Error().Msg("Empty ID in URI")
			c.JSON(404, gin.H{"error": "Empty ID in URI"})
			return
		}

		logger.Info().Msg("Updating table...")

		var input UpdateTableInput
		if err := c.ShouldBindJSON(&input); err != nil && err.Error() != "EOF" {
			logger.Error().Err(err).Msg("Failed to bind table")
			c.JSON(400, gin.H{"error": "Failed to bind table"})
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

		if err := db.Model(&Table{}).Where("id = ?", i).Updates(updates).Error; err != nil {
			logger.Error().Err(err).Msg("Failed to update table")
			c.JSON(500, gin.H{"error": "Failed to update table"})
			return
		}

		var updateTable Table
		if err := db.First(&updateTable, i).Error; err != nil {
			logger.Error().Err(err).Msg("Failed to get updated table")
			c.JSON(500, gin.H{"error": "Failed to get updated table"})
			return
		}

		logger.Debug().Msg("Table updated")
		logger.Info().Msg("Table updated")
		c.JSON(205, gin.H{"data": updateTable})
	}
}

func DeleteTable(db *gorm.DB, logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем ID базы из пути
		id := c.Param("id")
		if id == "" {
			logger.Error().Msg("Empty ID in URI")
			c.JSON(404, gin.H{"error": "Empty ID in URI"})
			return
		}

		logger.Debug().Msg("Deleting table...")

		i, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to parse ID")
			c.JSON(400, gin.H{"error": "Failed to parse ID"})
			return
		}

		if err := db.Delete(&Table{}, i).Error; err != nil {
			logger.Error().Err(err).Msg("Failed to delete table")
			c.JSON(500, gin.H{"error": "Failed to delete table"})
			return
		}

		logger.Debug().Msg("Table deleted")
		logger.Info().Msg("Table deleted")
		c.JSON(204, gin.H{})
	}
}
