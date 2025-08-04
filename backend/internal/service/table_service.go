package service

import (
	"dater/backend/internal/domain"
	"dater/backend/internal/repo"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog"
)

type TableService struct {
	repo repo.TableRepo
}

func NewTableService(repo repo.TableRepo) *TableService {
	return &TableService{
		repo: repo,
	}
}

// GetTableByID возвращает таблицу по её ID
func (s *TableService) GetTableByID(id int64, logger zerolog.Logger) (*domain.Table, error) {
	logger.Debug().Int64("table_id", id).Msg("[GetTableByID] Getting table by ID")

	if id <= 0 {
		logger.Error().Int64("table_id", id).Msg("[GetTableByID] Invalid table ID")
		return nil, errors.New("table ID must be greater than 0")
	}

	table, err := s.repo.GetByID(id)
	if err != nil {
		logger.Error().Err(err).Int64("table_id", id).Msg("[GetTableByID] Failed to get table")
		return nil, fmt.Errorf("failed to get table with ID %d: %w", id, err)
	}

	if table == nil {
		logger.Warn().Int64("table_id", id).Msg("[GetTableByID] Table not found")
		return nil, fmt.Errorf("table with ID %d not found", id)
	}

	logger.Debug().
		Int64("table_id", table.ID).
		Str("name", table.Name).
		Int64("base_id", table.BaseID).
		Msg("[GetTableByID] Table retrieved successfully")

	return table, nil
}

// CreateTable создает новую таблицу и возвращает её
func (s *TableService) CreateTable(name string, description *string, baseID int64, logger zerolog.Logger) (*domain.Table, error) {
	logger.Debug().
		Str("name", name).
		Int64("base_id", baseID).
		Msg("[CreateTable] Creating new table")

	// Validate input parameters
	if err := s.validateTableInput(name, baseID, logger); err != nil {
		return nil, err
	}

	// Create table object
	table := &domain.Table{
		Name:        name,
		Description: description,
		BaseID:      baseID,
		CreatedAt:   time.Now(),
	}

	// Save to database
	createdTable, err := s.repo.Create(table)
	if err != nil {
		logger.Error().Err(err).
			Str("name", name).
			Int64("base_id", baseID).
			Msg("[CreateTable] Failed to create table")
		return nil, fmt.Errorf("failed to create table '%s': %w", name, err)
	}

	if createdTable == nil {
		logger.Error().Msg("[CreateTable] Repository returned nil table")
		return nil, errors.New("failed to create table: repository returned nil")
	}

	logger.Info().
		Int64("table_id", createdTable.ID).
		Str("name", createdTable.Name).
		Int64("base_id", createdTable.BaseID).
		Msg("[CreateTable] Table created successfully")

	return createdTable, nil
}

// UpdateTable обновляет существующую таблицу и возвращает её
func (s *TableService) UpdateTable(id int64, name string, description *string, logger zerolog.Logger) (*domain.Table, error) {
	logger.Debug().
		Int64("table_id", id).
		Str("name", name).
		Msg("[UpdateTable] Updating table")

	// Validate input parameters
	if id <= 0 {
		logger.Error().Int64("table_id", id).Msg("[UpdateTable] Invalid table ID")
		return nil, errors.New("table ID must be greater than 0")
	}

	if name == "" {
		logger.Error().Msg("[UpdateTable] Table name cannot be empty")
		return nil, errors.New("table name cannot be empty")
	}

	// Get existing table
	existingTable, err := s.repo.GetByID(id)
	if err != nil {
		logger.Error().Err(err).Int64("table_id", id).Msg("[UpdateTable] Failed to get existing table")
		return nil, fmt.Errorf("failed to get table with ID %d: %w", id, err)
	}

	if existingTable == nil {
		logger.Warn().Int64("table_id", id).Msg("[UpdateTable] Table not found")
		return nil, fmt.Errorf("table with ID %d not found", id)
	}

	// Update fields
	existingTable.Name = name
	existingTable.Description = description

	// Save changes
	updatedTable, err := s.repo.Update(existingTable)
	if err != nil {
		logger.Error().Err(err).
			Int64("table_id", id).
			Str("name", name).
			Msg("[UpdateTable] Failed to update table")
		return nil, fmt.Errorf("failed to update table with ID %d: %w", id, err)
	}

	if updatedTable == nil {
		logger.Error().Msg("[UpdateTable] Repository returned nil table")
		return nil, errors.New("failed to update table: repository returned nil")
	}

	logger.Info().
		Int64("table_id", updatedTable.ID).
		Str("name", updatedTable.Name).
		Msg("[UpdateTable] Table updated successfully")

	return updatedTable, nil
}

