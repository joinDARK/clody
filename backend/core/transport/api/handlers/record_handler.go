package handlers

import (
	"clody/core/service"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type CreateRecordInput struct {
	TableID  int64 `json:"table_id" binding:"required"`
	Position int   `json:"position" binding:"required,min=1"`
}

type UpdateRecordInput struct {
	Position *int `json:"position" binding:"omitempty,min=1"`
}

func GetRecordByIDHandler(recordService *service.RecordService, logger zerolog.Logger) gin.HandlerFunc {
	service := *recordService

	return func(c *gin.Context) {
		logger.Debug().Msg("[GetRecordByIDHandler] Getting record by ID")

		id := c.Param("id")
		recordID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Error().Err(err).Str("id", id).Msg("[GetRecordByIDHandler] Failed to parse record ID")
			c.JSON(400, gin.H{"error": "Invalid record ID format"})
			return
		}

		record, err := service.GetRecordByID(recordID, logger)
		if err != nil {
			logger.Error().Err(err).Int64("record_id", recordID).Msg("[GetRecordByIDHandler] Failed to get record")
			if err.Error() == fmt.Sprintf("record with ID %d not found", recordID) {
				c.JSON(404, gin.H{"error": "Record not found"})
				return
			}
			c.JSON(500, gin.H{"error": "Failed to retrieve record"})
			return
		}

		// Check if additional data is requested
		_, withCells := c.GetQuery("with_cells")
		_, withSourceRecords := c.GetQuery("with_source_records")
		_, withTargetRecords := c.GetQuery("with_target_records")

		responseData := gin.H{
			"record": record,
		}

		if withCells {
			cells, err := service.GetRecordCells(recordID, logger)
			if err != nil {
				logger.Error().Err(err).Int64("record_id", recordID).Msg("[GetRecordByIDHandler] Failed to get record cells")
			} else {
				responseData["cells"] = cells
				responseData["cells_count"] = len(cells)
			}
		}

		if withSourceRecords {
			sourceRecords, err := service.GetRecordsBySourceID(recordID, logger)
			if err != nil {
				logger.Error().Err(err).Int64("record_id", recordID).Msg("[GetRecordByIDHandler] Failed to get source records")
			} else {
				responseData["source_records"] = sourceRecords
				responseData["source_records_count"] = len(sourceRecords)
			}
		}

		if withTargetRecords {
			targetRecords, err := service.GetRecordsByTargetID(recordID, logger)
			if err != nil {
				logger.Error().Err(err).Int64("record_id", recordID).Msg("[GetRecordByIDHandler] Failed to get target records")
			} else {
				responseData["target_records"] = targetRecords
				responseData["target_records_count"] = len(targetRecords)
			}
		}

		logger.Info().Int64("record_id", recordID).Msg("[GetRecordByIDHandler] Record retrieved successfully")
		c.JSON(200, gin.H{
			"data": responseData,
			"_links": map[string]map[string]string{
				"self": {
					"href":   fmt.Sprintf("/api/records/%d", record.ID),
					"params": "with_cells,with_source_records,with_target_records",
				},
				"table": {
					"href": fmt.Sprintf("/api/tables/%d", record.TableID),
				},
				"update": {
					"href":   fmt.Sprintf("/api/records/%d", record.ID),
					"method": "PUT",
				},
				"delete": {
					"href":   fmt.Sprintf("/api/records/%d", record.ID),
					"method": "DELETE",
				},
				"cells": {
					"href": fmt.Sprintf("/api/records/%d/cells", record.ID),
				},
				"source_records": {
					"href": fmt.Sprintf("/api/records/%d/source-records", record.ID),
				},
				"target_records": {
					"href": fmt.Sprintf("/api/records/%d/target-records", record.ID),
				},
			},
		})
	}
}

