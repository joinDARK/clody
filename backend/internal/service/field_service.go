package service

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"strconv"
)

type CreateFieldInput struct {
	TableID     int64            `json:"table_id" binding:"required"`
	Name        string           `json:"name" binding:"required,min=3,max=100"`
	Type        FieldType `json:"type" binding:"required"`
	Position    int              `json:"position" binding:"required"`
	Options     datatypes.JSON   `json:"options" binding:"omitempty"`
	Description *string          `json:"description" binding:"omitempty"`
}

type UpdateFieldInput struct {
	TableID     int64            `json:"table_id" binding:"omitempty"`
	Name        string           `json:"name" binding:"omitempty,min=3,max=100"`
	Type        FieldType `json:"type" binding:"omitempty"`
	Position    int              `json:"position" binding:"omitempty"`
	Options     datatypes.JSON   `json:"options" binding:"omitempty"`
	Description string           `json:"description" binding:"omitempty"`
}

func GetFields(db *gorm.DB, logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Debug().Msg("Getting Fields...")

		var fields []Field
		if err := db.Find(&fields).Error; err != nil {
			logger.Error().Err(err).Msg("Failed to get fields")
			c.JSON(500, gin.H{"error": "Failed to get fields"})
			return
		}

		logger.Debug().Msg("Fields retrieved")
		logger.Info().Msg("Fields retrieved")
		c.JSON(200, gin.H{"data": fields})
	}
}

func CreateField(db *gorm.DB, logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Debug().Msg("Creating Field...")

		var input CreateFieldInput
		if err := c.ShouldBindJSON(&input); err != nil {
			logger.Error().Err(err).Msg("Failed to bind JSON")
			c.JSON(400, gin.H{"error": "Failed to bind JSON"})
			return
		}

		field := Field{
			TableID:     input.TableID,
			Name:        input.Name,
			Type:        input.Type,
			Position:    input.Position,
			Options:     input.Options,
			Description: input.Description,
		}

		if err := db.Create(&field).Error; err != nil {
			logger.Error().Err(err).Msg("Failed to create field")
			c.JSON(500, gin.H{"error": "Failed to create field"})
			return
		}

		logger.Debug().Msg("Field created")
		logger.Info().Msg("Field created")
		c.JSON(201, gin.H{"data": field})
	}
}

func UpdateField(db *gorm.DB, logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Debug().Msg("Updating Field...")

		id := c.Param("id")
		if id == "" {
			logger.Error().Msg("Empty ID in URI")
			c.JSON(404, gin.H{"error": "Empty ID in URI"})
			return
		}

		var input UpdateFieldInput
		if err := c.ShouldBindJSON(&input); err != nil {
			logger.Error().Err(err).Msg("Failed to bind JSON")
			c.JSON(400, gin.H{"error": "Failed to bind JSON"})
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
		if input.Type != "" {
			updates["type"] = input.Type
		}
		if input.Position != 0 {
			updates["position"] = input.Position
		}
		if input.Options != nil {
			updates["options"] = input.Options
		}
		if input.TableID != 0 {
			updates["table_id"] = input.TableID
		}

		var updateField Field
		if err := db.First(&updateField, i).Updates(updates).Error; err != nil {
			logger.Error().Err(err).Msg("Failed to update field")
			c.JSON(500, gin.H{"error": "Failed to update field"})
			return
		}

		logger.Debug().Msg("Field updated")
		logger.Info().Msg("Field updated")
		c.JSON(200, gin.H{"data": updateField})
	}
}

func DeleteField(db *gorm.DB, logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Debug().Msg("Deleting Field...")

		id := c.Param("id")
		if id == "" {
			logger.Error().Msg("Empty ID in URI")
			c.JSON(404, gin.H{"error": "Empty ID in URI"})
			return
		}

		var field Field
		if err := db.First(&field, id).Error; err != nil {
			logger.Error().Err(err).Msg("Failed to find field")
			c.JSON(404, gin.H{"error": "Failed to find field"})
			return
		}

		if err := db.Delete(&field).Error; err != nil {
			logger.Error().Err(err).Msg("Failed to delete field")
			c.JSON(500, gin.H{"error": "Failed to delete field"})
			return
		}

		logger.Debug().Msg("Field deleted")
		logger.Info().Msg("Field deleted")
		c.JSON(204, gin.H{})
	}
}