// DeleteTable удаляет таблицу по ID. Возвращает ошибку, если удаление не удалось.
func (s *TableService) DeleteTable(id int64, logger zerolog.Logger) (int64, error) {
	logger.Debug().Int64("table_id", id).Msg("[DeleteTable] Deleting table")

	if id <= 0 {
		logger.Error().Int64("table_id", id).Msg("[DeleteTable] Invalid table ID")
		return 0, errors.New("table ID must be greater than 0")
	}

	// Check if table exists before deletion
	existingTable, err := s.repo.GetByID(id)
	if err != nil {
		logger.Error().Err(err).Int64("table_id", id).Msg("[DeleteTable] Failed to get table for deletion")
		return 0, fmt.Errorf("failed to verify table existence with ID %d: %w", id, err)
	}

	if existingTable == nil {
		logger.Warn().Int64("table_id", id).Msg("[DeleteTable] Table not found")
		return 0, fmt.Errorf("table with ID %d not found", id)
	}

	// Delete the table
	deletedID, err := s.repo.Delete(id)
	if err != nil {
		logger.Error().Err(err).Int64("table_id", id).Msg("[DeleteTable] Failed to delete table")
		return 0, fmt.Errorf("failed to delete table with ID %d: %w", id, err)
	}

	if deletedID != id {
		logger.Error().
			Int64("expected_id", id).
			Int64("deleted_id", deletedID).
			Msg("[DeleteTable] Unexpected deleted ID")
		return 0, fmt.Errorf("unexpected deleted ID: expected %d, got %d", id, deletedID)
	}

	logger.Info().Int64("table_id", id).Msg("[DeleteTable] Table deleted successfully")
	return deletedID, nil
}

// GetAllRecordsByTable возвращает все записи для конкретной таблицы
func (s *TableService) GetAllRecordsByTable(tableID int64, logger zerolog.Logger) ([]*domain.Record, error) {
	logger.Debug().Int64("table_id", tableID).Msg("[GetAllRecordsByTable] Getting all records by table")

	if tableID <= 0 {
		logger.Error().Int64("table_id", tableID).Msg("[GetAllRecordsByTable] Invalid table ID")
		return nil, errors.New("table ID must be greater than 0")
	}

	// Verify table exists
	table, err := s.repo.GetByID(tableID)
	if err != nil {
		logger.Error().Err(err).Int64("table_id", tableID).Msg("[GetAllRecordsByTable] Failed to verify table existence")
		return nil, fmt.Errorf("failed to verify table existence with ID %d: %w", tableID, err)
	}

	if table == nil {
		logger.Warn().Int64("table_id", tableID).Msg("[GetAllRecordsByTable] Table not found")
		return nil, fmt.Errorf("table with ID %d not found", tableID)
	}

	// Get records
	records, err := s.repo.GetAllRecords(tableID)
	if err != nil {
		logger.Error().Err(err).Int64("table_id", tableID).Msg("[GetAllRecordsByTable] Failed to get records")
		return nil, fmt.Errorf("failed to get records for table %d: %w", tableID, err)
	}

	if records == nil {
		records = []*domain.Record{}
	}

	logger.Debug().
		Int64("table_id", tableID).
		Int("count", len(records)).
		Msg("[GetAllRecordsByTable] Records retrieved successfully")

	return records, nil
}

// GetAllFieldsByTable возвращает все поля для конкретной таблицы
func (s *TableService) GetAllFieldsByTable(tableID int64, logger zerolog.Logger) ([]*domain.Field, error) {
	logger.Debug().Int64("table_id", tableID).Msg("[GetAllFieldsByTable] Getting all fields by table")

	if tableID <= 0 {
		logger.Error().Int64("table_id", tableID).Msg("[GetAllFieldsByTable] Invalid table ID")
		return nil, errors.New("table ID must be greater than 0")
	}

	// Verify table exists
	table, err := s.repo.GetByID(tableID)
	if err != nil {
		logger.Error().Err(err).Int64("table_id", tableID).Msg("[GetAllFieldsByTable] Failed to verify table existence")
		return nil, fmt.Errorf("failed to verify table existence with ID %d: %w", tableID, err)
	}

	if table == nil {
		logger.Warn().Int64("table_id", tableID).Msg("[GetAllFieldsByTable] Table not found")
		return nil, fmt.Errorf("table with ID %d not found", tableID)
	}

	// Get fields
	fields, err := s.repo.GetAllFields(tableID)
	if err != nil {
		logger.Error().Err(err).Int64("table_id", tableID).Msg("[GetAllFieldsByTable] Failed to get fields")
		return nil, fmt.Errorf("failed to get fields for table %d: %w", tableID, err)
	}

	if fields == nil {
		fields = []*domain.Field{}
	}

	logger.Debug().
		Int64("table_id", tableID).
		Int("count", len(fields)).
		Msg("[GetAllFieldsByTable] Fields retrieved successfully")

	return fields, nil
}

// validateTableInput валидирует входные данные таблицы
func (s *TableService) validateTableInput(name string, baseID int64, logger zerolog.Logger) error {
	if name == "" {
		logger.Error().Msg("[validateTableInput] Table name cannot be empty")
		return errors.New("table name cannot be empty")
	}

	if len(name) > 255 {
		logger.Error().Int("length", len(name)).Msg("[validateTableInput] Table name too long")
		return errors.New("table name cannot exceed 255 characters")
	}

	if baseID <= 0 {
		logger.Error().Int64("base_id", baseID).Msg("[validateTableInput] Invalid base ID")
		return errors.New("base ID must be greater than 0")
	}

	return nil
}
