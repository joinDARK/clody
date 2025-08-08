package handlers

import (
	"clody/core/service"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type CreateRecordRelationInput struct {
	FieldID        int64 `json:"field_id" binding:"required"`
	SourceRecordID int64 `json:"source_record_id" binding:"required"`
	TargetRecordID int64 `json:"target_record_id" binding:"required"`
}

type UpdateRecordRelationInput struct {
	FieldID        *int64 `json:"field_id" binding:"omitempty"`
	SourceRecordID *int64 `json:"source_record_id" binding:"omitempty"`
	TargetRecordID *int64 `json:"target_record_id" binding:"omitempty"`
}

func GetRecordRelationByIDHandler(recordRelationService *service.RecordRelationService, logger zerolog.Logger) gin.HandlerFunc {
	service := *recordRelationService

	return func(c *gin.Context) {
		logger.Debug().Msg("[GetRecordRelationByIDHandler] Getting record relation by ID")

		id := c.Param("id")
		recordRelationID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Error().Err(err).Str("id", id).Msg("[GetRecordRelationByIDHandler] Failed to parse record relation ID")
			c.JSON(400, gin.H{"error": "Invalid record relation ID format"})
			return
		}

		recordRelation, err := service.GetRecordRelationByID(recordRelationID, logger)
		if err != nil {
			logger.Error().Err(err).Int64("record_relation_id", recordRelationID).Msg("[GetRecordRelationByIDHandler] Failed to get record relation")
			if err.Error() == fmt.Sprintf("record relation with ID %d not found", recordRelationID) {
				c.JSON(404, gin.H{"error": "Record relation not found"})
				return
			}
			c.JSON(500, gin.H{"error": "Failed to retrieve record relation"})
			return
		}

		logger.Info().Int64("record_relation_id", recordRelationID).Msg("[GetRecordRelationByIDHandler] Record relation retrieved successfully")
		c.JSON(200, gin.H{
			"data": recordRelation,
			"_links": map[string]map[string]string{
				"self": {
					"href": fmt.Sprintf("/api/record-relations/%d", recordRelation.ID),
				},
				"field": {
					"href": fmt.Sprintf("/api/fields/%d", recordRelation.FieldID),
				},
				"source_record": {
					"href": fmt.Sprintf("/api/records/%d", recordRelation.SourceRecordID),
				},
				"target_record": {
					"href": fmt.Sprintf("/api/records/%d", recordRelation.TargetRecordID),
				},
				"update": {
					"href":   fmt.Sprintf("/api/record-relations/%d", recordRelation.ID),
					"method": "PUT",
				},
				"delete": {
					"href":   fmt.Sprintf("/api/record-relations/%d", recordRelation.ID),
					"method": "DELETE",
				},
			},
		})
	}
}

func CreateRecordRelationHandler(recordRelationService *service.RecordRelationService, logger zerolog.Logger) gin.HandlerFunc {
	service := *recordRelationService

	return func(c *gin.Context) {
		logger.Debug().Msg("[CreateRecordRelationHandler] Creating new record relation")

		var input CreateRecordRelationInput
		if err := c.ShouldBindJSON(&input); err != nil {
			logger.Error().Err(err).Msg("[CreateRecordRelationHandler] Failed to bind JSON input")
			c.JSON(400, gin.H{"error": "Invalid input data", "details": err.Error()})
			return
		}

		recordRelation, err := service.CreateRecordRelation(
			input.FieldID,
			input.SourceRecordID,
			input.TargetRecordID,
			logger,
		)
		if err != nil {
			logger.Error().Err(err).
				Int64("field_id", input.FieldID).
				Int64("source_record_id", input.SourceRecordID).
				Int64("target_record_id", input.TargetRecordID).
				Msg("[CreateRecordRelationHandler] Failed to create record relation")

			// Handle specific error cases
			if err.Error() == "record cannot be related to itself" {
				c.JSON(400, gin.H{"error": "Record cannot be related to itself"})
				return
			}
			if fmt.Sprintf("relation between records %d and %d already exists", input.SourceRecordID, input.TargetRecordID) == err.Error() {
				c.JSON(409, gin.H{"error": "Relation already exists"})
				return
			}

			c.JSON(500, gin.H{"error": "Failed to create record relation"})
			return
		}

		logger.Info().
			Int64("record_relation_id", recordRelation.ID).
			Int64("field_id", recordRelation.FieldID).
			Int64("source_record_id", recordRelation.SourceRecordID).
			Int64("target_record_id", recordRelation.TargetRecordID).
			Msg("[CreateRecordRelationHandler] Record relation created successfully")

		c.JSON(201, gin.H{
			"data": recordRelation,
			"_links": map[string]map[string]string{
				"self": {
					"href": fmt.Sprintf("/api/record-relations/%d", recordRelation.ID),
				},
				"field": {
					"href": fmt.Sprintf("/api/fields/%d", recordRelation.FieldID),
				},
				"source_record": {
					"href": fmt.Sprintf("/api/records/%d", recordRelation.SourceRecordID),
				},
				"target_record": {
					"href": fmt.Sprintf("/api/records/%d", recordRelation.TargetRecordID),
				},
			},
		})
	}
}

