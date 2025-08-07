package service

import (
	"dater/backend/internal/domain"
	"dater/backend/internal/repo"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog"
)

type RecordService struct {
	repo repo.RecordRepo
}

func NewRecordService(repo repo.RecordRepo) *RecordService {
	return &RecordService{
		repo: repo,
	}
}

// GetRecordByID возвращает запись по её ID
func (s *RecordService) GetRecordByID(id int64, logger zerolog.Logger) (*domain.Record, error) {
	logger.Debug().Int64("record_id", id).Msg("[GetRecordByID] Getting record by ID")

	if id <= 0 {
		logger.Error().Int64("record_id", id).Msg("[GetRecordByID] Invalid record ID")
		return nil, errors.New("record ID must be greater than 0")
	}

	record, err := s.repo.GetByID(id)
	if err != nil {
		logger.Error().Err(err).Int64("record_id", id).Msg("[GetRecordByID] Failed to get record")
		return nil, fmt.Errorf("failed to get record with ID %d: %w", id, err)
	}

	if record == nil {
		logger.Warn().Int64("record_id", id).Msg("[GetRecordByID] Record not found")
		return nil, fmt.Errorf("record with ID %d not found", id)
	}

	logger.Debug().
		Int64("record_id", record.ID).
		Int64("table_id", record.TableID).
		Int("position", record.Position).
		Msg("[GetRecordByID] Record retrieved successfully")

	return record, nil
}

// CreateRecord создает новую запись
func (s *RecordService) CreateRecord(tableID int64, position int, logger zerolog.Logger) (*domain.Record, error) {
	logger.Debug().
		Int64("table_id", tableID).
		Int("position", position).
		Msg("[CreateRecord] Creating new record")

	// Validate input parameters
	if err := s.validateRecordInput(tableID, position, logger); err != nil {
		return nil, err
	}

	// Create record object
	record := &domain.Record{
		TableID:   tableID,
		Position:  position,
		CreatedAt: time.Now(),
	}

	// Save to database
	err := s.repo.Create(record)
	if err != nil {
		logger.Error().Err(err).
			Int64("table_id", tableID).
			Int("position", position).
			Msg("[CreateRecord] Failed to create record")
		return nil, fmt.Errorf("failed to create record in table %d: %w", tableID, err)
	}

	logger.Info().
		Int64("record_id", record.ID).
		Int64("table_id", record.TableID).
		Int("position", record.Position).
		Msg("[CreateRecord] Record created successfully")

	return record, nil
}

// UpdateRecord обновляет существующую запись
func (s *RecordService) UpdateRecord(id int64, position *int, logger zerolog.Logger) (*domain.Record, error) {
	logger.Debug().
		Int64("record_id", id).
		Msg("[UpdateRecord] Updating record")

	// Validate record ID
	if id <= 0 {
		logger.Error().Int64("record_id", id).Msg("[UpdateRecord] Invalid record ID")
		return nil, errors.New("record ID must be greater than 0")
	}

	// Get existing record
	existingRecord, err := s.repo.GetByID(id)
	if err != nil {
		logger.Error().Err(err).Int64("record_id", id).Msg("[UpdateRecord] Failed to get existing record")
		return nil, fmt.Errorf("failed to get record with ID %d: %w", id, err)
	}

	if existingRecord == nil {
		logger.Warn().Int64("record_id", id).Msg("[UpdateRecord] Record not found")
		return nil, fmt.Errorf("record with ID %d not found", id)
	}

	// Update fields if provided
	if position != nil {
		if *position < 0 {
			logger.Error().Int("position", *position).Msg("[UpdateRecord] Invalid position")
			return nil, errors.New("record position cannot be negative")
		}
		existingRecord.Position = *position
	}

	// Save changes
	err = s.repo.Update(existingRecord)
	if err != nil {
		logger.Error().Err(err).
			Int64("record_id", id).
			Msg("[UpdateRecord] Failed to update record")
		return nil, fmt.Errorf("failed to update record with ID %d: %w", id, err)
	}

	logger.Info().
		Int64("record_id", existingRecord.ID).
		Int64("table_id", existingRecord.TableID).
		Int("position", existingRecord.Position).
		Msg("[UpdateRecord] Record updated successfully")

	return existingRecord, nil
}

// DeleteRecord удаляет запись по ID
func (s *RecordService) DeleteRecord(id int64, logger zerolog.Logger) error {
	logger.Debug().Int64("record_id", id).Msg("[DeleteRecord] Deleting record")

	if id <= 0 {
		logger.Error().Int64("record_id", id).Msg("[DeleteRecord] Invalid record ID")
		return errors.New("record ID must be greater than 0")
	}

	// Check if record exists before deletion
	existingRecord, err := s.repo.GetByID(id)
	if err != nil {
		logger.Error().Err(err).Int64("record_id", id).Msg("[DeleteRecord] Failed to get record for deletion")
		return fmt.Errorf("failed to verify record existence with ID %d: %w", id, err)
	}

	if existingRecord == nil {
		logger.Warn().Int64("record_id", id).Msg("[DeleteRecord] Record not found")
		return fmt.Errorf("record with ID %d not found", id)
	}

	// Delete the record
	err = s.repo.Delete(id)
	if err != nil {
		logger.Error().Err(err).Int64("record_id", id).Msg("[DeleteRecord] Failed to delete record")
		return fmt.Errorf("failed to delete record with ID %d: %w", id, err)
	}

	logger.Info().Int64("record_id", id).Msg("[DeleteRecord] Record deleted successfully")
	return nil
}

