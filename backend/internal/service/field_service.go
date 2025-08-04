package service

import (
	"dater/backend/internal/domain"
	"dater/backend/internal/repo"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/datatypes"
)

type FieldService struct {
	repo repo.FieldRepo
}

func NewFieldService(repo repo.FieldRepo) *FieldService {
	return &FieldService{
		repo: repo,
	}
}

// GetFieldByID возвращает поле по его ID
func (s *FieldService) GetFieldByID(id int64, logger zerolog.Logger) (*domain.Field, error) {
	logger.Debug().Int64("field_id", id).Msg("[GetFieldByID] Getting field by ID")

	if id <= 0 {
		logger.Error().Int64("field_id", id).Msg("[GetFieldByID] Invalid field ID")
		return nil, errors.New("field ID must be greater than 0")
	}

	field, err := s.repo.GetByID(id)
	if err != nil {
		logger.Error().Err(err).Int64("field_id", id).Msg("[GetFieldByID] Failed to get field")
		return nil, fmt.Errorf("failed to get field with ID %d: %w", id, err)
	}

	if field == nil {
		logger.Warn().Int64("field_id", id).Msg("[GetFieldByID] Field not found")
		return nil, fmt.Errorf("field with ID %d not found", id)
	}

	logger.Debug().
		Int64("field_id", field.ID).
		Str("name", field.Name).
		Str("type", string(field.Type)).
		Int64("table_id", field.TableID).
		Msg("[GetFieldByID] Field retrieved successfully")

	return field, nil
}

// CreateField создает новое поле
func (s *FieldService) CreateField(tableID int64, name string, fieldType domain.FieldType, position int, description *string, options datatypes.JSON, logger zerolog.Logger) (*domain.Field, error) {
	logger.Debug().
		Int64("table_id", tableID).
		Str("name", name).
		Str("type", string(fieldType)).
		Int("position", position).
		Msg("[CreateField] Creating new field")

	// Validate input parameters
	if err := s.validateFieldInput(tableID, name, fieldType, position, logger); err != nil {
		return nil, err
	}

	// Create field object
	field := &domain.Field{
		TableID:     tableID,
		Name:        name,
		Type:        fieldType,
		Position:    position,
		Description: description,
		Options:     options,
		CreatedAt:   time.Now(),
	}

	// Save to database
	err := s.repo.Create(field)
	if err != nil {
		logger.Error().Err(err).
			Str("name", name).
			Int64("table_id", tableID).
			Msg("[CreateField] Failed to create field")
		return nil, fmt.Errorf("failed to create field '%s': %w", name, err)
	}

	logger.Info().
		Int64("field_id", field.ID).
		Str("name", field.Name).
		Str("type", string(field.Type)).
		Int64("table_id", field.TableID).
		Msg("[CreateField] Field created successfully")

	return field, nil
}

// UpdateField обновляет существующее поле
func (s *FieldService) UpdateField(id int64, name string, fieldType domain.FieldType, position int, description *string, options datatypes.JSON, logger zerolog.Logger) (*domain.Field, error) {
	logger.Debug().
		Int64("field_id", id).
		Str("name", name).
		Str("type", string(fieldType)).
		Int("position", position).
		Msg("[UpdateField] Updating field")

	// Validate field ID
	if id <= 0 {
		logger.Error().Int64("field_id", id).Msg("[UpdateField] Invalid field ID")
		return nil, errors.New("field ID must be greater than 0")
	}

	// Get existing field
	existingField, err := s.repo.GetByID(id)
	if err != nil {
		logger.Error().Err(err).Int64("field_id", id).Msg("[UpdateField] Failed to get existing field")
		return nil, fmt.Errorf("failed to get field with ID %d: %w", id, err)
	}

	if existingField == nil {
		logger.Warn().Int64("field_id", id).Msg("[UpdateField] Field not found")
		return nil, fmt.Errorf("field with ID %d not found", id)
	}

	// Validate updated field data if provided
	if name != "" {
		if err := s.validateFieldName(name, logger); err != nil {
			return nil, err
		}
		existingField.Name = name
	}

	if fieldType != "" {
		if !fieldType.IsValid() {
			logger.Error().Str("type", string(fieldType)).Msg("[UpdateField] Invalid field type")
			return nil, fmt.Errorf("invalid field type: %s", fieldType)
		}
		existingField.Type = fieldType
	}

	if position > 0 {
		existingField.Position = position
	}

	existingField.Description = description
	existingField.Options = options

	// Save changes
	err = s.repo.Update(existingField)
	if err != nil {
		logger.Error().Err(err).
			Int64("field_id", id).
			Str("name", name).
			Msg("[UpdateField] Failed to update field")
		return nil, fmt.Errorf("failed to update field with ID %d: %w", id, err)
	}

	logger.Info().
		Int64("field_id", existingField.ID).
		Str("name", existingField.Name).
		Str("type", string(existingField.Type)).
		Msg("[UpdateField] Field updated successfully")

	return existingField, nil
}

