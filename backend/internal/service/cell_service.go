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

type CellService struct {
	repo repo.CellRepo
}

func NewCellService(repo repo.CellRepo) *CellService {
	return &CellService{
		repo: repo,
	}
}

// GetCellByID возвращает ячейку по её ID
func (s *CellService) GetCellByID(id int64, logger zerolog.Logger) (*domain.Cell, error) {
	logger.Debug().Int64("cell_id", id).Msg("[GetCellByID] Getting cell by ID")

	if id <= 0 {
		logger.Error().Int64("cell_id", id).Msg("[GetCellByID] Invalid cell ID")
		return nil, errors.New("cell ID must be greater than 0")
	}

	cell, err := s.repo.GetByID(id)
	if err != nil {
		logger.Error().Err(err).Int64("cell_id", id).Msg("[GetCellByID] Failed to get cell")
		return nil, fmt.Errorf("failed to get cell with ID %d: %w", id, err)
	}

	if cell == nil {
		logger.Warn().Int64("cell_id", id).Msg("[GetCellByID] Cell not found")
		return nil, fmt.Errorf("cell with ID %d not found", id)
	}

	logger.Debug().
		Int64("cell_id", cell.ID).
		Int64("field_id", cell.FieldID).
		Int64("record_id", cell.RecordID).
		Int("position", cell.Position).
		Msg("[GetCellByID] Cell retrieved successfully")

	return cell, nil
}

// CreateCell создает новую ячейку
func (s *CellService) CreateCell(fieldID, recordID int64, position int, value datatypes.JSON, logger zerolog.Logger) (*domain.Cell, error) {
	logger.Debug().
		Int64("field_id", fieldID).
		Int64("record_id", recordID).
		Int("position", position).
		Msg("[CreateCell] Creating new cell")

	// Validate input parameters
	if err := s.validateCellInput(fieldID, recordID, position, logger); err != nil {
		return nil, err
	}

	// Create cell object
	cell := &domain.Cell{
		FieldID:   fieldID,
		RecordID:  recordID,
		Position:  position,
		Value:     value,
		CreatedAt: time.Now(),
	}

	// Save to database
	err := s.repo.Create(cell)
	if err != nil {
		logger.Error().Err(err).
			Int64("field_id", fieldID).
			Int64("record_id", recordID).
			Msg("[CreateCell] Failed to create cell")
		return nil, fmt.Errorf("failed to create cell: %w", err)
	}

	logger.Info().
		Int64("cell_id", cell.ID).
		Int64("field_id", cell.FieldID).
		Int64("record_id", cell.RecordID).
		Int("position", cell.Position).
		Msg("[CreateCell] Cell created successfully")

	return cell, nil
}

// UpdateCell обновляет существующую ячейку
func (s *CellService) UpdateCell(id int64, position *int, value datatypes.JSON, logger zerolog.Logger) (*domain.Cell, error) {
	logger.Debug().
		Int64("cell_id", id).
		Msg("[UpdateCell] Updating cell")

	// Validate cell ID
	if id <= 0 {
		logger.Error().Int64("cell_id", id).Msg("[UpdateCell] Invalid cell ID")
		return nil, errors.New("cell ID must be greater than 0")
	}

	// Get existing cell
	existingCell, err := s.repo.GetByID(id)
	if err != nil {
		logger.Error().Err(err).Int64("cell_id", id).Msg("[UpdateCell] Failed to get existing cell")
		return nil, fmt.Errorf("failed to get cell with ID %d: %w", id, err)
	}

	if existingCell == nil {
		logger.Warn().Int64("cell_id", id).Msg("[UpdateCell] Cell not found")
		return nil, fmt.Errorf("cell with ID %d not found", id)
	}

	// Update fields if provided
	if position != nil {
		if *position < 0 {
			logger.Error().Int("position", *position).Msg("[UpdateCell] Invalid position")
			return nil, errors.New("cell position cannot be negative")
		}
		existingCell.Position = *position
	}

	if value != nil {
		existingCell.Value = value
	}

	// Save changes
	err = s.repo.Update(existingCell)
	if err != nil {
		logger.Error().Err(err).
			Int64("cell_id", id).
			Msg("[UpdateCell] Failed to update cell")
		return nil, fmt.Errorf("failed to update cell with ID %d: %w", id, err)
	}

	logger.Info().
		Int64("cell_id", existingCell.ID).
		Int64("field_id", existingCell.FieldID).
		Int64("record_id", existingCell.RecordID).
		Int("position", existingCell.Position).
		Msg("[UpdateCell] Cell updated successfully")

	return existingCell, nil
}

