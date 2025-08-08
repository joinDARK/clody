package handlers

import (
	"clody/core/domain"
	"clody/core/service"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/datatypes"
)

type CreateFieldInput struct {
	TableID     int64            `json:"table_id" binding:"required"`
	Name        string           `json:"name" binding:"required,min=1,max=255"`
	Type        domain.FieldType `json:"type" binding:"required"`
	Position    int              `json:"position" binding:"required,min=0"`
	Description *string          `json:"description" binding:"omitempty"`
	Options     datatypes.JSON   `json:"options" binding:"omitempty"`
}

type UpdateFieldInput struct {
	Name        string           `json:"name" binding:"omitempty,min=1,max=255"`
	Type        domain.FieldType `json:"type" binding:"omitempty"`
	Position    int              `json:"position" binding:"omitempty,min=0"`
	Description *string          `json:"description" binding:"omitempty"`
	Options     datatypes.JSON   `json:"options" binding:"omitempty"`
}

func GetFieldByIDHandler(fieldService *service.FieldService, logger zerolog.Logger) gin.HandlerFunc {
	service := *fieldService

	return func(c *gin.Context) {
		logger.Debug().Msg("[GetFieldByIDHandler] Getting field by ID")

		id := c.Param("id")
		fieldID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Error().Err(err).Str("id", id).Msg("[GetFieldByIDHandler] Failed to parse field ID")
			c.JSON(400, gin.H{"error": "Invalid field ID format"})
			return
		}

		field, err := service.GetFieldByID(fieldID, logger)
		if err != nil {
			logger.Error().Err(err).Int64("field_id", fieldID).Msg("[GetFieldByIDHandler] Failed to get field")
			if err.Error() == fmt.Sprintf("field with ID %d not found", fieldID) {
				c.JSON(404, gin.H{"error": "Field not found"})
				return
			}
			c.JSON(500, gin.H{"error": "Failed to retrieve field"})
			return
		}

		// Check if additional data is requested
		_, withCells := c.GetQuery("with_cells")
		_, withOptions := c.GetQuery("with_options")
		_, withRelations := c.GetQuery("with_relations")

		if withCells {
			cells, err := service.GetFieldCells(fieldID, logger)
			if err != nil {
				logger.Error().Err(err).Int64("field_id", fieldID).Msg("[GetFieldByIDHandler] Failed to get field cells")
			} else {
				// Add cells to response (would need to add this to the field struct or response)
				logger.Debug().Int("cells_count", len(cells)).Msg("[GetFieldByIDHandler] Retrieved field cells")
			}
		}

		if withOptions {
			options, err := service.GetFieldSelectOptions(fieldID, logger)
			if err != nil {
				logger.Error().Err(err).Int64("field_id", fieldID).Msg("[GetFieldByIDHandler] Failed to get field options")
			} else {
				field.SelectOptions = options
			}
		}

		if withRelations {
			relations, err := service.GetFieldRelationRecords(fieldID, logger)
			if err != nil {
				logger.Error().Err(err).Int64("field_id", fieldID).Msg("[GetFieldByIDHandler] Failed to get field relations")
			} else {
				logger.Debug().Int("relations_count", len(relations)).Msg("[GetFieldByIDHandler] Retrieved field relations")
			}
		}

		logger.Info().Int64("field_id", fieldID).Msg("[GetFieldByIDHandler] Field retrieved successfully")
		c.JSON(200, gin.H{
			"data": field,
			"_links": map[string]map[string]string{
				"self": {
					"href":   fmt.Sprintf("/api/fields/%d", field.ID),
					"params": "with_cells,with_options,with_relations",
				},
				"table": {
					"href": fmt.Sprintf("/api/tables/%d", field.TableID),
				},
				"update": {
					"href":   fmt.Sprintf("/api/fields/%d", field.ID),
					"method": "PUT",
				},
				"delete": {
					"href":   fmt.Sprintf("/api/fields/%d", field.ID),
					"method": "DELETE",
				},
			},
		})
	}
}

func CreateFieldHandler(fieldService *service.FieldService, logger zerolog.Logger) gin.HandlerFunc {
	service := *fieldService

	return func(c *gin.Context) {
		logger.Debug().Msg("[CreateFieldHandler] Creating new field")

		var input CreateFieldInput
		if err := c.ShouldBindJSON(&input); err != nil {
			logger.Error().Err(err).Msg("[CreateFieldHandler] Failed to bind JSON input")
			c.JSON(400, gin.H{"error": "Invalid input data", "details": err.Error()})
			return
		}

		// Validate field type
		if !input.Type.IsValid() {
			logger.Error().Str("type", string(input.Type)).Msg("[CreateFieldHandler] Invalid field type")
			c.JSON(400, gin.H{"error": "Invalid field type"})
			return
		}

		field, err := service.CreateField(
			input.TableID,
			input.Name,
			input.Type,
			input.Position,
			input.Description,
			input.Options,
			logger,
		)
		if err != nil {
			logger.Error().Err(err).
				Str("name", input.Name).
				Int64("table_id", input.TableID).
				Msg("[CreateFieldHandler] Failed to create field")
			c.JSON(500, gin.H{"error": "Failed to create field"})
			return
		}

		logger.Info().
			Int64("field_id", field.ID).
			Str("name", field.Name).
			Str("type", string(field.Type)).
			Msg("[CreateFieldHandler] Field created successfully")

		c.JSON(201, gin.H{
			"data": field,
			"_links": map[string]map[string]string{
				"self": {
					"href": fmt.Sprintf("/api/fields/%d", field.ID),
				},
				"table": {
					"href": fmt.Sprintf("/api/tables/%d", field.TableID),
				},
			},
		})
	}
}