// GetRecordCells возвращает все ячейки записи по ID записи
func (s *RecordService) GetRecordCells(recordID int64, logger zerolog.Logger) ([]*domain.Cell, error) {
	logger.Debug().Int64("record_id", recordID).Msg("[GetRecordCells] Getting record cells")

	if recordID <= 0 {
		logger.Error().Int64("record_id", recordID).Msg("[GetRecordCells] Invalid record ID")
		return nil, errors.New("record ID must be greater than 0")
	}

	// Verify record exists
	if err := s.validateRecordExists(recordID, logger); err != nil {
		return nil, err
	}

	cells, err := s.repo.GetAllCells(recordID)
	if err != nil {
		logger.Error().Err(err).Int64("record_id", recordID).Msg("[GetRecordCells] Failed to get cells")
		return nil, fmt.Errorf("failed to get cells for record %d: %w", recordID, err)
	}

	if cells == nil {
		cells = []*domain.Cell{}
	}

	logger.Debug().
		Int64("record_id", recordID).
		Int("count", len(cells)).
		Msg("[GetRecordCells] Cells retrieved successfully")

	return cells, nil
}

// GetRecordsBySourceID возвращает все записи, связанные как источник
func (s *RecordService) GetRecordsBySourceID(recordID int64, logger zerolog.Logger) ([]*domain.Record, error) {
	logger.Debug().Int64("source_record_id", recordID).Msg("[GetRecordsBySourceID] Getting records by source ID")

	if recordID <= 0 {
		logger.Error().Int64("source_record_id", recordID).Msg("[GetRecordsBySourceID] Invalid record ID")
		return nil, errors.New("record ID must be greater than 0")
	}

	// Verify record exists
	if err := s.validateRecordExists(recordID, logger); err != nil {
		return nil, err
	}

	records, err := s.repo.GetAllRecordsBySourceID(recordID)
	if err != nil {
		logger.Error().Err(err).Int64("source_record_id", recordID).Msg("[GetRecordsBySourceID] Failed to get records")
		return nil, fmt.Errorf("failed to get records by source ID %d: %w", recordID, err)
	}

	if records == nil {
		records = []*domain.Record{}
	}

	logger.Debug().
		Int64("source_record_id", recordID).
		Int("count", len(records)).
		Msg("[GetRecordsBySourceID] Records retrieved successfully")

	return records, nil
}

// GetRecordsByTargetID возвращает все записи, связанные как цель
func (s *RecordService) GetRecordsByTargetID(recordID int64, logger zerolog.Logger) ([]*domain.Record, error) {
	logger.Debug().Int64("target_record_id", recordID).Msg("[GetRecordsByTargetID] Getting records by target ID")

	if recordID <= 0 {
		logger.Error().Int64("target_record_id", recordID).Msg("[GetRecordsByTargetID] Invalid record ID")
		return nil, errors.New("record ID must be greater than 0")
	}

	// Verify record exists
	if err := s.validateRecordExists(recordID, logger); err != nil {
		return nil, err
	}

	records, err := s.repo.GetAllRecordsByTargetID(recordID)
	if err != nil {
		logger.Error().Err(err).Int64("target_record_id", recordID).Msg("[GetRecordsByTargetID] Failed to get records")
		return nil, fmt.Errorf("failed to get records by target ID %d: %w", recordID, err)
	}

	if records == nil {
		records = []*domain.Record{}
	}

	logger.Debug().
		Int64("target_record_id", recordID).
		Int("count", len(records)).
		Msg("[GetRecordsByTargetID] Records retrieved successfully")

	return records, nil
}

// validateRecordExists проверяет, существует ли запись
func (s *RecordService) validateRecordExists(id int64, logger zerolog.Logger) error {
	logger.Debug().Int64("record_id", id).Msg("[validateRecordExists] Validating record existence")

	if id <= 0 {
		logger.Error().Int64("record_id", id).Msg("[validateRecordExists] Invalid record ID")
		return errors.New("record ID must be greater than 0")
	}

	record, err := s.repo.GetByID(id)
	if err != nil {
		logger.Error().Err(err).Int64("record_id", id).Msg("[validateRecordExists] Failed to get record")
		return fmt.Errorf("failed to validate record existence with ID %d: %w", id, err)
	}

	if record == nil {
		logger.Warn().Int64("record_id", id).Msg("[validateRecordExists] Record not found")
		return fmt.Errorf("record with ID %d does not exist", id)
	}

	logger.Debug().Int64("record_id", id).Msg("[validateRecordExists] Record exists")
	return nil
}

// validateRecordInput валидирует входные параметры записи
func (s *RecordService) validateRecordInput(tableID int64, position int, logger zerolog.Logger) error {
	if tableID <= 0 {
		logger.Error().Int64("table_id", tableID).Msg("[validateRecordInput] Invalid table ID")
		return errors.New("table ID must be greater than 0")
	}

	if position < 0 {
		logger.Error().Int("position", position).Msg("[validateRecordInput] Invalid position")
		return errors.New("record position cannot be negative")
	}

	return nil
}