// DeleteCell удаляет ячейку по ID
func (s *CellService) DeleteCell(id int64, logger zerolog.Logger) error {
	logger.Debug().Int64("cell_id", id).Msg("[DeleteCell] Deleting cell")

	if id <= 0 {
		logger.Error().Int64("cell_id", id).Msg("[DeleteCell] Invalid cell ID")
		return errors.New("cell ID must be greater than 0")
	}

	// Check if cell exists before deletion
	existingCell, err := s.repo.GetByID(id)
	if err != nil {
		logger.Error().Err(err).Int64("cell_id", id).Msg("[DeleteCell] Failed to get cell for deletion")
		return fmt.Errorf("failed to verify cell existence with ID %d: %w", id, err)
	}

	if existingCell == nil {
		logger.Warn().Int64("cell_id", id).Msg("[DeleteCell] Cell not found")
		return fmt.Errorf("cell with ID %d not found", id)
	}

	// Delete the cell
	err = s.repo.Delete(id)
	if err != nil {
		logger.Error().Err(err).Int64("cell_id", id).Msg("[DeleteCell] Failed to delete cell")
		return fmt.Errorf("failed to delete cell with ID %d: %w", id, err)
	}

	logger.Info().Int64("cell_id", id).Msg("[DeleteCell] Cell deleted successfully")
	return nil
}

// GetCellsByRecordID возвращает все ячейки для конкретной записи
func (s *CellService) GetCellsByRecordID(recordID int64, logger zerolog.Logger) ([]*domain.Cell, error) {
	logger.Debug().Int64("record_id", recordID).Msg("[GetCellsByRecordID] Getting cells by record ID")

	if recordID <= 0 {
		logger.Error().Int64("record_id", recordID).Msg("[GetCellsByRecordID] Invalid record ID")
		return nil, errors.New("record ID must be greater than 0")
	}

	cells, err := s.repo.GetAllCellsByRecordID(recordID)
	if err != nil {
		logger.Error().Err(err).Int64("record_id", recordID).Msg("[GetCellsByRecordID] Failed to get cells")
		return nil, fmt.Errorf("failed to get cells for record %d: %w", recordID, err)
	}

	if cells == nil {
		cells = []*domain.Cell{}
	}

	logger.Debug().
		Int64("record_id", recordID).
		Int("count", len(cells)).
		Msg("[GetCellsByRecordID] Cells retrieved successfully")

	return cells, nil
}

// GetCellsByFieldID возвращает все ячейки для конкретного поля
func (s *CellService) GetCellsByFieldID(fieldID int64, logger zerolog.Logger) ([]*domain.Cell, error) {
	logger.Debug().Int64("field_id", fieldID).Msg("[GetCellsByFieldID] Getting cells by field ID")

	if fieldID <= 0 {
		logger.Error().Int64("field_id", fieldID).Msg("[GetCellsByFieldID] Invalid field ID")
		return nil, errors.New("field ID must be greater than 0")
	}

	cells, err := s.repo.GetAllCellsByFieldID(fieldID)
	if err != nil {
		logger.Error().Err(err).Int64("field_id", fieldID).Msg("[GetCellsByFieldID] Failed to get cells")
		return nil, fmt.Errorf("failed to get cells for field %d: %w", fieldID, err)
	}

	if cells == nil {
		cells = []*domain.Cell{}
	}

	logger.Debug().
		Int64("field_id", fieldID).
		Int("count", len(cells)).
		Msg("[GetCellsByFieldID] Cells retrieved successfully")

	return cells, nil
}

// GetCellByFieldAndRecord возвращает ячейку для конкретного поля и записи
func (s *CellService) GetCellByFieldAndRecord(fieldID, recordID int64, logger zerolog.Logger) (*domain.Cell, error) {
	logger.Debug().
		Int64("field_id", fieldID).
		Int64("record_id", recordID).
		Msg("[GetCellByFieldAndRecord] Getting cell by field and record ID")

	if fieldID <= 0 {
		logger.Error().Int64("field_id", fieldID).Msg("[GetCellByFieldAndRecord] Invalid field ID")
		return nil, errors.New("field ID must be greater than 0")
	}

	if recordID <= 0 {
		logger.Error().Int64("record_id", recordID).Msg("[GetCellByFieldAndRecord] Invalid record ID")
		return nil, errors.New("record ID must be greater than 0")
	}

	// Get all cells for the record
	cells, err := s.repo.GetAllCellsByRecordID(recordID)
	if err != nil {
		logger.Error().Err(err).Int64("record_id", recordID).Msg("[GetCellByFieldAndRecord] Failed to get cells")
		return nil, fmt.Errorf("failed to get cells for record %d: %w", recordID, err)
	}

	// Find cell with matching field ID
	for _, cell := range cells {
		if cell.FieldID == fieldID {
			logger.Debug().
				Int64("cell_id", cell.ID).
				Int64("field_id", fieldID).
				Int64("record_id", recordID).
				Msg("[GetCellByFieldAndRecord] Cell found")
			return cell, nil
		}
	}

	logger.Debug().
		Int64("field_id", fieldID).
		Int64("record_id", recordID).
		Msg("[GetCellByFieldAndRecord] Cell not found")
	return nil, fmt.Errorf("cell not found for field %d and record %d", fieldID, recordID)
}

