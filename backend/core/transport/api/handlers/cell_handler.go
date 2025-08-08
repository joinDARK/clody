package handlers

import (
	"clody/core/service"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/datatypes"
)

type CreateCellInput struct {
	FieldID  int64          `json:"field_id" binding:"required"`
	RecordID int64          `json:"record_id" binding:"required"`
	Position int            `json:"position" binding:"required,min=0"`
	Value    datatypes.JSON `json:"value" binding:"omitempty"`
}

type UpdateCellInput struct {
	Position *int           `json:"position" binding:"omitempty,min=0"`
	Value    datatypes.JSON `json:"value" binding:"omitempty"`
}

type UpdateCellValueInput struct {
	Value datatypes.JSON `json:"value" binding:"required"`
}

type BulkUpdateCellsInput struct {
	Updates []struct {
		ID    int64          `json:"id" binding:"required"`
		Value datatypes.JSON `json:"value" binding:"required"`
	} `json:"updates" binding:"required"`
}

func GetCellByIDHandler(cellService *service.CellService, logger zerolog.Logger) gin.HandlerFunc {
	service := *cellService

	return func(c *gin.Context) {
		logger.Debug().Msg("[GetCellByIDHandler] Getting cell by ID")

		id := c.Param("id")
		cellID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Error().Err(err).Str("id", id).Msg("[GetCellByIDHandler] Failed to parse cell ID")
			c.JSON(400, gin.H{"error": "Invalid cell ID format"})
			return
		}

		cell, err := service.GetCellByID(cellID, logger)
		if err != nil {
			logger.Error().Err(err).Int64("cell_id", cellID).Msg("[GetCellByIDHandler] Failed to get cell")
			if err.Error() == fmt.Sprintf("cell with ID %d not found", cellID) {
				c.JSON(404, gin.H{"error": "Cell not found"})
				return
			}
			c.JSON(500, gin.H{"error": "Failed to retrieve cell"})
			return
		}

		logger.Info().Int64("cell_id", cellID).Msg("[GetCellByIDHandler] Cell retrieved successfully")
		c.JSON(200, gin.H{
			"data": cell,
			"_links": map[string]map[string]string{
				"self": {
					"href": fmt.Sprintf("/api/cells/%d", cell.ID),
				},
				"field": {
					"href": fmt.Sprintf("/api/fields/%d", cell.FieldID),
				},
				"record": {
					"href": fmt.Sprintf("/api/records/%d", cell.RecordID),
				},
				"update": {
					"href":   fmt.Sprintf("/api/cells/%d", cell.ID),
					"method": "PUT",
				},
				"delete": {
					"href":   fmt.Sprintf("/api/cells/%d", cell.ID),
					"method": "DELETE",
				},
			},
		})
	}
}

func CreateCellHandler(cellService *service.CellService, logger zerolog.Logger) gin.HandlerFunc {
	service := *cellService

	return func(c *gin.Context) {
		logger.Debug().Msg("[CreateCellHandler] Creating new cell")

		var input CreateCellInput
		if err := c.ShouldBindJSON(&input); err != nil {
			logger.Error().Err(err).Msg("[CreateCellHandler] Failed to bind JSON input")
			c.JSON(400, gin.H{"error": "Invalid input data", "details": err.Error()})
			return
		}

		cell, err := service.CreateCell(
			input.FieldID,
			input.RecordID,
			input.Position,
			input.Value,
			logger,
		)
		if err != nil {
			logger.Error().Err(err).
				Int64("field_id", input.FieldID).
				Int64("record_id", input.RecordID).
				Msg("[CreateCellHandler] Failed to create cell")
			c.JSON(500, gin.H{"error": "Failed to create cell"})
			return
		}

		logger.Info().
			Int64("cell_id", cell.ID).
			Int64("field_id", cell.FieldID).
			Int64("record_id", cell.RecordID).
			Msg("[CreateCellHandler] Cell created successfully")

		c.JSON(201, gin.H{
			"data": cell,
			"_links": map[string]map[string]string{
				"self": {
					"href": fmt.Sprintf("/api/cells/%d", cell.ID),
				},
				"field": {
					"href": fmt.Sprintf("/api/fields/%d", cell.FieldID),
				},
				"record": {
					"href": fmt.Sprintf("/api/records/%d", cell.RecordID),
				},
			},
		})
	}
}

