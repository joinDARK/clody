package service

import (
	"clody/core/domain"
	"clody/core/repo"
	"errors"
	"fmt"

	"github.com/rs/zerolog"
)

type SelectOptionService struct {
	repo repo.SelectOptionRepo
}

func NewSelectOptionService(repo repo.SelectOptionRepo) *SelectOptionService {
	return &SelectOptionService{
		repo: repo,
	}
}

// GetSelectOptionByID возвращает опцию выбора по её ID
func (s *SelectOptionService) GetSelectOptionByID(id int64, logger zerolog.Logger) (*domain.SelectOption, error) {
	logger.Debug().Int64("select_option_id", id).Msg("[GetSelectOptionByID] Getting select option by ID")

	if id <= 0 {
		logger.Error().Int64("select_option_id", id).Msg("[GetSelectOptionByID] Invalid select option ID")
		return nil, errors.New("select option ID must be greater than 0")
	}

	selectOption, err := s.repo.GetByID(id)
	if err != nil {
		logger.Error().Err(err).Int64("select_option_id", id).Msg("[GetSelectOptionByID] Failed to get select option")
		return nil, fmt.Errorf("failed to get select option with ID %d: %w", id, err)
	}

	if selectOption == nil {
		logger.Warn().Int64("select_option_id", id).Msg("[GetSelectOptionByID] Select option not found")
		return nil, fmt.Errorf("select option with ID %d not found", id)
	}

	logger.Debug().
		Int64("select_option_id", selectOption.ID).
		Str("label", selectOption.Label).
		Str("color", selectOption.Color).
		Int64("field_id", selectOption.FieldID).
		Msg("[GetSelectOptionByID] Select option retrieved successfully")

	return selectOption, nil
}

// CreateSelectOption создает новую опцию выбора
func (s *SelectOptionService) CreateSelectOption(fieldID int64, label, color string, position int, logger zerolog.Logger) (*domain.SelectOption, error) {
	logger.Debug().
		Int64("field_id", fieldID).
		Str("label", label).
		Str("color", color).
		Int("position", position).
		Msg("[CreateSelectOption] Creating new select option")

	// Validate input parameters
	if err := s.validateSelectOptionInput(fieldID, label, color, position, logger); err != nil {
		return nil, err
	}

	// Create select option object
	selectOption := &domain.SelectOption{
		FieldID:  fieldID,
		Label:    label,
		Color:    color,
		Position: position,
	}

	// Save to database
	err := s.repo.Create(selectOption)
	if err != nil {
		logger.Error().Err(err).
			Str("label", label).
			Int64("field_id", fieldID).
			Msg("[CreateSelectOption] Failed to create select option")
		return nil, fmt.Errorf("failed to create select option '%s': %w", label, err)
	}

	logger.Info().
		Int64("select_option_id", selectOption.ID).
		Str("label", selectOption.Label).
		Str("color", selectOption.Color).
		Int64("field_id", selectOption.FieldID).
		Msg("[CreateSelectOption] Select option created successfully")

	return selectOption, nil
}

// UpdateSelectOption обновляет существующую опцию выбора
func (s *SelectOptionService) UpdateSelectOption(id int64, label, color string, position *int, logger zerolog.Logger) (*domain.SelectOption, error) {
	logger.Debug().
		Int64("select_option_id", id).
		Str("label", label).
		Str("color", color).
		Msg("[UpdateSelectOption] Updating select option")

	// Validate select option ID
	if id <= 0 {
		logger.Error().Int64("select_option_id", id).Msg("[UpdateSelectOption] Invalid select option ID")
		return nil, errors.New("select option ID must be greater than 0")
	}

	// Get existing select option
	existingSelectOption, err := s.repo.GetByID(id)
	if err != nil {
		logger.Error().Err(err).Int64("select_option_id", id).Msg("[UpdateSelectOption] Failed to get existing select option")
		return nil, fmt.Errorf("failed to get select option with ID %d: %w", id, err)
	}

	if existingSelectOption == nil {
		logger.Warn().Int64("select_option_id", id).Msg("[UpdateSelectOption] Select option not found")
		return nil, fmt.Errorf("select option with ID %d not found", id)
	}

	// Validate and update fields if provided
	if label != "" {
		if err := s.validateLabel(label, logger); err != nil {
			return nil, err
		}
		existingSelectOption.Label = label
	}

	if color != "" {
		if err := s.validateColor(color, logger); err != nil {
			return nil, err
		}
		existingSelectOption.Color = color
	}

	if position != nil {
		if *position < 0 {
			logger.Error().Int("position", *position).Msg("[UpdateSelectOption] Invalid position")
			return nil, errors.New("select option position cannot be negative")
		}
		existingSelectOption.Position = *position
	}

	// Save changes
	err = s.repo.Update(existingSelectOption)
	if err != nil {
		logger.Error().Err(err).
			Int64("select_option_id", id).
			Str("label", label).
			Msg("[UpdateSelectOption] Failed to update select option")
		return nil, fmt.Errorf("failed to update select option with ID %d: %w", id, err)
	}

	logger.Info().
		Int64("select_option_id", existingSelectOption.ID).
		Str("label", existingSelectOption.Label).
		Str("color", existingSelectOption.Color).
		Msg("[UpdateSelectOption] Select option updated successfully")

	return existingSelectOption, nil
}

