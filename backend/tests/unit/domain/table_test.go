package domain_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"clody/core/domain"
)

// ПОЗИТИВНЫЕ ТЕСТЫ - тестирование нормальных сценариев
func TestTableStruct_Positive(t *testing.T) {
	t.Run("Table struct creation with valid data", func(t *testing.T) {
		// Arrange
		description := "Test table description"
		
		table := &domain.Table{
			ID:          1,
			BaseID:      10,
			Name:        "Test Table",
			Description: &description,
			CreatedAt:   time.Now(),
		}

		// Act & Assert
		assert.NotNil(t, table, "Table struct should not be nil")
		assert.Equal(t, int64(1), table.ID, "ID should be 1")
		assert.Equal(t, int64(10), table.BaseID, "BaseID should be 10")
		assert.Equal(t, "Test Table", table.Name, "Name should be 'Test Table'")
		assert.Equal(t, description, *table.Description, "Description should match")
		assert.NotZero(t, table.CreatedAt, "CreatedAt should not be zero")
	})

	t.Run("Table with relations", func(t *testing.T) {
		// Arrange
		field1 := &domain.Field{ID: 1, Name: "Field 1"}
		field2 := &domain.Field{ID: 2, Name: "Field 2"}
		record1 := &domain.Record{ID: 1, TableID: 1}
		record2 := &domain.Record{ID: 2, TableID: 1}
		base := &domain.Base{ID: 100, Name: "Test Base"}

		table := &domain.Table{
			ID:          2,
			BaseID:      100,
			Name:        "Table with Relations",
			CreatedAt:   time.Now(),
			Fields:      []*domain.Field{field1, field2},
			Records:     []*domain.Record{record1, record2},
			Base:        base,
		}

		// Act & Assert
		assert.NotNil(t, table, "Table should not be nil")
		assert.Equal(t, 2, len(table.Fields), "Should have 2 fields")
		assert.Equal(t, 2, len(table.Records), "Should have 2 records")
		assert.NotNil(t, table.Base, "Base relation should not be nil")
		
		// Проверка полей
		assert.Equal(t, int64(1), table.Fields[0].ID, "First field ID should be 1")
		assert.Equal(t, "Field 1", table.Fields[0].Name, "First field name should match")
		assert.Equal(t, int64(2), table.Fields[1].ID, "Second field ID should be 2")
		assert.Equal(t, "Field 2", table.Fields[1].Name, "Second field name should match")
		
		// Проверка записей
		assert.Equal(t, int64(1), table.Records[0].ID, "First record ID should be 1")
		assert.Equal(t, int64(2), table.Records[1].ID, "Second record ID should be 2")
		
		// Проверка базы
		assert.Equal(t, int64(100), table.Base.ID, "Base ID should be 100")
		assert.Equal(t, "Test Base", table.Base.Name, "Base name should match")
	})

	t.Run("Table with empty slices", func(t *testing.T) {
		// Arrange
		table := &domain.Table{
			ID:        3,
			BaseID:    20,
			Name:      "Table with Empty Relations",
			CreatedAt: time.Now(),
			Fields:    []*domain.Field{}, // Пустой слайс
			Records:   []*domain.Record{}, // Пустой слайс
		}

		// Act & Assert
		assert.NotNil(t, table, "Table should not be nil")
		assert.Empty(t, table.Fields, "Fields should be empty")
		assert.Empty(t, table.Records, "Records should be empty")
		assert.NotNil(t, table.Fields, "Fields slice should not be nil")
		assert.NotNil(t, table.Records, "Records slice should not be nil")
	})
}

// НЕГАТИВНЫЕ ТЕСТЫ - тестирование граничных и "плохих" сценариев
func TestTableStruct_Negative(t *testing.T) {
	t.Run("Table with zero values", func(t *testing.T) {
		// Arrange
		table := &domain.Table{
			ID:        0,
			BaseID:    0,
			Name:      "",
			CreatedAt: time.Time{}, // Zero time
		}

		// Act & Assert
		assert.NotNil(t, table, "Table should not be nil")
		assert.Equal(t, int64(0), table.ID, "ID can be zero")
		assert.Equal(t, int64(0), table.BaseID, "BaseID can be zero")
		assert.Equal(t, "", table.Name, "Name can be empty")
		assert.Equal(t, time.Time{}, table.CreatedAt, "CreatedAt can be zero time")
	})

	t.Run("Table with nil description", func(t *testing.T) {
		// Arrange
		table := &domain.Table{
			ID:          4,
			BaseID:      30,
			Name:        "Table with Nil Description",
			Description: nil, // Nil description
			CreatedAt:   time.Now(),
		}

		// Act & Assert
		assert.NotNil(t, table, "Table should not be nil")
		assert.Nil(t, table.Description, "Description should be nil")
	})

	t.Run("Table with nil relations", func(t *testing.T) {
		// Arrange
		table := &domain.Table{
			ID:        5,
			BaseID:    40,
			Name:      "Table with Nil Relations",
			CreatedAt: time.Now(),
			Fields:    nil, // Nil fields
			Records:   nil, // Nil records
			Base:      nil, // Nil base
		}

		// Act & Assert
		assert.NotNil(t, table, "Table should not be nil")
		assert.Nil(t, table.Fields, "Fields can be nil")
		assert.Nil(t, table.Records, "Records can be nil")
		assert.Nil(t, table.Base, "Base can be nil")
	})

	t.Run("Table with very long name", func(t *testing.T) {
		// Arrange
		longName := "VeryLongTableNameThatExceedsNormalLengthButShouldStillWorkAccordingToGORMConstraints"
		table := &domain.Table{
			ID:        6,
			BaseID:    50,
			Name:      longName,
			CreatedAt: time.Now(),
		}

		// Act & Assert
		assert.NotNil(t, table, "Table should not be nil")
		assert.Equal(t, longName, table.Name, "Long name should be supported")
	})
}