func UpdateCellHandler(cellService *service.CellService, logger zerolog.Logger) gin.HandlerFunc {
	service := *cellService

	return func(c *gin.Context) {
		logger.Debug().Msg("[UpdateCellHandler] Updating cell")

		id := c.Param("id")
		cellID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Error().Err(err).Str("id", id).Msg("[UpdateCellHandler] Failed to parse cell ID")
			c.JSON(400, gin.H{"error": "Invalid cell ID format"})
			return
		}

		var input UpdateCellInput
		if err := c.ShouldBindJSON(&input); err != nil {
			logger.Error().Err(err).Msg("[UpdateCellHandler] Failed to bind JSON input")
			c.JSON(400, gin.H{"error": "Invalid input data", "details": err.Error()})
			return
		}

		cell, err := service.UpdateCell(
			cellID,
			input.Position,
			input.Value,
			logger,
		)
		if err != nil {
			logger.Error().Err(err).Int64("cell_id", cellID).Msg("[UpdateCellHandler] Failed to update cell")
			if err.Error() == fmt.Sprintf("cell with ID %d not found", cellID) {
				c.JSON(404, gin.H{"error": "Cell not found"})
				return
			}
			c.JSON(500, gin.H{"error": "Failed to update cell"})
			return
		}

		logger.Info().
			Int64("cell_id", cell.ID).
			Int64("field_id", cell.FieldID).
			Int64("record_id", cell.RecordID).
			Msg("[UpdateCellHandler] Cell updated successfully")

		c.JSON(200, gin.H{
			"data": cell,
			"_links": map[string]map[string]string{
				"self": {
					"href": fmt.Sprintf("/api/cells/%d", cell.ID),
				},
				"field": {
					"href": fmt.Sprintf("/api/fields/%d", cell.FieldID),
				},
				"record": {
					"href": fmt.Sprintf("/api/records/%d", cell.RecordID),
				},
			},
		})
	}
}

func UpdateCellValueHandler(cellService *service.CellService, logger zerolog.Logger) gin.HandlerFunc {
	service := *cellService

	return func(c *gin.Context) {
		logger.Debug().Msg("[UpdateCellValueHandler] Updating cell value")

		id := c.Param("id")
		cellID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Error().Err(err).Str("id", id).Msg("[UpdateCellValueHandler] Failed to parse cell ID")
			c.JSON(400, gin.H{"error": "Invalid cell ID format"})
			return
		}

		var input UpdateCellValueInput
		if err := c.ShouldBindJSON(&input); err != nil {
			logger.Error().Err(err).Msg("[UpdateCellValueHandler] Failed to bind JSON input")
			c.JSON(400, gin.H{"error": "Invalid input data", "details": err.Error()})
			return
		}

		cell, err := service.UpdateCellValue(cellID, input.Value, logger)
		if err != nil {
			logger.Error().Err(err).Int64("cell_id", cellID).Msg("[UpdateCellValueHandler] Failed to update cell value")
			if err.Error() == fmt.Sprintf("cell with ID %d not found", cellID) {
				c.JSON(404, gin.H{"error": "Cell not found"})
				return
			}
			c.JSON(500, gin.H{"error": "Failed to update cell value"})
			return
		}

		logger.Info().Int64("cell_id", cell.ID).Msg("[UpdateCellValueHandler] Cell value updated successfully")

		c.JSON(200, gin.H{
			"data": cell,
			"_links": map[string]map[string]string{
				"self": {
					"href": fmt.Sprintf("/api/cells/%d", cell.ID),
				},
				"field": {
					"href": fmt.Sprintf("/api/fields/%d", cell.FieldID),
				},
				"record": {
					"href": fmt.Sprintf("/api/records/%d", cell.RecordID),
				},
			},
		})
	}
}

