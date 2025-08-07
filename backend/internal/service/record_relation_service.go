package service

import (
	"dater/backend/internal/domain"
	"dater/backend/internal/repo"
	"errors"
	"fmt"

	"github.com/rs/zerolog"
)

type RecordRelationService struct {
	repo repo.RecordRelationRepo
}

func NewRecordRelationService(repo repo.RecordRelationRepo) *RecordRelationService {
	return &RecordRelationService{
		repo: repo,
	}
}

// GetRecordRelationByID возвращает связь записи по её ID
func (s *RecordRelationService) GetRecordRelationByID(id int64, logger zerolog.Logger) (*domain.RecordRelation, error) {
	logger.Debug().Int64("record_relation_id", id).Msg("[GetRecordRelationByID] Getting record relation by ID")

	if id <= 0 {
		logger.Error().Int64("record_relation_id", id).Msg("[GetRecordRelationByID] Invalid record relation ID")
		return nil, errors.New("record relation ID must be greater than 0")
	}

	recordRelation, err := s.repo.GetByID(id)
	if err != nil {
		logger.Error().Err(err).Int64("record_relation_id", id).Msg("[GetRecordRelationByID] Failed to get record relation")
		return nil, fmt.Errorf("failed to get record relation with ID %d: %w", id, err)
	}

	if recordRelation == nil {
		logger.Warn().Int64("record_relation_id", id).Msg("[GetRecordRelationByID] Record relation not found")
		return nil, fmt.Errorf("record relation with ID %d not found", id)
	}

	logger.Debug().
		Int64("record_relation_id", recordRelation.ID).
		Int64("field_id", recordRelation.FieldID).
		Int64("source_record_id", recordRelation.SourceRecordID).
		Int64("target_record_id", recordRelation.TargetRecordID).
		Msg("[GetRecordRelationByID] Record relation retrieved successfully")

	return recordRelation, nil
}

// CreateRecordRelation создает новую связь записи
func (s *RecordRelationService) CreateRecordRelation(fieldID, sourceRecordID, targetRecordID int64, logger zerolog.Logger) (*domain.RecordRelation, error) {
	logger.Debug().
		Int64("field_id", fieldID).
		Int64("source_record_id", sourceRecordID).
		Int64("target_record_id", targetRecordID).
		Msg("[CreateRecordRelation] Creating new record relation")

	// Validate input parameters
	if err := s.validateRecordRelationInput(fieldID, sourceRecordID, targetRecordID, logger); err != nil {
		return nil, err
	}

	// Check if relation already exists
	if err := s.validateRelationUniqueness(fieldID, sourceRecordID, targetRecordID, logger); err != nil {
		return nil, err
	}

	// Create record relation object
	recordRelation := &domain.RecordRelation{
		FieldID:        fieldID,
		SourceRecordID: sourceRecordID,
		TargetRecordID: targetRecordID,
	}

	// Save to database
	err := s.repo.Create(recordRelation)
	if err != nil {
		logger.Error().Err(err).
			Int64("field_id", fieldID).
			Int64("source_record_id", sourceRecordID).
			Int64("target_record_id", targetRecordID).
			Msg("[CreateRecordRelation] Failed to create record relation")
		return nil, fmt.Errorf("failed to create record relation: %w", err)
	}

	logger.Info().
		Int64("record_relation_id", recordRelation.ID).
		Int64("field_id", recordRelation.FieldID).
		Int64("source_record_id", recordRelation.SourceRecordID).
		Int64("target_record_id", recordRelation.TargetRecordID).
		Msg("[CreateRecordRelation] Record relation created successfully")

	return recordRelation, nil
}

