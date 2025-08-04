package handlers

import (
	"dater/backend/internal/service"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type CreateBaseInput struct {
	Name        string  `json:"name" binding:"required,min=3,max=100"`
	Description *string `json:"description" binding:"omitempty"`
}

type UpdateBaseInput struct {
	Name        string  `json:"name" binding:"omitempty,min=3,max=100"`
	Description *string `json:"description" binding:"omitempty"`
}

func GetInfoBaseHandler(baseService *service.BaseService, logger zerolog.Logger) gin.HandlerFunc {
	service := *baseService

	return func(c *gin.Context) {
		logger.Debug().Msg("Getting base...")

		id := c.Param("id")
		i, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to parse ID")
			c.JSON(400, gin.H{"error": "Failed to parse ID"})
			return
		}

		base, err := service.GetBaseInfo(i, logger)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to get base")
			c.JSON(500, gin.H{"error": "Failed to get base"})
			return
		}

		_, withTables := c.GetQuery("WithTables")
		if withTables {
			base.Tables, err = service.GetBaseTables(i, logger)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to get base tables")
				c.JSON(500, gin.H{"error": "Failed to get base tables"})
				return
			}
		}

		logger.Info().Msg("Base retrieved")
		c.JSON(200, gin.H{
			"data": base,
			"_links": map[string]map[string]string{
				"self": {
					"href":   fmt.Sprintf("/api/bases/%d", base.ID),
					"params": "withTables",
				},
				"create": {
					"href":   "/api/bases",
					"method": "POST",
				},
				"update": {
					"href":   fmt.Sprintf("/api/bases/%d", base.ID),
					"method": "PUT",
				},
				"delete": {
					"href":   fmt.Sprintf("/api/bases/%d", base.ID),
					"method": "DELETE",
				},
			},
		})
	}
}

func CreateBaseHandler(baseService *service.BaseService, logger zerolog.Logger) gin.HandlerFunc {
	service := *baseService

	return func(c *gin.Context) {
		logger.Debug().Msg("Creating base...")

		var input CreateBaseInput
		if err := c.ShouldBindJSON(&input); err != nil {
			logger.Error().Err(err).Msg("Failed to bind JSON")
			c.JSON(400, gin.H{"error": "Failed to bind JSON"})
			return
		}

		newBase, err := service.CreateBase(input.Name, input.Description, logger)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to create base")
			c.JSON(500, gin.H{"error": "Failed to create base"})
			return
		}

		logger.Info().Msg("Base created")
		c.JSON(201, gin.H{"data": newBase})
	}
}

func UpdateBaseHandler(baseService *service.BaseService, logger zerolog.Logger) gin.HandlerFunc {
	service := *baseService

	return func(c *gin.Context) {
		logger.Debug().Msg("Updating base...")

		id := c.Param("id")
		i, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to parse ID")
			c.JSON(400, gin.H{"error": "Failed to parse ID"})
			return
		}

		var input UpdateBaseInput
		if err := c.ShouldBindJSON(&input); err != nil {
			logger.Error().Err(err).Msg("Failed to bind JSON")
			c.JSON(400, gin.H{"error": "Failed to bind JSON"})
			return
		}

		updBase, err := service.UpdateBase(i, input.Name, input.Description, logger)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to update base")
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		logger.Info().Msg("Base updated")
		c.JSON(200, gin.H{"data": updBase})
	}
}

func DeleteBaseHandler(baseService *service.BaseService, logger zerolog.Logger) gin.HandlerFunc {
	service := *baseService

	return func(c *gin.Context) {
		logger.Debug().Msg("Deleting base...")

		id := c.Param("id")
		i, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to parse ID")
			c.JSON(400, gin.H{"error": "Failed to parse ID"})
			return
		}

		deleteID, err := service.DeleteBase(i, logger)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to delete base")
			c.JSON(500, gin.H{"error": "Failed to delete base"})
			return
		}

		logger.Info().Msg("Base deleted")
		c.JSON(200, gin.H{"data": deleteID})
	}
}