func DeleteCellHandler(cellService *service.CellService, logger zerolog.Logger) gin.HandlerFunc {
	service := *cellService

	return func(c *gin.Context) {
		logger.Debug().Msg("[DeleteCellHandler] Deleting cell")

		id := c.Param("id")
		cellID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Error().Err(err).Str("id", id).Msg("[DeleteCellHandler] Failed to parse cell ID")
			c.JSON(400, gin.H{"error": "Invalid cell ID format"})
			return
		}

		err = service.DeleteCell(cellID, logger)
		if err != nil {
			logger.Error().Err(err).Int64("cell_id", cellID).Msg("[DeleteCellHandler] Failed to delete cell")
			if err.Error() == fmt.Sprintf("cell with ID %d not found", cellID) {
				c.JSON(404, gin.H{"error": "Cell not found"})
				return
			}
			c.JSON(500, gin.H{"error": "Failed to delete cell"})
			return
		}

		logger.Info().Int64("cell_id", cellID).Msg("[DeleteCellHandler] Cell deleted successfully")
		c.JSON(204, gin.H{"message": "Cell deleted successfully"})
	}
}

func GetCellsByRecordIDHandler(cellService *service.CellService, logger zerolog.Logger) gin.HandlerFunc {
	service := *cellService

	return func(c *gin.Context) {
		logger.Debug().Msg("[GetCellsByRecordIDHandler] Getting cells by record ID")

		recordID := c.Param("record_id")
		recordIDInt, err := strconv.ParseInt(recordID, 10, 64)
		if err != nil {
			logger.Error().Err(err).Str("record_id", recordID).Msg("[GetCellsByRecordIDHandler] Failed to parse record ID")
			c.JSON(400, gin.H{"error": "Invalid record ID format"})
			return
		}

		cells, err := service.GetCellsByRecordID(recordIDInt, logger)
		if err != nil {
			logger.Error().Err(err).Int64("record_id", recordIDInt).Msg("[GetCellsByRecordIDHandler] Failed to get cells")
			c.JSON(500, gin.H{"error": "Failed to retrieve cells"})
			return
		}

		logger.Info().
			Int64("record_id", recordIDInt).
			Int("count", len(cells)).
			Msg("[GetCellsByRecordIDHandler] Cells retrieved successfully")

		c.JSON(200, gin.H{
			"data":  cells,
			"count": len(cells),
			"_links": map[string]map[string]string{
				"record": {
					"href": fmt.Sprintf("/api/records/%d", recordIDInt),
				},
			},
		})
	}
}

func GetCellsByFieldIDHandler(cellService *service.CellService, logger zerolog.Logger) gin.HandlerFunc {
	service := *cellService

	return func(c *gin.Context) {
		logger.Debug().Msg("[GetCellsByFieldIDHandler] Getting cells by field ID")

		fieldID := c.Param("field_id")
		fieldIDInt, err := strconv.ParseInt(fieldID, 10, 64)
		if err != nil {
			logger.Error().Err(err).Str("field_id", fieldID).Msg("[GetCellsByFieldIDHandler] Failed to parse field ID")
			c.JSON(400, gin.H{"error": "Invalid field ID format"})
			return
		}

		cells, err := service.GetCellsByFieldID(fieldIDInt, logger)
		if err != nil {
			logger.Error().Err(err).Int64("field_id", fieldIDInt).Msg("[GetCellsByFieldIDHandler] Failed to get cells")
			c.JSON(500, gin.H{"error": "Failed to retrieve cells"})
			return
		}

		logger.Info().
			Int64("field_id", fieldIDInt).
			Int("count", len(cells)).
			Msg("[GetCellsByFieldIDHandler] Cells retrieved successfully")

		c.JSON(200, gin.H{
			"data":  cells,
			"count": len(cells),
			"_links": map[string]map[string]string{
				"field": {
					"href": fmt.Sprintf("/api/fields/%d", fieldIDInt),
				},
			},
		})
	}
}