// UpdateRecordRelation обновляет существующую связь записи
func (s *RecordRelationService) UpdateRecordRelation(id int64, fieldID, sourceRecordID, targetRecordID *int64, logger zerolog.Logger) (*domain.RecordRelation, error) {
	logger.Debug().
		Int64("record_relation_id", id).
		Msg("[UpdateRecordRelation] Updating record relation")

	// Validate record relation ID
	if id <= 0 {
		logger.Error().Int64("record_relation_id", id).Msg("[UpdateRecordRelation] Invalid record relation ID")
		return nil, errors.New("record relation ID must be greater than 0")
	}

	// Get existing record relation
	existingRecordRelation, err := s.repo.GetByID(id)
	if err != nil {
		logger.Error().Err(err).Int64("record_relation_id", id).Msg("[UpdateRecordRelation] Failed to get existing record relation")
		return nil, fmt.Errorf("failed to get record relation with ID %d: %w", id, err)
	}

	if existingRecordRelation == nil {
		logger.Warn().Int64("record_relation_id", id).Msg("[UpdateRecordRelation] Record relation not found")
		return nil, fmt.Errorf("record relation with ID %d not found", id)
	}

	// Update fields if provided
	if fieldID != nil {
		if *fieldID <= 0 {
			logger.Error().Int64("field_id", *fieldID).Msg("[UpdateRecordRelation] Invalid field ID")
			return nil, errors.New("field ID must be greater than 0")
		}
		existingRecordRelation.FieldID = *fieldID
	}

	if sourceRecordID != nil {
		if *sourceRecordID <= 0 {
			logger.Error().Int64("source_record_id", *sourceRecordID).Msg("[UpdateRecordRelation] Invalid source record ID")
			return nil, errors.New("source record ID must be greater than 0")
		}
		existingRecordRelation.SourceRecordID = *sourceRecordID
	}

	if targetRecordID != nil {
		if *targetRecordID <= 0 {
			logger.Error().Int64("target_record_id", *targetRecordID).Msg("[UpdateRecordRelation] Invalid target record ID")
			return nil, errors.New("target record ID must be greater than 0")
		}
		existingRecordRelation.TargetRecordID = *targetRecordID
	}

	// Prevent self-relation
	if existingRecordRelation.SourceRecordID == existingRecordRelation.TargetRecordID {
		logger.Error().
			Int64("source_record_id", existingRecordRelation.SourceRecordID).
			Int64("target_record_id", existingRecordRelation.TargetRecordID).
			Msg("[UpdateRecordRelation] Self-relation not allowed")
		return nil, errors.New("record cannot be related to itself")
	}

	// Save changes
	err = s.repo.Update(existingRecordRelation)
	if err != nil {
		logger.Error().Err(err).
			Int64("record_relation_id", id).
			Msg("[UpdateRecordRelation] Failed to update record relation")
		return nil, fmt.Errorf("failed to update record relation with ID %d: %w", id, err)
	}

	logger.Info().
		Int64("record_relation_id", existingRecordRelation.ID).
		Int64("field_id", existingRecordRelation.FieldID).
		Int64("source_record_id", existingRecordRelation.SourceRecordID).
		Int64("target_record_id", existingRecordRelation.TargetRecordID).
		Msg("[UpdateRecordRelation] Record relation updated successfully")

	return existingRecordRelation, nil
}

// DeleteRecordRelation удаляет связь записи по ID
func (s *RecordRelationService) DeleteRecordRelation(id int64, logger zerolog.Logger) error {
	logger.Debug().Int64("record_relation_id", id).Msg("[DeleteRecordRelation] Deleting record relation")

	if id <= 0 {
		logger.Error().Int64("record_relation_id", id).Msg("[DeleteRecordRelation] Invalid record relation ID")
		return errors.New("record relation ID must be greater than 0")
	}

	// Check if record relation exists before deletion
	existingRecordRelation, err := s.repo.GetByID(id)
	if err != nil {
		logger.Error().Err(err).Int64("record_relation_id", id).Msg("[DeleteRecordRelation] Failed to get record relation for deletion")
		return fmt.Errorf("failed to verify record relation existence with ID %d: %w", id, err)
	}

	if existingRecordRelation == nil {
		logger.Warn().Int64("record_relation_id", id).Msg("[DeleteRecordRelation] Record relation not found")
		return fmt.Errorf("record relation with ID %d not found", id)
	}

	// Delete the record relation
	err = s.repo.Delete(id)
	if err != nil {
		logger.Error().Err(err).Int64("record_relation_id", id).Msg("[DeleteRecordRelation] Failed to delete record relation")
		return fmt.Errorf("failed to delete record relation with ID %d: %w", id, err)
	}

	logger.Info().Int64("record_relation_id", id).Msg("[DeleteRecordRelation] Record relation deleted successfully")
	return nil
}