// ГРАНИЧНЫЕ ТЕСТЫ - тестирование крайних значений
func TestTableStruct_Boundary(t *testing.T) {
	t.Run("Table with maximum int64 values", func(t *testing.T) {
		// Arrange
		maxDescription := "Maximum values test description"
		
		table := &domain.Table{
			ID:          9223372036854775807,        // math.MaxInt64
			BaseID:      9223372036854775807,       // math.MaxInt64
			Name:        "Max Values Table",
			Description: &maxDescription,
			CreatedAt:   time.Now(),
		}

		// Act & Assert
		assert.NotNil(t, table, "Table should not be nil")
		assert.Equal(t, int64(9223372036854775807), table.ID, "Max ID should be supported")
		assert.Equal(t, int64(9223372036854775807), table.BaseID, "Max BaseID should be supported")
	})

	t.Run("Table with minimum int64 values", func(t *testing.T) {
		// Arrange
		table := &domain.Table{
			ID:        -9223372036854775808,       // math.MinInt64
			BaseID:    -9223372036854775808,      // math.MinInt64
			Name:      "Min Values Table",
			CreatedAt: time.Time{},
		}

		// Act & Assert
		assert.NotNil(t, table, "Table should not be nil")
		assert.Equal(t, int64(-9223372036854775808), table.ID, "Min ID should be supported")
		assert.Equal(t, int64(-9223372036854775808), table.BaseID, "Min BaseID should be supported")
	})

	t.Run("Table with many relations", func(t *testing.T) {
		// Arrange
		var fields []*domain.Field
		var records []*domain.Record
		
		// Создаем много полей и записей
		for i := 0; i < 10; i++ {
			fields = append(fields, &domain.Field{ID: int64(i), Name: "Field " + string(rune(i+48))})
			records = append(records, &domain.Record{ID: int64(i), TableID: 7})
		}

		table := &domain.Table{
			ID:        7,
			BaseID:    60,
			Name:      "Table with Many Relations",
			CreatedAt: time.Now(),
			Fields:    fields,
			Records:   records,
		}

		// Act & Assert
		assert.NotNil(t, table, "Table should not be nil")
		assert.Equal(t, 10, len(table.Fields), "Should have 10 fields")
		assert.Equal(t, 10, len(table.Records), "Should have 10 records")
	})
}

// ТЕСТЫ ДЛЯ ОСОБЫХ СЛУЧАЕВ
func TestTableStruct_SpecialCases(t *testing.T) {
	t.Run("Table with empty name (should be allowed by struct)", func(t *testing.T) {
		// Arrange
		table := &domain.Table{
			ID:        8,
			BaseID:    70,
			Name:      "", // Empty name
			CreatedAt: time.Now(),
		}

		// Act & Assert
		assert.NotNil(t, table, "Table should not be nil")
		assert.Equal(t, "", table.Name, "Empty name should be allowed at struct level")
	})

	t.Run("Table with same ID and BaseID", func(t *testing.T) {
		// Arrange
		description := "Same ID test"
		
		table := &domain.Table{
			ID:          100,
			BaseID:      100, // Same as ID
			Name:        "Same ID Table",
			Description: &description,
			CreatedAt:   time.Now(),
		}

		// Act & Assert
		assert.NotNil(t, table, "Table should not be nil")
		assert.Equal(t, table.ID, table.BaseID, "ID and BaseID can be equal")
	})

	t.Run("Table with special characters in name", func(t *testing.T) {
		// Arrange
		specialName := "Table_With-Special.Characters@123"
		
		table := &domain.Table{
			ID:        9,
			BaseID:    80,
			Name:      specialName,
			CreatedAt: time.Now(),
		}

		// Act & Assert
		assert.NotNil(t, table, "Table should not be nil")
		assert.Equal(t, specialName, table.Name, "Special characters should be supported")
	})
}