func GetCellByFieldAndRecordHandler(cellService *service.CellService, logger zerolog.Logger) gin.HandlerFunc {
	service := *cellService

	return func(c *gin.Context) {
		logger.Debug().Msg("[GetCellByFieldAndRecordHandler] Getting cell by field and record ID")

		fieldID := c.Param("field_id")
		recordID := c.Param("record_id")

		fieldIDInt, err := strconv.ParseInt(fieldID, 10, 64)
		if err != nil {
			logger.Error().Err(err).Str("field_id", fieldID).Msg("[GetCellByFieldAndRecordHandler] Failed to parse field ID")
			c.JSON(400, gin.H{"error": "Invalid field ID format"})
			return
		}

		recordIDInt, err := strconv.ParseInt(recordID, 10, 64)
		if err != nil {
			logger.Error().Err(err).Str("record_id", recordID).Msg("[GetCellByFieldAndRecordHandler] Failed to parse record ID")
			c.JSON(400, gin.H{"error": "Invalid record ID format"})
			return
		}

		cell, err := service.GetCellByFieldAndRecord(fieldIDInt, recordIDInt, logger)
		if err != nil {
			logger.Error().Err(err).
				Int64("field_id", fieldIDInt).
				Int64("record_id", recordIDInt).
				Msg("[GetCellByFieldAndRecordHandler] Failed to get cell")
			if err.Error() == fmt.Sprintf("cell not found for field %d and record %d", fieldIDInt, recordIDInt) {
				c.JSON(404, gin.H{"error": "Cell not found"})
				return
			}
			c.JSON(500, gin.H{"error": "Failed to retrieve cell"})
			return
		}

		logger.Info().
			Int64("cell_id", cell.ID).
			Int64("field_id", fieldIDInt).
			Int64("record_id", recordIDInt).
			Msg("[GetCellByFieldAndRecordHandler] Cell retrieved successfully")

		c.JSON(200, gin.H{
			"data": cell,
			"_links": map[string]map[string]string{
				"self": {
					"href": fmt.Sprintf("/api/cells/%d", cell.ID),
				},
				"field": {
					"href": fmt.Sprintf("/api/fields/%d", cell.FieldID),
				},
				"record": {
					"href": fmt.Sprintf("/api/records/%d", cell.RecordID),
				},
			},
		})
	}
}

func BulkUpdateCellsHandler(cellService *service.CellService, logger zerolog.Logger) gin.HandlerFunc {
	service := *cellService

	return func(c *gin.Context) {
		logger.Debug().Msg("[BulkUpdateCellsHandler] Bulk updating cells")

		var input BulkUpdateCellsInput
		if err := c.ShouldBindJSON(&input); err != nil {
			logger.Error().Err(err).Msg("[BulkUpdateCellsHandler] Failed to bind JSON input")
			c.JSON(400, gin.H{"error": "Invalid input data", "details": err.Error()})
			return
		}

		if len(input.Updates) == 0 {
			logger.Warn().Msg("[BulkUpdateCellsHandler] No updates provided")
			c.JSON(400, gin.H{"error": "No updates provided"})
			return
		}

		// Convert input to service format
		updates := make([]struct {
			ID    int64
			Value datatypes.JSON
		}, len(input.Updates))

		for i, update := range input.Updates {
			updates[i] = struct {
				ID    int64
				Value datatypes.JSON
			}{
				ID:    update.ID,
				Value: update.Value,
			}
		}

		updatedCells, err := service.BulkUpdateCells(updates, logger)
		if err != nil {
			logger.Error().Err(err).Msg("[BulkUpdateCellsHandler] Failed to bulk update cells")
			c.JSON(500, gin.H{"error": "Failed to bulk update cells"})
			return
		}

		logger.Info().
			Int("requested", len(input.Updates)).
			Int("updated", len(updatedCells)).
			Msg("[BulkUpdateCellsHandler] Cells bulk updated successfully")

		c.JSON(200, gin.H{
			"data": gin.H{
				"updated_cells":   updatedCells,
				"requested_count": len(input.Updates),
				"updated_count":   len(updatedCells),
			},
			"message": fmt.Sprintf("Successfully updated %d out of %d cells", len(updatedCells), len(input.Updates)),
		})
	}
}