// DeleteSelectOption удаляет опцию выбора по ID
func (s *SelectOptionService) DeleteSelectOption(id int64, logger zerolog.Logger) error {
	logger.Debug().Int64("select_option_id", id).Msg("[DeleteSelectOption] Deleting select option")

	if id <= 0 {
		logger.Error().Int64("select_option_id", id).Msg("[DeleteSelectOption] Invalid select option ID")
		return errors.New("select option ID must be greater than 0")
	}

	// Check if select option exists before deletion
	existingSelectOption, err := s.repo.GetByID(id)
	if err != nil {
		logger.Error().Err(err).Int64("select_option_id", id).Msg("[DeleteSelectOption] Failed to get select option for deletion")
		return fmt.Errorf("failed to verify select option existence with ID %d: %w", id, err)
	}

	if existingSelectOption == nil {
		logger.Warn().Int64("select_option_id", id).Msg("[DeleteSelectOption] Select option not found")
		return fmt.Errorf("select option with ID %d not found", id)
	}

	// Delete the select option
	err = s.repo.Delete(id)
	if err != nil {
		logger.Error().Err(err).Int64("select_option_id", id).Msg("[DeleteSelectOption] Failed to delete select option")
		return fmt.Errorf("failed to delete select option with ID %d: %w", id, err)
	}

	logger.Info().Int64("select_option_id", id).Msg("[DeleteSelectOption] Select option deleted successfully")
	return nil
}

// GetSelectOptionsByFieldID возвращает все опции выбора для конкретного поля
func (s *SelectOptionService) GetSelectOptionsByFieldID(fieldID int64, logger zerolog.Logger) ([]*domain.SelectOption, error) {
	logger.Debug().Int64("field_id", fieldID).Msg("[GetSelectOptionsByFieldID] Getting select options by field ID")

	if fieldID <= 0 {
		logger.Error().Int64("field_id", fieldID).Msg("[GetSelectOptionsByFieldID] Invalid field ID")
		return nil, errors.New("field ID must be greater than 0")
	}

	selectOptions, err := s.repo.GetAllByFieldID(fieldID)
	if err != nil {
		logger.Error().Err(err).Int64("field_id", fieldID).Msg("[GetSelectOptionsByFieldID] Failed to get select options")
		return nil, fmt.Errorf("failed to get select options for field %d: %w", fieldID, err)
	}

	if selectOptions == nil {
		selectOptions = []*domain.SelectOption{}
	}

	logger.Debug().
		Int64("field_id", fieldID).
		Int("count", len(selectOptions)).
		Msg("[GetSelectOptionsByFieldID] Select options retrieved successfully")

	return selectOptions, nil
}

// validateSelectOptionExists проверяет, существует ли опция выбора
func (s *SelectOptionService) validateSelectOptionExists(id int64, logger zerolog.Logger) error {
	logger.Debug().Int64("select_option_id", id).Msg("[validateSelectOptionExists] Validating select option existence")

	if id <= 0 {
		logger.Error().Int64("select_option_id", id).Msg("[validateSelectOptionExists] Invalid select option ID")
		return errors.New("select option ID must be greater than 0")
	}

	selectOption, err := s.repo.GetByID(id)
	if err != nil {
		logger.Error().Err(err).Int64("select_option_id", id).Msg("[validateSelectOptionExists] Failed to get select option")
		return fmt.Errorf("failed to validate select option existence with ID %d: %w", id, err)
	}

	if selectOption == nil {
		logger.Warn().Int64("select_option_id", id).Msg("[validateSelectOptionExists] Select option not found")
		return fmt.Errorf("select option with ID %d does not exist", id)
	}

	logger.Debug().Int64("select_option_id", id).Msg("[validateSelectOptionExists] Select option exists")
	return nil
}

// validateSelectOptionInput валидирует входные параметры опции выбора
func (s *SelectOptionService) validateSelectOptionInput(fieldID int64, label, color string, position int, logger zerolog.Logger) error {
	if fieldID <= 0 {
		logger.Error().Int64("field_id", fieldID).Msg("[validateSelectOptionInput] Invalid field ID")
		return errors.New("field ID must be greater than 0")
	}

	if err := s.validateLabel(label, logger); err != nil {
		return err
	}

	if err := s.validateColor(color, logger); err != nil {
		return err
	}

	if position < 0 {
		logger.Error().Int("position", position).Msg("[validateSelectOptionInput] Invalid position")
		return errors.New("select option position cannot be negative")
	}

	return nil
}

// validateLabel валидирует метку опции выбора
func (s *SelectOptionService) validateLabel(label string, logger zerolog.Logger) error {
	if label == "" {
		logger.Error().Msg("[validateLabel] Select option label cannot be empty")
		return errors.New("select option label cannot be empty")
	}

	if len(label) > 255 {
		logger.Error().Int("length", len(label)).Msg("[validateLabel] Select option label too long")
		return errors.New("select option label cannot exceed 255 characters")
	}

	return nil
}

// validateColor валидирует цвет опции выбора
func (s *SelectOptionService) validateColor(color string, logger zerolog.Logger) error {
	if color == "" {
		logger.Error().Msg("[validateColor] Select option color cannot be empty")
		return errors.New("select option color cannot be empty")
	}

	if len(color) > 255 {
		logger.Error().Int("length", len(color)).Msg("[validateColor] Select option color too long")
		return errors.New("select option color cannot exceed 255 characters")
	}

	// Basic hex color validation
	if len(color) > 0 && color[0] == '#' && (len(color) == 7 || len(color) == 4) {
		for i := 1; i < len(color); i++ {
			c := color[i]
			if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
				logger.Error().Str("color", color).Msg("[validateColor] Invalid hex color format")
				return errors.New("invalid hex color format")
			}
		}
	}

	return nil
}