func CreateRecordHandler(recordService *service.RecordService, logger zerolog.Logger) gin.HandlerFunc {
	service := *recordService

	return func(c *gin.Context) {
		logger.Debug().Msg("[CreateRecordHandler] Creating new record")

		var input CreateRecordInput
		if err := c.ShouldBindJSON(&input); err != nil {
			logger.Error().Err(err).Msg("[CreateRecordHandler] Failed to bind JSON input")
			c.JSON(400, gin.H{"error": "Invalid input data", "details": err.Error()})
			return
		}

		record, err := service.CreateRecord(
			input.TableID,
			input.Position,
			logger,
		)
		if err != nil {
			logger.Error().Err(err).
				Int64("table_id", input.TableID).
				Int("position", input.Position).
				Msg("[CreateRecordHandler] Failed to create record")
			c.JSON(500, gin.H{"error": "Failed to create record"})
			return
		}

		logger.Info().
			Int64("record_id", record.ID).
			Int64("table_id", record.TableID).
			Int("position", record.Position).
			Msg("[CreateRecordHandler] Record created successfully")

		c.JSON(201, gin.H{
			"data": record,
			"_links": map[string]map[string]string{
				"self": {
					"href": fmt.Sprintf("/api/records/%d", record.ID),
				},
				"table": {
					"href": fmt.Sprintf("/api/tables/%d", record.TableID),
				},
			},
		})
	}
}

func UpdateRecordHandler(recordService *service.RecordService, logger zerolog.Logger) gin.HandlerFunc {
	service := *recordService

	return func(c *gin.Context) {
		logger.Debug().Msg("[UpdateRecordHandler] Updating record")

		id := c.Param("id")
		recordID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Error().Err(err).Str("id", id).Msg("[UpdateRecordHandler] Failed to parse record ID")
			c.JSON(400, gin.H{"error": "Invalid record ID format"})
			return
		}

		var input UpdateRecordInput
		if err := c.ShouldBindJSON(&input); err != nil {
			logger.Error().Err(err).Msg("[UpdateRecordHandler] Failed to bind JSON input")
			c.JSON(400, gin.H{"error": "Invalid input data", "details": err.Error()})
			return
		}

		record, err := service.UpdateRecord(
			recordID,
			input.Position,
			logger,
		)
		if err != nil {
			logger.Error().Err(err).Int64("record_id", recordID).Msg("[UpdateRecordHandler] Failed to update record")
			if err.Error() == fmt.Sprintf("record with ID %d not found", recordID) {
				c.JSON(404, gin.H{"error": "Record not found"})
				return
			}
			c.JSON(500, gin.H{"error": "Failed to update record"})
			return
		}

		logger.Info().
			Int64("record_id", record.ID).
			Int64("table_id", record.TableID).
			Int("position", record.Position).
			Msg("[UpdateRecordHandler] Record updated successfully")

		c.JSON(200, gin.H{
			"data": record,
			"_links": map[string]map[string]string{
				"self": {
					"href": fmt.Sprintf("/api/records/%d", record.ID),
				},
				"table": {
					"href": fmt.Sprintf("/api/tables/%d", record.TableID),
				},
			},
		})
	}
}

func DeleteRecordHandler(recordService *service.RecordService, logger zerolog.Logger) gin.HandlerFunc {
	service := *recordService

	return func(c *gin.Context) {
		logger.Debug().Msg("[DeleteRecordHandler] Deleting record")

		id := c.Param("id")
		recordID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Error().Err(err).Str("id", id).Msg("[DeleteRecordHandler] Failed to parse record ID")
			c.JSON(400, gin.H{"error": "Invalid record ID format"})
			return
		}

		err = service.DeleteRecord(recordID, logger)
		if err != nil {
			logger.Error().Err(err).Int64("record_id", recordID).Msg("[DeleteRecordHandler] Failed to delete record")
			if err.Error() == fmt.Sprintf("record with ID %d not found", recordID) {
				c.JSON(404, gin.H{"error": "Record not found"})
				return
			}
			c.JSON(500, gin.H{"error": "Failed to delete record"})
			return
		}

		logger.Info().Int64("record_id", recordID).Msg("[DeleteRecordHandler] Record deleted successfully")
		c.JSON(204, gin.H{"message": "Record deleted successfully"})
	}
}