// DeleteField удаляет поле по ID
func (s *FieldService) DeleteField(id int64, logger zerolog.Logger) (int64, error) {
	logger.Debug().Int64("field_id", id).Msg("[DeleteField] Deleting field")

	if id <= 0 {
		logger.Error().Int64("field_id", id).Msg("[DeleteField] Invalid field ID")
		return 0, errors.New("field ID must be greater than 0")
	}

	// Check if field exists before deletion
	existingField, err := s.repo.GetByID(id)
	if err != nil {
		logger.Error().Err(err).Int64("field_id", id).Msg("[DeleteField] Failed to get field for deletion")
		return 0, fmt.Errorf("failed to verify field existence with ID %d: %w", id, err)
	}

	if existingField == nil {
		logger.Warn().Int64("field_id", id).Msg("[DeleteField] Field not found")
		return 0, fmt.Errorf("field with ID %d not found", id)
	}

	// Delete the field
	id, err = s.repo.Delete(id)
	if err != nil {
		logger.Error().Err(err).Int64("field_id", id).Msg("[DeleteField] Failed to delete field")
		return 0, fmt.Errorf("failed to delete field with ID %d: %w", id, err)
	}

	logger.Info().Int64("field_id", id).Msg("[DeleteField] Field deleted successfully")
	return id, nil
}

// GetFieldCells возвращает все ячейки в поле по ID поля
func (s *FieldService) GetFieldCells(fieldID int64, logger zerolog.Logger) ([]*domain.Cell, error) {
	logger.Debug().Int64("field_id", fieldID).Msg("[GetFieldCells] Getting field cells")

	if fieldID <= 0 {
		logger.Error().Int64("field_id", fieldID).Msg("[GetFieldCells] Invalid field ID")
		return nil, errors.New("field ID must be greater than 0")
	}

	// Verify field exists
	if err := s.validateFieldExists(fieldID, logger); err != nil {
		return nil, err
	}

	cells, err := s.repo.GetAllCells(fieldID)
	if err != nil {
		logger.Error().Err(err).Int64("field_id", fieldID).Msg("[GetFieldCells] Failed to get cells")
		return nil, fmt.Errorf("failed to get cells for field %d: %w", fieldID, err)
	}

	if cells == nil {
		cells = []*domain.Cell{}
	}

	logger.Debug().
		Int64("field_id", fieldID).
		Int("count", len(cells)).
		Msg("[GetFieldCells] Cells retrieved successfully")

	return cells, nil
}

// GetFieldSelectOptions возвращает все опции выбора для конкретного поля
func (s *FieldService) GetFieldSelectOptions(fieldID int64, logger zerolog.Logger) ([]*domain.SelectOption, error) {
	logger.Debug().Int64("field_id", fieldID).Msg("[GetFieldSelectOptions] Getting field select options")

	if fieldID <= 0 {
		logger.Error().Int64("field_id", fieldID).Msg("[GetFieldSelectOptions] Invalid field ID")
		return nil, errors.New("field ID must be greater than 0")
	}

	// Verify field exists
	field, err := s.repo.GetByID(fieldID)
	if err != nil {
		logger.Error().Err(err).Int64("field_id", fieldID).Msg("[GetFieldSelectOptions] Failed to verify field existence")
		return nil, fmt.Errorf("failed to verify field existence with ID %d: %w", fieldID, err)
	}

	if field == nil {
		logger.Warn().Int64("field_id", fieldID).Msg("[GetFieldSelectOptions] Field not found")
		return nil, fmt.Errorf("field with ID %d not found", fieldID)
	}

	// Check if field type supports select options
	if field.Type != domain.FieldTypeSingleSelect && field.Type != domain.FieldTypeMultiSelect {
		logger.Warn().
			Int64("field_id", fieldID).
			Str("type", string(field.Type)).
			Msg("[GetFieldSelectOptions] Field type does not support select options")
		return []*domain.SelectOption{}, nil
	}

	options, err := s.repo.GetAllSelectOptions(fieldID)
	if err != nil {
		logger.Error().Err(err).Int64("field_id", fieldID).Msg("[GetFieldSelectOptions] Failed to get select options")
		return nil, fmt.Errorf("failed to get select options for field %d: %w", fieldID, err)
	}

	if options == nil {
		options = []*domain.SelectOption{}
	}

	logger.Debug().
		Int64("field_id", fieldID).
		Int("count", len(options)).
		Msg("[GetFieldSelectOptions] Select options retrieved successfully")

	return options, nil
}

