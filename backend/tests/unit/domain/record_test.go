package domain_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"clody/core/domain"
)

// ПОЗИТИВНЫЕ ТЕСТЫ - тестирование нормальных сценариев
func TestRecordStruct_Positive(t *testing.T) {
	t.Run("Record struct creation with valid data", func(t *testing.T) {
		// Arrange
		record := &domain.Record{
			ID:        1,
			TableID:   10,
			Position:  1,
			CreatedAt: time.Now(),
		}

		// Act & Assert
		assert.NotNil(t, record, "Record struct should not be nil")
		assert.Equal(t, int64(1), record.ID, "ID should be 1")
		assert.Equal(t, int64(10), record.TableID, "TableID should be 10")
		assert.Equal(t, 1, record.Position, "Position should be 1")
		assert.NotZero(t, record.CreatedAt, "CreatedAt should not be zero")
	})

	t.Run("Record with relations", func(t *testing.T) {
		// Arrange
		table := &domain.Table{ID: 100, Name: "Test Table"}
		cell1 := &domain.Cell{ID: 1, FieldID: 1, RecordID: 1}
		cell2 := &domain.Cell{ID: 2, FieldID: 2, RecordID: 1}

		record := &domain.Record{
			ID:        2,
			TableID:   100,
			Position:  2,
			CreatedAt: time.Now(),
			Table:     table,
			Cells:     []*domain.Cell{cell1, cell2},
		}

		// Act & Assert
		assert.NotNil(t, record, "Record should not be nil")
		assert.NotNil(t, record.Table, "Table relation should not be nil")
		assert.Equal(t, 2, len(record.Cells), "Should have 2 cells")
		
		// Проверка таблицы
		assert.Equal(t, int64(100), record.Table.ID, "Table ID should match")
		assert.Equal(t, "Test Table", record.Table.Name, "Table name should match")
		
		// Проверка ячеек
		assert.Equal(t, int64(1), record.Cells[0].ID, "First cell ID should be 1")
		assert.Equal(t, int64(1), record.Cells[0].RecordID, "First cell RecordID should match")
		assert.Equal(t, int64(2), record.Cells[1].ID, "Second cell ID should be 2")
		assert.Equal(t, int64(1), record.Cells[1].RecordID, "Second cell RecordID should match")
	})

	t.Run("Record with empty cells slice", func(t *testing.T) {
		// Arrange
		record := &domain.Record{
			ID:        3,
			TableID:   20,
			Position:  3,
			CreatedAt: time.Now(),
			Cells:     []*domain.Cell{}, // Пустой слайс
		}

		// Act & Assert
		assert.NotNil(t, record, "Record should not be nil")
		assert.Empty(t, record.Cells, "Cells should be empty")
		assert.NotNil(t, record.Cells, "Cells slice should not be nil")
	})
}

// НЕГАТИВНЫЕ ТЕСТЫ - тестирование граничных и "плохих" сценариев
func TestRecordStruct_Negative(t *testing.T) {
	t.Run("Record with zero values", func(t *testing.T) {
		// Arrange
		record := &domain.Record{
			ID:        0,
			TableID:   0,
			Position:  0,
			CreatedAt: time.Time{}, // Zero time
		}

		// Act & Assert
		assert.NotNil(t, record, "Record should not be nil")
		assert.Equal(t, int64(0), record.ID, "ID can be zero")
		assert.Equal(t, int64(0), record.TableID, "TableID can be zero")
		assert.Equal(t, 0, record.Position, "Position can be zero")
		assert.Equal(t, time.Time{}, record.CreatedAt, "CreatedAt can be zero time")
	})

	t.Run("Record with nil relations", func(t *testing.T) {
		// Arrange
		record := &domain.Record{
			ID:        4,
			TableID:   30,
			Position:  4,
			CreatedAt: time.Now(),
			Table:     nil, // Nil table
			Cells:     nil, // Nil cells
		}

		// Act & Assert
		assert.NotNil(t, record, "Record should not be nil")
		assert.Nil(t, record.Table, "Table can be nil")
		assert.Nil(t, record.Cells, "Cells can be nil")
	})

	t.Run("Record with negative position", func(t *testing.T) {
		// Arrange
		record := &domain.Record{
			ID:        5,
			TableID:   40,
			Position:  -1, // Negative position
			CreatedAt: time.Now(),
		}

		// Act & Assert
		assert.NotNil(t, record, "Record should not be nil")
		assert.Equal(t, -1, record.Position, "Negative position should be allowed")
	})

	t.Run("Record with negative IDs", func(t *testing.T) {
		// Arrange
		record := &domain.Record{
			ID:        -1,
			TableID:   -10,
			Position:  5,
			CreatedAt: time.Now(),
		}

		// Act & Assert
		assert.NotNil(t, record, "Record should not be nil")
		assert.Equal(t, int64(-1), record.ID, "Negative ID should be supported")
		assert.Equal(t, int64(-10), record.TableID, "Negative TableID should be supported")
	})
}