func UpdateFieldHandler(fieldService *service.FieldService, logger zerolog.Logger) gin.HandlerFunc {
	service := *fieldService

	return func(c *gin.Context) {
		logger.Debug().Msg("[UpdateFieldHandler] Updating field")

		id := c.Param("id")
		fieldID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Error().Err(err).Str("id", id).Msg("[UpdateFieldHandler] Failed to parse field ID")
			c.JSON(400, gin.H{"error": "Invalid field ID format"})
			return
		}

		var input UpdateFieldInput
		if err := c.ShouldBindJSON(&input); err != nil {
			logger.Error().Err(err).Msg("[UpdateFieldHandler] Failed to bind JSON input")
			c.JSON(400, gin.H{"error": "Invalid input data", "details": err.Error()})
			return
		}

		// Validate field type if provided
		if input.Type != "" && !input.Type.IsValid() {
			logger.Error().Str("type", string(input.Type)).Msg("[UpdateFieldHandler] Invalid field type")
			c.JSON(400, gin.H{"error": "Invalid field type"})
			return
		}

		field, err := service.UpdateField(
			fieldID,
			input.Name,
			input.Type,
			input.Position,
			input.Description,
			input.Options,
			logger,
		)
		if err != nil {
			logger.Error().Err(err).Int64("field_id", fieldID).Msg("[UpdateFieldHandler] Failed to update field")
			if err.Error() == fmt.Sprintf("field with ID %d not found", fieldID) {
				c.JSON(404, gin.H{"error": "Field not found"})
				return
			}
			c.JSON(500, gin.H{"error": "Failed to update field"})
			return
		}

		logger.Info().
			Int64("field_id", field.ID).
			Str("name", field.Name).
			Msg("[UpdateFieldHandler] Field updated successfully")

		c.JSON(200, gin.H{
			"data": field,
			"_links": map[string]map[string]string{
				"self": {
					"href": fmt.Sprintf("/api/fields/%d", field.ID),
				},
				"table": {
					"href": fmt.Sprintf("/api/tables/%d", field.TableID),
				},
			},
		})
	}
}

func DeleteFieldHandler(fieldService *service.FieldService, logger zerolog.Logger) gin.HandlerFunc {
	service := *fieldService

	return func(c *gin.Context) {
		logger.Debug().Msg("[DeleteFieldHandler] Deleting field")

		id := c.Param("id")
		fieldID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Error().Err(err).Str("id", id).Msg("[DeleteFieldHandler] Failed to parse field ID")
			c.JSON(400, gin.H{"error": "Invalid field ID format"})
			return
		}

		deleteID, err := service.DeleteField(fieldID, logger)
		if err != nil {
			logger.Error().Err(err).Int64("field_id", fieldID).Msg("[DeleteFieldHandler] Failed to delete field")
			if err.Error() == fmt.Sprintf("field with ID %d not found", fieldID) {
				c.JSON(404, gin.H{"error": "Field not found"})
				return
			}
			c.JSON(500, gin.H{"error": "Failed to delete field"})
			return
		}

		logger.Info().Int64("field_id", fieldID).Msg("[DeleteFieldHandler] Field deleted successfully")
		c.JSON(204, gin.H{"data": deleteID})
	}
}

func GetFieldCellsHandler(fieldService *service.FieldService, logger zerolog.Logger) gin.HandlerFunc {
	service := *fieldService

	return func(c *gin.Context) {
		logger.Debug().Msg("[GetFieldCellsHandler] Getting field cells")

		id := c.Param("id")
		fieldID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Error().Err(err).Str("id", id).Msg("[GetFieldCellsHandler] Failed to parse field ID")
			c.JSON(400, gin.H{"error": "Invalid field ID format"})
			return
		}

		cells, err := service.GetFieldCells(fieldID, logger)
		if err != nil {
			logger.Error().Err(err).Int64("field_id", fieldID).Msg("[GetFieldCellsHandler] Failed to get field cells")
			if err.Error() == fmt.Sprintf("field with ID %d not found", fieldID) {
				c.JSON(404, gin.H{"error": "Field not found"})
				return
			}
			c.JSON(500, gin.H{"error": "Failed to retrieve field cells"})
			return
		}

		logger.Info().
			Int64("field_id", fieldID).
			Int("count", len(cells)).
			Msg("[GetFieldCellsHandler] Field cells retrieved successfully")

		c.JSON(200, gin.H{
			"data":  cells,
			"count": len(cells),
			"_links": map[string]map[string]string{
				"field": {
					"href": fmt.Sprintf("/api/fields/%d", fieldID),
				},
			},
		})
	}
}