// GetRecordRelationsByFieldID возвращает все связи записей для конкретного поля
func (s *RecordRelationService) GetRecordRelationsByFieldID(fieldID int64, logger zerolog.Logger) ([]*domain.RecordRelation, error) {
	logger.Debug().Int64("field_id", fieldID).Msg("[GetRecordRelationsByFieldID] Getting record relations by field ID")

	if fieldID <= 0 {
		logger.Error().Int64("field_id", fieldID).Msg("[GetRecordRelationsByFieldID] Invalid field ID")
		return nil, errors.New("field ID must be greater than 0")
	}

	recordRelations, err := s.repo.GetByFieldID(fieldID)
	if err != nil {
		logger.Error().Err(err).Int64("field_id", fieldID).Msg("[GetRecordRelationsByFieldID] Failed to get record relations")
		return nil, fmt.Errorf("failed to get record relations for field %d: %w", fieldID, err)
	}

	if recordRelations == nil {
		recordRelations = []*domain.RecordRelation{}
	}

	logger.Debug().
		Int64("field_id", fieldID).
		Int("count", len(recordRelations)).
		Msg("[GetRecordRelationsByFieldID] Record relations retrieved successfully")

	return recordRelations, nil
}

// GetRecordRelationsBySourceRecordID возвращает все связи записей по ID исходной записи
func (s *RecordRelationService) GetRecordRelationsBySourceRecordID(recordID int64, logger zerolog.Logger) ([]*domain.RecordRelation, error) {
	logger.Debug().Int64("source_record_id", recordID).Msg("[GetRecordRelationsBySourceRecordID] Getting record relations by source record ID")

	if recordID <= 0 {
		logger.Error().Int64("source_record_id", recordID).Msg("[GetRecordRelationsBySourceRecordID] Invalid record ID")
		return nil, errors.New("record ID must be greater than 0")
	}

	recordRelations, err := s.repo.GetBySourceRecordID(recordID)
	if err != nil {
		logger.Error().Err(err).Int64("source_record_id", recordID).Msg("[GetRecordRelationsBySourceRecordID] Failed to get record relations")
		return nil, fmt.Errorf("failed to get record relations by source record ID %d: %w", recordID, err)
	}

	if recordRelations == nil {
		recordRelations = []*domain.RecordRelation{}
	}

	logger.Debug().
		Int64("source_record_id", recordID).
		Int("count", len(recordRelations)).
		Msg("[GetRecordRelationsBySourceRecordID] Record relations retrieved successfully")

	return recordRelations, nil
}

// GetRecordRelationsByTargetRecordID возвращает все связи записей по ID целевой записи
func (s *RecordRelationService) GetRecordRelationsByTargetRecordID(recordID int64, logger zerolog.Logger) ([]*domain.RecordRelation, error) {
	logger.Debug().Int64("target_record_id", recordID).Msg("[GetRecordRelationsByTargetRecordID] Getting record relations by target record ID")

	if recordID <= 0 {
		logger.Error().Int64("target_record_id", recordID).Msg("[GetRecordRelationsByTargetRecordID] Invalid record ID")
		return nil, errors.New("record ID must be greater than 0")
	}

	recordRelations, err := s.repo.GetByTargetRecordID(recordID)
	if err != nil {
		logger.Error().Err(err).Int64("target_record_id", recordID).Msg("[GetRecordRelationsByTargetRecordID] Failed to get record relations")
		return nil, fmt.Errorf("failed to get record relations by target record ID %d: %w", recordID, err)
	}

	if recordRelations == nil {
		recordRelations = []*domain.RecordRelation{}
	}

	logger.Debug().
		Int64("target_record_id", recordID).
		Int("count", len(recordRelations)).
		Msg("[GetRecordRelationsByTargetRecordID] Record relations retrieved successfully")

	return recordRelations, nil
}