// ГРАНИЧНЫЕ ТЕСТЫ - тестирование крайних значений
func TestRecordStruct_Boundary(t *testing.T) {
	t.Run("Record with maximum int64 values", func(t *testing.T) {
		// Arrange
		record := &domain.Record{
			ID:        9223372036854775807,        // math.MaxInt64
			TableID:   9223372036854775807,       // math.MaxInt64
			Position:  2147483647,                // math.MaxInt32
			CreatedAt: time.Now(),
		}

		// Act & Assert
		assert.NotNil(t, record, "Record should not be nil")
		assert.Equal(t, int64(9223372036854775807), record.ID, "Max ID should be supported")
		assert.Equal(t, int64(9223372036854775807), record.TableID, "Max TableID should be supported")
		assert.Equal(t, 2147483647, record.Position, "Max position should be supported")
	})

	t.Run("Record with minimum int64 values", func(t *testing.T) {
		// Arrange
		record := &domain.Record{
			ID:        -9223372036854775808,       // math.MinInt64
			TableID:   -9223372036854775808,      // math.MinInt64
			Position:  -2147483648,               // math.MinInt32
			CreatedAt: time.Time{},
		}

		// Act & Assert
		assert.NotNil(t, record, "Record should not be nil")
		assert.Equal(t, int64(-9223372036854775808), record.ID, "Min ID should be supported")
		assert.Equal(t, int64(-9223372036854775808), record.TableID, "Min TableID should be supported")
		assert.Equal(t, -2147483648, record.Position, "Min position should be supported")
	})

	t.Run("Record with many cells", func(t *testing.T) {
		// Arrange
		var cells []*domain.Cell
		
		// Создаем много ячеек
		for i := 0; i < 10; i++ {
			cells = append(cells, &domain.Cell{ID: int64(i), RecordID: 6})
		}

		record := &domain.Record{
			ID:        6,
			TableID:   50,
			Position:  6,
			CreatedAt: time.Now(),
			Cells:     cells,
		}

		// Act & Assert
		assert.NotNil(t, record, "Record should not be nil")
		assert.Equal(t, 10, len(record.Cells), "Should have 10 cells")
	})
}

// ТЕСТЫ ДЛЯ ОСОБЫХ СЛУЧАЕВ
func TestRecordStruct_SpecialCases(t *testing.T) {
	t.Run("Record with same ID and TableID", func(t *testing.T) {
		// Arrange
		record := &domain.Record{
			ID:        100,
			TableID:   100, // Same as ID
			Position:  7,
			CreatedAt: time.Now(),
		}

		// Act & Assert
		assert.NotNil(t, record, "Record should not be nil")
		assert.Equal(t, record.ID, record.TableID, "ID and TableID can be equal")
	})

	t.Run("Record with zero position", func(t *testing.T) {
		// Arrange
		record := &domain.Record{
			ID:        7,
			TableID:   60,
			Position:  0, // Zero position
			CreatedAt: time.Now(),
		}

		// Act & Assert
		assert.NotNil(t, record, "Record should not be nil")
		assert.Equal(t, 0, record.Position, "Zero position should be allowed")
	})

	t.Run("Record with very large position", func(t *testing.T) {
		// Arrange
		largePosition := 1000000
		
		record := &domain.Record{
			ID:        8,
			TableID:   70,
			Position:  largePosition,
			CreatedAt: time.Now(),
		}

		// Act & Assert
		assert.NotNil(t, record, "Record should not be nil")
		assert.Equal(t, largePosition, record.Position, "Large position should be supported")
	})

	t.Run("Record with same position as ID", func(t *testing.T) {
		// Arrange
		sameValue := 200
		
		record := &domain.Record{
			ID:        int64(sameValue),
			TableID:   80,
			Position:  sameValue, // Same as ID
			CreatedAt: time.Now(),
		}

		// Act & Assert
		assert.NotNil(t, record, "Record should not be nil")
		assert.Equal(t, int64(sameValue), record.ID, "ID should match")
		assert.Equal(t, sameValue, record.Position, "Position should match ID")
	})
}

// ТЕСТЫ СВЯЗЕЙ И ОГРАНИЧЕНИЙ
func TestRecordStruct_Relations(t *testing.T) {
	t.Run("Record with partial relations", func(t *testing.T) {
		// Arrange
		table := &domain.Table{ID: 200, Name: "Partial Table"}

		record := &domain.Record{
			ID:        9,
			TableID:   200,
			Position:  9,
			CreatedAt: time.Now(),
			Table:     table,    // Only table relation
			Cells:     nil,      // No cells
		}

		// Act & Assert
		assert.NotNil(t, record, "Record should not be nil")
		assert.NotNil(t, record.Table, "Table should not be nil")
		assert.Nil(t, record.Cells, "Cells can be nil")
		
		assert.Equal(t, "Partial Table", record.Table.Name, "Table name should match")
	})

	t.Run("Record with cells but no table", func(t *testing.T) {
		// Arrange
		cell1 := &domain.Cell{ID: 10, FieldID: 10, RecordID: 10}
		cell2 := &domain.Cell{ID: 11, FieldID: 11, RecordID: 10}

		record := &domain.Record{
			ID:        10,
			TableID:   90,
			Position:  10,
			CreatedAt: time.Now(),
			Table:     nil,              // No table relation
			Cells:     []*domain.Cell{cell1, cell2}, // Only cells
		}

		// Act & Assert
		assert.NotNil(t, record, "Record should not be nil")
		assert.Nil(t, record.Table, "Table can be nil")
		assert.NotNil(t, record.Cells, "Cells should not be nil")
		assert.Equal(t, 2, len(record.Cells), "Should have 2 cells")
	})

	t.Run("Record with no relations", func(t *testing.T) {
		// Arrange
		record := &domain.Record{
			ID:        11,
			TableID:   100,
			Position:  11,
			CreatedAt: time.Now(),
			Table:     nil,    // No table
			Cells:     nil,    // No cells
		}

		// Act & Assert
		assert.NotNil(t, record, "Record should not be nil")
		assert.Nil(t, record.Table, "Table should be nil")
		assert.Nil(t, record.Cells, "Cells should be nil")
	})
}