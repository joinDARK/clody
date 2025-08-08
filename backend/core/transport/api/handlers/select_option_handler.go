package handlers

import (
	"clody/core/service"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type CreateSelectOptionInput struct {
	FieldID  int64  `json:"field_id" binding:"required"`
	Label    string `json:"label" binding:"required,min=1,max=255"`
	Color    string `json:"color" binding:"required,min=1,max=255"`
	Position int    `json:"position" binding:"required,min=0"`
}

type UpdateSelectOptionInput struct {
	Label    string `json:"label" binding:"omitempty,min=1,max=255"`
	Color    string `json:"color" binding:"omitempty,min=1,max=255"`
	Position *int   `json:"position" binding:"omitempty,min=0"`
}

func GetSelectOptionByIDHandler(selectOptionService *service.SelectOptionService, logger zerolog.Logger) gin.HandlerFunc {
	service := *selectOptionService

	return func(c *gin.Context) {
		logger.Debug().Msg("[GetSelectOptionByIDHandler] Getting select option by ID")

		id := c.Param("id")
		selectOptionID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Error().Err(err).Str("id", id).Msg("[GetSelectOptionByIDHandler] Failed to parse select option ID")
			c.JSON(400, gin.H{"error": "Invalid select option ID format"})
			return
		}

		selectOption, err := service.GetSelectOptionByID(selectOptionID, logger)
		if err != nil {
			logger.Error().Err(err).Int64("select_option_id", selectOptionID).Msg("[GetSelectOptionByIDHandler] Failed to get select option")
			if err.Error() == fmt.Sprintf("select option with ID %d not found", selectOptionID) {
				c.JSON(404, gin.H{"error": "Select option not found"})
				return
			}
			c.JSON(500, gin.H{"error": "Failed to retrieve select option"})
			return
		}

		logger.Info().Int64("select_option_id", selectOptionID).Msg("[GetSelectOptionByIDHandler] Select option retrieved successfully")
		c.JSON(200, gin.H{
			"data": selectOption,
			"_links": map[string]map[string]string{
				"self": {
					"href": fmt.Sprintf("/api/select-options/%d", selectOption.ID),
				},
				"field": {
					"href": fmt.Sprintf("/api/fields/%d", selectOption.FieldID),
				},
				"update": {
					"href":   fmt.Sprintf("/api/select-options/%d", selectOption.ID),
					"method": "PUT",
				},
				"delete": {
					"href":   fmt.Sprintf("/api/select-options/%d", selectOption.ID),
					"method": "DELETE",
				},
			},
		})
	}
}

func CreateSelectOptionHandler(selectOptionService *service.SelectOptionService, logger zerolog.Logger) gin.HandlerFunc {
	service := *selectOptionService

	return func(c *gin.Context) {
		logger.Debug().Msg("[CreateSelectOptionHandler] Creating new select option")

		var input CreateSelectOptionInput
		if err := c.ShouldBindJSON(&input); err != nil {
			logger.Error().Err(err).Msg("[CreateSelectOptionHandler] Failed to bind JSON input")
			c.JSON(400, gin.H{"error": "Invalid input data", "details": err.Error()})
			return
		}

		selectOption, err := service.CreateSelectOption(
			input.FieldID,
			input.Label,
			input.Color,
			input.Position,
			logger,
		)
		if err != nil {
			logger.Error().Err(err).
				Str("label", input.Label).
				Int64("field_id", input.FieldID).
				Msg("[CreateSelectOptionHandler] Failed to create select option")
			c.JSON(500, gin.H{"error": "Failed to create select option"})
			return
		}

		logger.Info().
			Int64("select_option_id", selectOption.ID).
			Str("label", selectOption.Label).
			Str("color", selectOption.Color).
			Msg("[CreateSelectOptionHandler] Select option created successfully")

		c.JSON(201, gin.H{
			"data": selectOption,
			"_links": map[string]map[string]string{
				"self": {
					"href": fmt.Sprintf("/api/select-options/%d", selectOption.ID),
				},
				"field": {
					"href": fmt.Sprintf("/api/fields/%d", selectOption.FieldID),
				},
			},
		})
	}
}