// validateRecordRelationExists проверяет, существует ли связь записи
func (s *RecordRelationService) validateRecordRelationExists(id int64, logger zerolog.Logger) error {
	logger.Debug().Int64("record_relation_id", id).Msg("[validateRecordRelationExists] Validating record relation existence")

	if id <= 0 {
		logger.Error().Int64("record_relation_id", id).Msg("[validateRecordRelationExists] Invalid record relation ID")
		return errors.New("record relation ID must be greater than 0")
	}

	recordRelation, err := s.repo.GetByID(id)
	if err != nil {
		logger.Error().Err(err).Int64("record_relation_id", id).Msg("[validateRecordRelationExists] Failed to get record relation")
		return fmt.Errorf("failed to validate record relation existence with ID %d: %w", id, err)
	}

	if recordRelation == nil {
		logger.Warn().Int64("record_relation_id", id).Msg("[validateRecordRelationExists] Record relation not found")
		return fmt.Errorf("record relation with ID %d does not exist", id)
	}

	logger.Debug().Int64("record_relation_id", id).Msg("[validateRecordRelationExists] Record relation exists")
	return nil
}

// validateRecordRelationInput валидирует входные параметры связи записи
func (s *RecordRelationService) validateRecordRelationInput(fieldID, sourceRecordID, targetRecordID int64, logger zerolog.Logger) error {
	if fieldID <= 0 {
		logger.Error().Int64("field_id", fieldID).Msg("[validateRecordRelationInput] Invalid field ID")
		return errors.New("field ID must be greater than 0")
	}

	if sourceRecordID <= 0 {
		logger.Error().Int64("source_record_id", sourceRecordID).Msg("[validateRecordRelationInput] Invalid source record ID")
		return errors.New("source record ID must be greater than 0")
	}

	if targetRecordID <= 0 {
		logger.Error().Int64("target_record_id", targetRecordID).Msg("[validateRecordRelationInput] Invalid target record ID")
		return errors.New("target record ID must be greater than 0")
	}

	if sourceRecordID == targetRecordID {
		logger.Error().
			Int64("source_record_id", sourceRecordID).
			Int64("target_record_id", targetRecordID).
			Msg("[validateRecordRelationInput] Self-relation not allowed")
		return errors.New("record cannot be related to itself")
	}

	return nil
}

// validateRelationUniqueness проверяет уникальность связи
func (s *RecordRelationService) validateRelationUniqueness(fieldID, sourceRecordID, targetRecordID int64, logger zerolog.Logger) error {
	logger.Debug().
		Int64("field_id", fieldID).
		Int64("source_record_id", sourceRecordID).
		Int64("target_record_id", targetRecordID).
		Msg("[validateRelationUniqueness] Validating relation uniqueness")

	// Get existing relations by field ID
	existingRelations, err := s.repo.GetByFieldID(fieldID)
	if err != nil {
		logger.Error().Err(err).Int64("field_id", fieldID).Msg("[validateRelationUniqueness] Failed to get existing relations")
		return fmt.Errorf("failed to check relation uniqueness: %w", err)
	}

	// Check for duplicate relations
	for _, relation := range existingRelations {
		if relation.SourceRecordID == sourceRecordID && relation.TargetRecordID == targetRecordID {
			logger.Warn().
				Int64("field_id", fieldID).
				Int64("source_record_id", sourceRecordID).
				Int64("target_record_id", targetRecordID).
				Msg("[validateRelationUniqueness] Relation already exists")
			return fmt.Errorf("relation between records %d and %d already exists in field %d", sourceRecordID, targetRecordID, fieldID)
		}
	}

	logger.Debug().Msg("[validateRelationUniqueness] Relation is unique")
	return nil
}