func GetFieldSelectOptionsHandler(fieldService *service.FieldService, logger zerolog.Logger) gin.HandlerFunc {
	service := *fieldService

	return func(c *gin.Context) {
		logger.Debug().Msg("[GetFieldSelectOptionsHandler] Getting field select options")

		id := c.Param("id")
		fieldID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Error().Err(err).Str("id", id).Msg("[GetFieldSelectOptionsHandler] Failed to parse field ID")
			c.JSON(400, gin.H{"error": "Invalid field ID format"})
			return
		}

		options, err := service.GetFieldSelectOptions(fieldID, logger)
		if err != nil {
			logger.Error().Err(err).Int64("field_id", fieldID).Msg("[GetFieldSelectOptionsHandler] Failed to get field select options")
			if err.Error() == fmt.Sprintf("field with ID %d not found", fieldID) {
				c.JSON(404, gin.H{"error": "Field not found"})
				return
			}
			c.JSON(500, gin.H{"error": "Failed to retrieve field select options"})
			return
		}

		logger.Info().
			Int64("field_id", fieldID).
			Int("count", len(options)).
			Msg("[GetFieldSelectOptionsHandler] Field select options retrieved successfully")

		c.JSON(200, gin.H{
			"data":  options,
			"count": len(options),
			"_links": map[string]map[string]string{
				"field": {
					"href": fmt.Sprintf("/api/fields/%d", fieldID),
				},
			},
		})
	}
}

func GetFieldRelationRecordsHandler(fieldService *service.FieldService, logger zerolog.Logger) gin.HandlerFunc {
	service := *fieldService

	return func(c *gin.Context) {
		logger.Debug().Msg("[GetFieldRelationRecordsHandler] Getting field relation records")

		id := c.Param("id")
		fieldID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Error().Err(err).Str("id", id).Msg("[GetFieldRelationRecordsHandler] Failed to parse field ID")
			c.JSON(400, gin.H{"error": "Invalid field ID format"})
			return
		}

		records, err := service.GetFieldRelationRecords(fieldID, logger)
		if err != nil {
			logger.Error().Err(err).Int64("field_id", fieldID).Msg("[GetFieldRelationRecordsHandler] Failed to get field relation records")
			if err.Error() == fmt.Sprintf("field with ID %d not found", fieldID) {
				c.JSON(404, gin.H{"error": "Field not found"})
				return
			}
			c.JSON(500, gin.H{"error": "Failed to retrieve field relation records"})
			return
		}

		logger.Info().
			Int64("field_id", fieldID).
			Int("count", len(records)).
			Msg("[GetFieldRelationRecordsHandler] Field relation records retrieved successfully")

		c.JSON(200, gin.H{
			"data":  records,
			"count": len(records),
			"_links": map[string]map[string]string{
				"field": {
					"href": fmt.Sprintf("/api/fields/%d", fieldID),
				},
			},
		})
	}
}

func GetFieldTypesHandler(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Debug().Msg("[GetFieldTypesHandler] Getting available field types")

		fieldTypes := []map[string]interface{}{
			{"value": domain.FieldTypeText, "label": "Text", "description": "Single line text"},
			{"value": domain.FieldTypeNumber, "label": "Number", "description": "Numeric values"},
			{"value": domain.FieldTypeCheckbox, "label": "Checkbox", "description": "Boolean true/false"},
			{"value": domain.FieldTypeDate, "label": "Date", "description": "Date values"},
			{"value": domain.FieldTypeSingleSelect, "label": "Single Select", "description": "Select one option"},
			{"value": domain.FieldTypeMultiSelect, "label": "Multi Select", "description": "Select multiple options"},
			{"value": domain.FieldTypeMany2One, "label": "Many to One", "description": "Link to one record"},
			{"value": domain.FieldTypeOne2Many, "label": "One to Many", "description": "Link to multiple records"},
			{"value": domain.FieldTypeMany2Many, "label": "Many to Many", "description": "Bidirectional links"},
			{"value": domain.FieldTypeFormula, "label": "Formula", "description": "Calculated field"},
		}

		logger.Info().Int("count", len(fieldTypes)).Msg("[GetFieldTypesHandler] Field types retrieved")
		c.JSON(200, gin.H{
			"data":  fieldTypes,
			"count": len(fieldTypes),
		})
	}
}