func GetRecordCellsHandler(recordService *service.RecordService, logger zerolog.Logger) gin.HandlerFunc {
	service := *recordService

	return func(c *gin.Context) {
		logger.Debug().Msg("[GetRecordCellsHandler] Getting record cells")

		id := c.Param("id")
		recordID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Error().Err(err).Str("id", id).Msg("[GetRecordCellsHandler] Failed to parse record ID")
			c.JSON(400, gin.H{"error": "Invalid record ID format"})
			return
		}

		cells, err := service.GetRecordCells(recordID, logger)
		if err != nil {
			logger.Error().Err(err).Int64("record_id", recordID).Msg("[GetRecordCellsHandler] Failed to get record cells")
			if err.Error() == fmt.Sprintf("record with ID %d does not exist", recordID) {
				c.JSON(404, gin.H{"error": "Record not found"})
				return
			}
			c.JSON(500, gin.H{"error": "Failed to retrieve record cells"})
			return
		}

		logger.Info().
			Int64("record_id", recordID).
			Int("count", len(cells)).
			Msg("[GetRecordCellsHandler] Record cells retrieved successfully")

		c.JSON(200, gin.H{
			"data":  cells,
			"count": len(cells),
			"_links": map[string]map[string]string{
				"record": {
					"href": fmt.Sprintf("/api/records/%d", recordID),
				},
			},
		})
	}
}

func GetRecordSourceRecordsHandler(recordService *service.RecordService, logger zerolog.Logger) gin.HandlerFunc {
	service := *recordService

	return func(c *gin.Context) {
		logger.Debug().Msg("[GetRecordSourceRecordsHandler] Getting record source records")

		id := c.Param("id")
		recordID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Error().Err(err).Str("id", id).Msg("[GetRecordSourceRecordsHandler] Failed to parse record ID")
			c.JSON(400, gin.H{"error": "Invalid record ID format"})
			return
		}

		sourceRecords, err := service.GetRecordsBySourceID(recordID, logger)
		if err != nil {
			logger.Error().Err(err).Int64("record_id", recordID).Msg("[GetRecordSourceRecordsHandler] Failed to get source records")
			if err.Error() == fmt.Sprintf("record with ID %d does not exist", recordID) {
				c.JSON(404, gin.H{"error": "Record not found"})
				return
			}
			c.JSON(500, gin.H{"error": "Failed to retrieve source records"})
			return
		}

		logger.Info().
			Int64("record_id", recordID).
			Int("count", len(sourceRecords)).
			Msg("[GetRecordSourceRecordsHandler] Source records retrieved successfully")

		c.JSON(200, gin.H{
			"data":  sourceRecords,
			"count": len(sourceRecords),
			"_links": map[string]map[string]string{
				"record": {
					"href": fmt.Sprintf("/api/records/%d", recordID),
				},
			},
		})
	}
}

func GetRecordTargetRecordsHandler(recordService *service.RecordService, logger zerolog.Logger) gin.HandlerFunc {
	service := *recordService

	return func(c *gin.Context) {
		logger.Debug().Msg("[GetRecordTargetRecordsHandler] Getting record target records")

		id := c.Param("id")
		recordID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Error().Err(err).Str("id", id).Msg("[GetRecordTargetRecordsHandler] Failed to parse record ID")
			c.JSON(400, gin.H{"error": "Invalid record ID format"})
			return
		}

		targetRecords, err := service.GetRecordsByTargetID(recordID, logger)
		if err != nil {
			logger.Error().Err(err).Int64("record_id", recordID).Msg("[GetRecordTargetRecordsHandler] Failed to get target records")
			if err.Error() == fmt.Sprintf("record with ID %d does not exist", recordID) {
				c.JSON(404, gin.H{"error": "Record not found"})
				return
			}
			c.JSON(500, gin.H{"error": "Failed to retrieve target records"})
			return
		}

		logger.Info().
			Int64("record_id", recordID).
			Int("count", len(targetRecords)).
			Msg("[GetRecordTargetRecordsHandler] Target records retrieved successfully")

		c.JSON(200, gin.H{
			"data":  targetRecords,
			"count": len(targetRecords),
			"_links": map[string]map[string]string{
				"record": {
					"href": fmt.Sprintf("/api/records/%d", recordID),
				},
			},
		})
	}
}