// UpdateCellValue обновляет только значение ячейки
func (s *CellService) UpdateCellValue(id int64, value datatypes.JSON, logger zerolog.Logger) (*domain.Cell, error) {
	logger.Debug().
		Int64("cell_id", id).
		Msg("[UpdateCellValue] Updating cell value")

	// Validate cell ID
	if id <= 0 {
		logger.Error().Int64("cell_id", id).Msg("[UpdateCellValue] Invalid cell ID")
		return nil, errors.New("cell ID must be greater than 0")
	}

	// Get existing cell
	existingCell, err := s.repo.GetByID(id)
	if err != nil {
		logger.Error().Err(err).Int64("cell_id", id).Msg("[UpdateCellValue] Failed to get existing cell")
		return nil, fmt.Errorf("failed to get cell with ID %d: %w", id, err)
	}

	if existingCell == nil {
		logger.Warn().Int64("cell_id", id).Msg("[UpdateCellValue] Cell not found")
		return nil, fmt.Errorf("cell with ID %d not found", id)
	}

	// Update value
	existingCell.Value = value

	// Save changes
	err = s.repo.Update(existingCell)
	if err != nil {
		logger.Error().Err(err).
			Int64("cell_id", id).
			Msg("[UpdateCellValue] Failed to update cell value")
		return nil, fmt.Errorf("failed to update cell value with ID %d: %w", id, err)
	}

	logger.Info().
		Int64("cell_id", existingCell.ID).
		Msg("[UpdateCellValue] Cell value updated successfully")

	return existingCell, nil
}

// validateCellExists проверяет, существует ли ячейка
func (s *CellService) validateCellExists(id int64, logger zerolog.Logger) error {
	logger.Debug().Int64("cell_id", id).Msg("[validateCellExists] Validating cell existence")

	if id <= 0 {
		logger.Error().Int64("cell_id", id).Msg("[validateCellExists] Invalid cell ID")
		return errors.New("cell ID must be greater than 0")
	}

	cell, err := s.repo.GetByID(id)
	if err != nil {
		logger.Error().Err(err).Int64("cell_id", id).Msg("[validateCellExists] Failed to get cell")
		return fmt.Errorf("failed to validate cell existence with ID %d: %w", id, err)
	}

	if cell == nil {
		logger.Warn().Int64("cell_id", id).Msg("[validateCellExists] Cell not found")
		return fmt.Errorf("cell with ID %d does not exist", id)
	}

	logger.Debug().Int64("cell_id", id).Msg("[validateCellExists] Cell exists")
	return nil
}

// validateCellInput валидирует входные параметры ячейки
func (s *CellService) validateCellInput(fieldID, recordID int64, position int, logger zerolog.Logger) error {
	if fieldID <= 0 {
		logger.Error().Int64("field_id", fieldID).Msg("[validateCellInput] Invalid field ID")
		return errors.New("field ID must be greater than 0")
	}

	if recordID <= 0 {
		logger.Error().Int64("record_id", recordID).Msg("[validateCellInput] Invalid record ID")
		return errors.New("record ID must be greater than 0")
	}

	if position < 0 {
		logger.Error().Int("position", position).Msg("[validateCellInput] Invalid position")
		return errors.New("cell position cannot be negative")
	}

	return nil
}

// BulkUpdateCells обновляет несколько ячеек одновременно
func (s *CellService) BulkUpdateCells(updates []struct {
	ID    int64
	Value datatypes.JSON
}, logger zerolog.Logger) ([]*domain.Cell, error) {
	logger.Debug().
		Int("count", len(updates)).
		Msg("[BulkUpdateCells] Bulk updating cells")

	if len(updates) == 0 {
		logger.Warn().Msg("[BulkUpdateCells] No updates provided")
		return []*domain.Cell{}, nil
	}

	var updatedCells []*domain.Cell

	for _, update := range updates {
		cell, err := s.UpdateCellValue(update.ID, update.Value, logger)
		if err != nil {
			logger.Error().Err(err).
				Int64("cell_id", update.ID).
				Msg("[BulkUpdateCells] Failed to update cell")
			// Continue with other updates but log the error
			continue
		}
		updatedCells = append(updatedCells, cell)
	}

	logger.Info().
		Int("requested", len(updates)).
		Int("updated", len(updatedCells)).
		Msg("[BulkUpdateCells] Bulk update completed")

	return updatedCells, nil
}