func UpdateRecordRelationHandler(recordRelationService *service.RecordRelationService, logger zerolog.Logger) gin.HandlerFunc {
	service := *recordRelationService

	return func(c *gin.Context) {
		logger.Debug().Msg("[UpdateRecordRelationHandler] Updating record relation")

		id := c.Param("id")
		recordRelationID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Error().Err(err).Str("id", id).Msg("[UpdateRecordRelationHandler] Failed to parse record relation ID")
			c.JSON(400, gin.H{"error": "Invalid record relation ID format"})
			return
		}

		var input UpdateRecordRelationInput
		if err := c.ShouldBindJSON(&input); err != nil {
			logger.Error().Err(err).Msg("[UpdateRecordRelationHandler] Failed to bind JSON input")
			c.JSON(400, gin.H{"error": "Invalid input data", "details": err.Error()})
			return
		}

		recordRelation, err := service.UpdateRecordRelation(
			recordRelationID,
			input.FieldID,
			input.SourceRecordID,
			input.TargetRecordID,
			logger,
		)
		if err != nil {
			logger.Error().Err(err).Int64("record_relation_id", recordRelationID).Msg("[UpdateRecordRelationHandler] Failed to update record relation")
			if err.Error() == fmt.Sprintf("record relation with ID %d not found", recordRelationID) {
				c.JSON(404, gin.H{"error": "Record relation not found"})
				return
			}
			if err.Error() == "record cannot be related to itself" {
				c.JSON(400, gin.H{"error": "Record cannot be related to itself"})
				return
			}
			c.JSON(500, gin.H{"error": "Failed to update record relation"})
			return
		}

		logger.Info().
			Int64("record_relation_id", recordRelation.ID).
			Int64("field_id", recordRelation.FieldID).
			Int64("source_record_id", recordRelation.SourceRecordID).
			Int64("target_record_id", recordRelation.TargetRecordID).
			Msg("[UpdateRecordRelationHandler] Record relation updated successfully")

		c.JSON(200, gin.H{
			"data": recordRelation,
			"_links": map[string]map[string]string{
				"self": {
					"href": fmt.Sprintf("/api/record-relations/%d", recordRelation.ID),
				},
				"field": {
					"href": fmt.Sprintf("/api/fields/%d", recordRelation.FieldID),
				},
				"source_record": {
					"href": fmt.Sprintf("/api/records/%d", recordRelation.SourceRecordID),
				},
				"target_record": {
					"href": fmt.Sprintf("/api/records/%d", recordRelation.TargetRecordID),
				},
			},
		})
	}
}

func DeleteRecordRelationHandler(recordRelationService *service.RecordRelationService, logger zerolog.Logger) gin.HandlerFunc {
	service := *recordRelationService

	return func(c *gin.Context) {
		logger.Debug().Msg("[DeleteRecordRelationHandler] Deleting record relation")

		id := c.Param("id")
		recordRelationID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Error().Err(err).Str("id", id).Msg("[DeleteRecordRelationHandler] Failed to parse record relation ID")
			c.JSON(400, gin.H{"error": "Invalid record relation ID format"})
			return
		}

		err = service.DeleteRecordRelation(recordRelationID, logger)
		if err != nil {
			logger.Error().Err(err).Int64("record_relation_id", recordRelationID).Msg("[DeleteRecordRelationHandler] Failed to delete record relation")
			if err.Error() == fmt.Sprintf("record relation with ID %d not found", recordRelationID) {
				c.JSON(404, gin.H{"error": "Record relation not found"})
				return
			}
			c.JSON(500, gin.H{"error": "Failed to delete record relation"})
			return
		}

		logger.Info().Int64("record_relation_id", recordRelationID).Msg("[DeleteRecordRelationHandler] Record relation deleted successfully")
		c.JSON(204, gin.H{"message": "Record relation deleted successfully"})
	}
}

