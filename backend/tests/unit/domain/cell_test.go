package domain_test

import (
	"testing"
	"time"

	"clody/core/domain"
	"clody/utils"

	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

// ПОЗИТИВНЫЕ ТЕСТЫ - тестирование нормальных сценариев
func TestCellStruct_Positive(t *testing.T) {
	t.Run("Cell struct creation with valid data", func(t *testing.T) {
		// Arrange
		jsonValue := utils.CreateJSONFromString(t, `{"text": "test value"}`)

		cell := &domain.Cell{
			ID:        1,
			FieldID:   10,
			RecordID:  20,
			Position:  1,
			Value:     jsonValue,
			CreatedAt: time.Now(),
		}

		// Act & Assert
		assert.NotNil(t, cell, "Cell struct should not be nil")
		assert.Equal(t, int64(1), cell.ID, "ID should be 1")
		assert.Equal(t, int64(10), cell.FieldID, "FieldID should be 10")
		assert.Equal(t, int64(20), cell.RecordID, "RecordID should be 20")
		assert.Equal(t, 1, cell.Position, "Position should be 1")
		assert.Equal(t, jsonValue, cell.Value, "Value should match")
		assert.NotZero(t, cell.CreatedAt, "CreatedAt should not be zero")
	})

	t.Run("Cell with different JSON values", func(t *testing.T) {
		// Arrange
		testCases := []struct {
			name    string
			jsonStr string
		}{
			{
				name:    "String value",
				jsonStr: `{"value": "text"}`,
			},
			{
				name:    "Number value",
				jsonStr: `{"value": 123}`,
			},
			{
				name:    "Boolean value",
				jsonStr: `{"value": true}`,
			},
			{
				name:    "Array value",
				jsonStr: `{"value": [1, 2, 3]}`,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				jsonValue := utils.CreateJSONFromString(t, tc.jsonStr)

				cell := &domain.Cell{
					ID:        2,
					FieldID:   15,
					RecordID:  25,
					Position:  2,
					Value:     jsonValue,
					CreatedAt: time.Now(),
				}

				assert.NotNil(t, cell, "Cell should not be nil")
				assert.Equal(t, jsonValue, cell.Value, "JSON value should match")
			})
		}
	})

	t.Run("Cell with relations", func(t *testing.T) {
		// Arrange
		jsonValue := utils.CreateJSONFromString(t, `{"data": "test"}`)
		field := &domain.Field{ID: 100, Name: "Test Field"}
		record := &domain.Record{ID: 200, TableID: 1}

		cell := &domain.Cell{
			ID:        3,
			FieldID:   100,
			RecordID:  200,
			Position:  0,
			Value:     jsonValue,
			CreatedAt: time.Now(),
			Field:     field,
			Record:    record,
		}

		// Act & Assert
		assert.NotNil(t, cell, "Cell should not be nil")
		assert.NotNil(t, cell.Field, "Field relation should not be nil")
		assert.NotNil(t, cell.Record, "Record relation should not be nil")
		assert.Equal(t, int64(100), cell.Field.ID, "Field ID should match")
		assert.Equal(t, int64(200), cell.Record.ID, "Record ID should match")
	})
}

// НЕГАТИВНЫЕ ТЕСТЫ - тестирование граничных и "плохих" сценариев
func TestCellStruct_Negative(t *testing.T) {
	t.Run("Cell with zero values", func(t *testing.T) {
		// Arrange
		jsonValue := utils.CreateJSONFromString(t, `{}`)

		cell := &domain.Cell{
			ID:        0,
			FieldID:   0,
			RecordID:  0,
			Position:  0,
			Value:     jsonValue,
			CreatedAt: time.Time{}, // Zero time
		}

		// Act & Assert
		assert.NotNil(t, cell, "Cell should not be nil")
		assert.Equal(t, int64(0), cell.ID, "ID can be zero")
		assert.Equal(t, int64(0), cell.FieldID, "FieldID can be zero")
		assert.Equal(t, int64(0), cell.RecordID, "RecordID can be zero")
		assert.Equal(t, 0, cell.Position, "Position can be zero")
		assert.Equal(t, time.Time{}, cell.CreatedAt, "CreatedAt can be zero time")
	})

	t.Run("Cell with negative position", func(t *testing.T) {
		// Arrange
		jsonValue := utils.CreateJSONFromString(t, `{"negative": true}`)

		cell := &domain.Cell{
			ID:        4,
			FieldID:   30,
			RecordID:  40,
			Position:  -1, // Negative position
			Value:     jsonValue,
			CreatedAt: time.Now(),
		}

		// Act & Assert
		assert.NotNil(t, cell, "Cell should not be nil")
		assert.Equal(t, -1, cell.Position, "Negative position should be allowed")
	})

	t.Run("Cell with nil relations", func(t *testing.T) {
		// Arrange
		jsonValue := utils.CreateJSONFromString(t, `{"nil_relations": true}`)

		cell := &domain.Cell{
			ID:        5,
			FieldID:   35,
			RecordID:  45,
			Position:  3,
			Value:     jsonValue,
			CreatedAt: time.Now(),
			Field:     nil, // Nil field relation
			Record:    nil, // Nil record relation
		}

		// Act & Assert
		assert.NotNil(t, cell, "Cell should not be nil")
		assert.Nil(t, cell.Field, "Field can be nil")
		assert.Nil(t, cell.Record, "Record can be nil")
	})

	t.Run("Cell with empty JSON", func(t *testing.T) {
		// Arrange
		jsonValue := utils.CreateJSONFromString(t, `{}`)

		cell := &domain.Cell{
			ID:        6,
			FieldID:   50,
			RecordID:  60,
			Position:  4,
			Value:     jsonValue,
			CreatedAt: time.Now(),
		}

		// Act & Assert
		assert.NotNil(t, cell, "Cell should not be nil")
		assert.Equal(t, jsonValue, cell.Value, "Empty JSON should be valid")
	})
}