// GetFieldRelationRecords возвращает все связанные записи для конкретного поля
func (s *FieldService) GetFieldRelationRecords(fieldID int64, logger zerolog.Logger) ([]*domain.Record, error) {
	logger.Debug().Int64("field_id", fieldID).Msg("[GetFieldRelationRecords] Getting field relation records")

	if fieldID <= 0 {
		logger.Error().Int64("field_id", fieldID).Msg("[GetFieldRelationRecords] Invalid field ID")
		return nil, errors.New("field ID must be greater than 0")
	}

	// Verify field exists and check if it's a relation field
	field, err := s.repo.GetByID(fieldID)
	if err != nil {
		logger.Error().Err(err).Int64("field_id", fieldID).Msg("[GetFieldRelationRecords] Failed to verify field existence")
		return nil, fmt.Errorf("failed to verify field existence with ID %d: %w", fieldID, err)
	}

	if field == nil {
		logger.Warn().Int64("field_id", fieldID).Msg("[GetFieldRelationRecords] Field not found")
		return nil, fmt.Errorf("field with ID %d not found", fieldID)
	}

	// Check if field type supports relations
	if !s.isRelationField(field.Type) {
		logger.Warn().
			Int64("field_id", fieldID).
			Str("type", string(field.Type)).
			Msg("[GetFieldRelationRecords] Field type does not support relations")
		return []*domain.Record{}, nil
	}

	records, err := s.repo.GetAllRelationsRecords(fieldID)
	if err != nil {
		logger.Error().Err(err).Int64("field_id", fieldID).Msg("[GetFieldRelationRecords] Failed to get relation records")
		return nil, fmt.Errorf("failed to get relation records for field %d: %w", fieldID, err)
	}

	if records == nil {
		records = []*domain.Record{}
	}

	logger.Debug().
		Int64("field_id", fieldID).
		Int("count", len(records)).
		Msg("[GetFieldRelationRecords] Relation records retrieved successfully")

	return records, nil
}

// validateFieldExists проверяет, существует ли поле
func (s *FieldService) validateFieldExists(id int64, logger zerolog.Logger) error {
	logger.Debug().Int64("field_id", id).Msg("[validateFieldExists] Validating field existence")

	if id <= 0 {
		logger.Error().Int64("field_id", id).Msg("[validateFieldExists] Invalid field ID")
		return errors.New("field ID must be greater than 0")
	}

	field, err := s.repo.GetByID(id)
	if err != nil {
		logger.Error().Err(err).Int64("field_id", id).Msg("[validateFieldExists] Failed to get field")
		return fmt.Errorf("failed to validate field existence with ID %d: %w", id, err)
	}

	if field == nil {
		logger.Warn().Int64("field_id", id).Msg("[validateFieldExists] Field not found")
		return fmt.Errorf("field with ID %d does not exist", id)
	}

	logger.Debug().Int64("field_id", id).Msg("[validateFieldExists] Field exists")
	return nil
}

// validateFieldInput валидирует входные параметры поля
func (s *FieldService) validateFieldInput(tableID int64, name string, fieldType domain.FieldType, position int, logger zerolog.Logger) error {
	if tableID <= 0 {
		logger.Error().Int64("table_id", tableID).Msg("[validateFieldInput] Invalid table ID")
		return errors.New("table ID must be greater than 0")
	}

	if err := s.validateFieldName(name, logger); err != nil {
		return err
	}

	if !fieldType.IsValid() {
		logger.Error().Str("type", string(fieldType)).Msg("[validateFieldInput] Invalid field type")
		return fmt.Errorf("invalid field type: %s", fieldType)
	}

	if position < 0 {
		logger.Error().Int("position", position).Msg("[validateFieldInput] Invalid position")
		return errors.New("field position cannot be negative")
	}

	return nil
}

// validateFieldName валидирует имя поля
func (s *FieldService) validateFieldName(name string, logger zerolog.Logger) error {
	if name == "" {
		logger.Error().Msg("[validateFieldName] Field name cannot be empty")
		return errors.New("field name cannot be empty")
	}

	if len(name) > 255 {
		logger.Error().Int("length", len(name)).Msg("[validateFieldName] Field name too long")
		return errors.New("field name cannot exceed 255 characters")
	}

	return nil
}

// isRelationField проверяет, является ли тип поля отношением
func (s *FieldService) isRelationField(fieldType domain.FieldType) bool {
	return fieldType == domain.FieldTypeMany2One ||
		fieldType == domain.FieldTypeOne2Many ||
		fieldType == domain.FieldTypeMany2Many
}