func GetRecordRelationsByFieldIDHandler(recordRelationService *service.RecordRelationService, logger zerolog.Logger) gin.HandlerFunc {
	service := *recordRelationService

	return func(c *gin.Context) {
		logger.Debug().Msg("[GetRecordRelationsByFieldIDHandler] Getting record relations by field ID")

		fieldID := c.Param("field_id")
		fieldIDInt, err := strconv.ParseInt(fieldID, 10, 64)
		if err != nil {
			logger.Error().Err(err).Str("field_id", fieldID).Msg("[GetRecordRelationsByFieldIDHandler] Failed to parse field ID")
			c.JSON(400, gin.H{"error": "Invalid field ID format"})
			return
		}

		recordRelations, err := service.GetRecordRelationsByFieldID(fieldIDInt, logger)
		if err != nil {
			logger.Error().Err(err).Int64("field_id", fieldIDInt).Msg("[GetRecordRelationsByFieldIDHandler] Failed to get record relations")
			c.JSON(500, gin.H{"error": "Failed to retrieve record relations"})
			return
		}

		logger.Info().
			Int64("field_id", fieldIDInt).
			Int("count", len(recordRelations)).
			Msg("[GetRecordRelationsByFieldIDHandler] Record relations retrieved successfully")

		c.JSON(200, gin.H{
			"data":  recordRelations,
			"count": len(recordRelations),
			"_links": map[string]map[string]string{
				"field": {
					"href": fmt.Sprintf("/api/fields/%d", fieldIDInt),
				},
			},
		})
	}
}

func GetRecordRelationsBySourceRecordIDHandler(recordRelationService *service.RecordRelationService, logger zerolog.Logger) gin.HandlerFunc {
	service := *recordRelationService

	return func(c *gin.Context) {
		logger.Debug().Msg("[GetRecordRelationsBySourceRecordIDHandler] Getting record relations by source record ID")

		recordID := c.Param("record_id")
		recordIDInt, err := strconv.ParseInt(recordID, 10, 64)
		if err != nil {
			logger.Error().Err(err).Str("record_id", recordID).Msg("[GetRecordRelationsBySourceRecordIDHandler] Failed to parse record ID")
			c.JSON(400, gin.H{"error": "Invalid record ID format"})
			return
		}

		recordRelations, err := service.GetRecordRelationsBySourceRecordID(recordIDInt, logger)
		if err != nil {
			logger.Error().Err(err).Int64("source_record_id", recordIDInt).Msg("[GetRecordRelationsBySourceRecordIDHandler] Failed to get record relations")
			c.JSON(500, gin.H{"error": "Failed to retrieve record relations"})
			return
		}

		logger.Info().
			Int64("source_record_id", recordIDInt).
			Int("count", len(recordRelations)).
			Msg("[GetRecordRelationsBySourceRecordIDHandler] Record relations retrieved successfully")

		c.JSON(200, gin.H{
			"data":  recordRelations,
			"count": len(recordRelations),
			"_links": map[string]map[string]string{
				"source_record": {
					"href": fmt.Sprintf("/api/records/%d", recordIDInt),
				},
			},
		})
	}
}

func GetRecordRelationsByTargetRecordIDHandler(recordRelationService *service.RecordRelationService, logger zerolog.Logger) gin.HandlerFunc {
	service := *recordRelationService

	return func(c *gin.Context) {
		logger.Debug().Msg("[GetRecordRelationsByTargetRecordIDHandler] Getting record relations by target record ID")

		recordID := c.Param("record_id")
		recordIDInt, err := strconv.ParseInt(recordID, 10, 64)
		if err != nil {
			logger.Error().Err(err).Str("record_id", recordID).Msg("[GetRecordRelationsByTargetRecordIDHandler] Failed to parse record ID")
			c.JSON(400, gin.H{"error": "Invalid record ID format"})
			return
		}

		recordRelations, err := service.GetRecordRelationsByTargetRecordID(recordIDInt, logger)
		if err != nil {
			logger.Error().Err(err).Int64("target_record_id", recordIDInt).Msg("[GetRecordRelationsByTargetRecordIDHandler] Failed to get record relations")
			c.JSON(500, gin.H{"error": "Failed to retrieve record relations"})
			return
		}

		logger.Info().
			Int64("target_record_id", recordIDInt).
			Int("count", len(recordRelations)).
			Msg("[GetRecordRelationsByTargetRecordIDHandler] Record relations retrieved successfully")

		c.JSON(200, gin.H{
			"data":  recordRelations,
			"count": len(recordRelations),
			"_links": map[string]map[string]string{
				"target_record": {
					"href": fmt.Sprintf("/api/records/%d", recordIDInt),
				},
			},
		})
	}
}