// ГРАНИЧНЫЕ ТЕСТЫ - тестирование крайних значений
func TestCellStruct_Boundary(t *testing.T) {
	t.Run("Cell with maximum int64 values", func(t *testing.T) {
		// Arrange
		jsonValue := utils.CreateJSONFromString(t, `{"max_values": true}`)

		cell := &domain.Cell{
			ID:        9223372036854775807, // math.MaxInt64
			FieldID:   9223372036854775807, // math.MaxInt64
			RecordID:  9223372036854775807, // math.MaxInt64
			Position:  2147483647,          // math.MaxInt32 (max for int)
			Value:     jsonValue,
			CreatedAt: time.Now(),
		}

		// Act & Assert
		assert.NotNil(t, cell, "Cell should not be nil")
		assert.Equal(t, int64(9223372036854775807), cell.ID, "Max ID should be supported")
		assert.Equal(t, int64(9223372036854775807), cell.FieldID, "Max FieldID should be supported")
		assert.Equal(t, int64(9223372036854775807), cell.RecordID, "Max RecordID should be supported")
		assert.Equal(t, 2147483647, cell.Position, "Max int position should be supported")
	})

	t.Run("Cell with minimum int64 values", func(t *testing.T) {
		// Arrange
		jsonValue := utils.CreateJSONFromString(t, `{"min_values": true}`)

		cell := &domain.Cell{
			ID:        -9223372036854775808, // math.MinInt64
			FieldID:   -9223372036854775808, // math.MinInt64
			RecordID:  -9223372036854775808, // math.MinInt64
			Position:  -2147483648,          // math.MinInt32 (min for int)
			Value:     jsonValue,
			CreatedAt: time.Time{},
		}

		// Act & Assert
		assert.NotNil(t, cell, "Cell should not be nil")
		assert.Equal(t, int64(-9223372036854775808), cell.ID, "Min ID should be supported")
		assert.Equal(t, int64(-9223372036854775808), cell.FieldID, "Min FieldID should be supported")
		assert.Equal(t, int64(-9223372036854775808), cell.RecordID, "Min RecordID should be supported")
		assert.Equal(t, -2147483648, cell.Position, "Min int position should be supported")
	})

	t.Run("Cell with complex JSON structure", func(t *testing.T) {
		// Arrange
		complexJSON := `{
			"nested": {
				"array": [1, 2, {"key": "value"}],
				"string": "test string",
				"number": 42,
				"boolean": true,
				"null": null
			},
			"simple": "value"
		}`

		jsonValue := utils.CreateJSONFromString(t, complexJSON)

		cell := &domain.Cell{
			ID:        7,
			FieldID:   70,
			RecordID:  80,
			Position:  5,
			Value:     jsonValue,
			CreatedAt: time.Now(),
		}

		// Act & Assert
		assert.NotNil(t, cell, "Cell should not be nil")
		assert.Equal(t, jsonValue, cell.Value, "Complex JSON should be stored correctly")
	})
}

// ТЕСТЫ ДЛЯ ОСОБЫХ СЛУЧАЕВ
func TestCellStruct_SpecialCases(t *testing.T) {
	t.Run("Cell with same FieldID and RecordID", func(t *testing.T) {
		// Arrange
		jsonValue := utils.CreateJSONFromString(t, `{"same_ids": true}`)

		cell := &domain.Cell{
			ID:        8,
			FieldID:   100,
			RecordID:  100, // Same as FieldID
			Position:  6,
			Value:     jsonValue,
			CreatedAt: time.Now(),
		}

		// Act & Assert
		assert.NotNil(t, cell, "Cell should not be nil")
		assert.Equal(t, cell.FieldID, cell.RecordID, "FieldID and RecordID can be equal")
	})

	t.Run("Cell with zero JSON value", func(t *testing.T) {
		// Arrange
		var jsonValue datatypes.JSON // Zero value

		cell := &domain.Cell{
			ID:        9,
			FieldID:   110,
			RecordID:  120,
			Position:  7,
			Value:     jsonValue, // Zero value
			CreatedAt: time.Now(),
		}

		// Act & Assert
		assert.NotNil(t, cell, "Cell should not be nil")
		// Note: Zero JSON value is valid in Go, but might need validation elsewhere
	})
}