func UpdateSelectOptionHandler(selectOptionService *service.SelectOptionService, logger zerolog.Logger) gin.HandlerFunc {
	service := *selectOptionService

	return func(c *gin.Context) {
		logger.Debug().Msg("[UpdateSelectOptionHandler] Updating select option")

		id := c.Param("id")
		selectOptionID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Error().Err(err).Str("id", id).Msg("[UpdateSelectOptionHandler] Failed to parse select option ID")
			c.JSON(400, gin.H{"error": "Invalid select option ID format"})
			return
		}

		var input UpdateSelectOptionInput
		if err := c.ShouldBindJSON(&input); err != nil {
			logger.Error().Err(err).Msg("[UpdateSelectOptionHandler] Failed to bind JSON input")
			c.JSON(400, gin.H{"error": "Invalid input data", "details": err.Error()})
			return
		}

		selectOption, err := service.UpdateSelectOption(
			selectOptionID,
			input.Label,
			input.Color,
			input.Position,
			logger,
		)
		if err != nil {
			logger.Error().Err(err).Int64("select_option_id", selectOptionID).Msg("[UpdateSelectOptionHandler] Failed to update select option")
			if err.Error() == fmt.Sprintf("select option with ID %d not found", selectOptionID) {
				c.JSON(404, gin.H{"error": "Select option not found"})
				return
			}
			c.JSON(500, gin.H{"error": "Failed to update select option"})
			return
		}

		logger.Info().
			Int64("select_option_id", selectOption.ID).
			Str("label", selectOption.Label).
			Msg("[UpdateSelectOptionHandler] Select option updated successfully")

		c.JSON(200, gin.H{
			"data": selectOption,
			"_links": map[string]map[string]string{
				"self": {
					"href": fmt.Sprintf("/api/select-options/%d", selectOption.ID),
				},
				"field": {
					"href": fmt.Sprintf("/api/fields/%d", selectOption.FieldID),
				},
			},
		})
	}
}

func DeleteSelectOptionHandler(selectOptionService *service.SelectOptionService, logger zerolog.Logger) gin.HandlerFunc {
	service := *selectOptionService

	return func(c *gin.Context) {
		logger.Debug().Msg("[DeleteSelectOptionHandler] Deleting select option")

		id := c.Param("id")
		selectOptionID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Error().Err(err).Str("id", id).Msg("[DeleteSelectOptionHandler] Failed to parse select option ID")
			c.JSON(400, gin.H{"error": "Invalid select option ID format"})
			return
		}

		err = service.DeleteSelectOption(selectOptionID, logger)
		if err != nil {
			logger.Error().Err(err).Int64("select_option_id", selectOptionID).Msg("[DeleteSelectOptionHandler] Failed to delete select option")
			if err.Error() == fmt.Sprintf("select option with ID %d not found", selectOptionID) {
				c.JSON(404, gin.H{"error": "Select option not found"})
				return
			}
			c.JSON(500, gin.H{"error": "Failed to delete select option"})
			return
		}

		logger.Info().Int64("select_option_id", selectOptionID).Msg("[DeleteSelectOptionHandler] Select option deleted successfully")
		c.JSON(204, gin.H{"message": "Select option deleted successfully"})
	}
}

func GetSelectOptionsByFieldIDHandler(selectOptionService *service.SelectOptionService, logger zerolog.Logger) gin.HandlerFunc {
	service := *selectOptionService

	return func(c *gin.Context) {
		logger.Debug().Msg("[GetSelectOptionsByFieldIDHandler] Getting select options by field ID")

		fieldID := c.Param("field_id")
		fieldIDInt, err := strconv.ParseInt(fieldID, 10, 64)
		if err != nil {
			logger.Error().Err(err).Str("field_id", fieldID).Msg("[GetSelectOptionsByFieldIDHandler] Failed to parse field ID")
			c.JSON(400, gin.H{"error": "Invalid field ID format"})
			return
		}

		selectOptions, err := service.GetSelectOptionsByFieldID(fieldIDInt, logger)
		if err != nil {
			logger.Error().Err(err).Int64("field_id", fieldIDInt).Msg("[GetSelectOptionsByFieldIDHandler] Failed to get select options")
			c.JSON(500, gin.H{"error": "Failed to retrieve select options"})
			return
		}

		logger.Info().
			Int64("field_id", fieldIDInt).
			Int("count", len(selectOptions)).
			Msg("[GetSelectOptionsByFieldIDHandler] Select options retrieved successfully")

		c.JSON(200, gin.H{
			"data":  selectOptions,
			"count": len(selectOptions),
			"_links": map[string]map[string]string{
				"field": {
					"href": fmt.Sprintf("/api/fields/%d", fieldIDInt),
				},
			},
		})
	}
}
