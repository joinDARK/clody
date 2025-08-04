package handlers

import (
	"dater/backend/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type CreateTableInput struct {
	BaseID      int64   `json:"base_id" binding:"required"`
	Name        string  `json:"name" binding:"required,min=3,max=100"`
	Description *string `json:"description" binding:"omitempty"`
}

type UpdateTableInput struct {
	BaseID      int64   `json:"base_id" binding:"omitempty"`
	Name        string  `json:"name" binding:"omitempty,min=3,max=100"`
	Description *string `json:"description" binding:"omitempty"`
}

func GetInfoTableHandler(tableService *service.TableService, logger zerolog.Logger) gin.HandlerFunc {
	service := *tableService

	return func(c *gin.Context) {
		logger.Debug().Msg("Getting table info...")

		id := c.Param("id")
		i, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to parse ID")
			c.JSON(400, gin.H{"error": "Failed to parse ID"})
			return
		}

		tableInfo, err := service.GetTableByID(i, logger)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to get table info")
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		logger.Info().Msg("Table info retrieved")
		c.JSON(200, gin.H{"data": tableInfo})
	}
}

func CreateTableHandler(tableService *service.TableService, logger zerolog.Logger) gin.HandlerFunc {
	service := *tableService

	return func(c *gin.Context) {
		logger.Debug().Msg("Creating table...")

		var input CreateTableInput
		if err := c.ShouldBindJSON(&input); err != nil {
			logger.Error().Err(err).Msg("Failed to bind JSON")
			c.JSON(400, gin.H{"error": "Failed to bind JSON"})
			return
		}

		tableID, err := service.CreateTable(input.Name, input.Description, input.BaseID, logger)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to create table")
			c.JSON(500, gin.H{"error": "Failed to create table"})
			return
		}

		logger.Info().Msg("Table created")
		c.JSON(201, gin.H{"data": tableID})
	}
}

func UpdateTableHandler(tableService *service.TableService, logger zerolog.Logger) gin.HandlerFunc {
	service := *tableService

	return func(c *gin.Context) {
		logger.Debug().Msg("Updating table...")

		id := c.Param("id")
		i, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to parse ID")
			c.JSON(400, gin.H{"error": "Failed to parse ID"})
			return
		}

		var input UpdateTableInput
		if err := c.ShouldBindJSON(&input); err != nil {
			logger.Error().Err(err).Msg("Failed to bind JSON")
			c.JSON(400, gin.H{"error": "Failed to bind JSON"})
			return
		}

		updTable, err := service.UpdateTable(i, input.Name, input.Description, logger)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to update table")
			c.JSON(500, gin.H{"error": "Failed to update table"})
			return
		}

		logger.Info().Msg("Table updated")
		c.JSON(200, gin.H{"data": updTable})
	}
}

func DeleteTableHandler(tableService *service.TableService, logger zerolog.Logger) gin.HandlerFunc {
	service := *tableService

	return func(c *gin.Context) {
		logger.Debug().Msg("Deleting table...")

		id := c.Param("id")
		i, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to parse ID")
			c.JSON(400, gin.H{"error": "Failed to parse ID"})
			return
		}

		deleteID, err := service.DeleteTable(i, logger)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to delete table")
			c.JSON(500, gin.H{"error": "Failed to delete table"})
			return
		}

		logger.Info().Msg("Table deleted")
		c.JSON(204, gin.H{"data": deleteID})
	}
}

func GetAllRecordsByTableHandler(tableService *service.TableService, logger zerolog.Logger) gin.HandlerFunc {
	service := *tableService

	return func(c *gin.Context) {
		logger.Debug().Msg("Getting all records...")

		id := c.Param("id")
		i, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to parse ID")
			c.JSON(400, gin.H{"error": "Failed to parse ID"})
			return
		}

		records, err := service.GetAllRecordsByTable(i, logger)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to get all records")
			c.JSON(500, gin.H{"error": "Failed to get all records"})
			return
		}

		logger.Info().Msg("All records retrieved")
		c.JSON(200, gin.H{"data": records})
	}
}

func GetAllFieldsByTableHandler(tableService *service.TableService, logger zerolog.Logger) gin.HandlerFunc {
	service := *tableService

	return func(c *gin.Context) {
		logger.Debug().Msg("Getting all fields...")

		id := c.Param("id")
		i, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to parse ID")
			c.JSON(400, gin.H{"error": "Failed to parse ID"})
			return
		}

		fields, err := service.GetAllFieldsByTable(i, logger)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to get all fields")
			c.JSON(500, gin.H{"error": "Failed to get all fields"})
			return
		}

		logger.Info().Msg("All fields retrieved")
		c.JSON(200, gin.H{"data": fields})
	}
}
