package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"clody/core/domain"
)

// ПОЗИТИВНЫЕ ТЕСТЫ - тестирование нормальных сценариев
func TestSelectOptionStruct_Positive(t *testing.T) {
	t.Run("SelectOption struct creation with valid data", func(t *testing.T) {
		// Arrange
		selectOption := &domain.SelectOption{
			ID:       1,
			FieldID:  10,
			Label:    "Test Option",
			Color:    "#FF0000",
			Position: 1,
		}

		// Act & Assert
		assert.NotNil(t, selectOption, "SelectOption struct should not be nil")
		assert.Equal(t, int64(1), selectOption.ID, "ID should be 1")
		assert.Equal(t, int64(10), selectOption.FieldID, "FieldID should be 10")
		assert.Equal(t, "Test Option", selectOption.Label, "Label should be 'Test Option'")
		assert.Equal(t, "#FF0000", selectOption.Color, "Color should be '#FF0000'")
		assert.Equal(t, 1, selectOption.Position, "Position should be 1")
	})

	t.Run("SelectOption with field relation", func(t *testing.T) {
		// Arrange
		field := &domain.Field{ID: 100, Name: "Test Field"}

		selectOption := &domain.SelectOption{
			ID:       2,
			FieldID:  100,
			Label:    "Option with Field",
			Color:    "#00FF00",
			Position: 2,
			Field:    field,
		}

		// Act & Assert
		assert.NotNil(t, selectOption, "SelectOption should not be nil")
		assert.NotNil(t, selectOption.Field, "Field relation should not be nil")
		assert.Equal(t, int64(100), selectOption.Field.ID, "Field ID should match")
		assert.Equal(t, "Test Field", selectOption.Field.Name, "Field name should match")
	})

	t.Run("SelectOption with various color formats", func(t *testing.T) {
		// Arrange
		colorFormats := []struct {
			name  string
			color string
		}{
			{"Hex color", "#FF5733"},
			{"RGB color", "rgb(255, 87, 51)"},
			{"Named color", "red"},
			{"HSL color", "hsl(20, 100%, 50%)"},
			{"Empty color", ""},
		}

		for _, format := range colorFormats {
			t.Run(format.name, func(t *testing.T) {
				selectOption := &domain.SelectOption{
					ID:       3,
					FieldID:  20,
					Label:    format.name,
					Color:    format.color,
					Position: 3,
				}

				assert.NotNil(t, selectOption, "SelectOption should not be nil")
				assert.Equal(t, format.color, selectOption.Color, "Color should match")
			})
		}
	})
}

// НЕГАТИВНЫЕ ТЕСТЫ - тестирование граничных и "плохих" сценариев
func TestSelectOptionStruct_Negative(t *testing.T) {
	t.Run("SelectOption with zero values", func(t *testing.T) {
		// Arrange
		selectOption := &domain.SelectOption{
			ID:       0,
			FieldID:  0,
			Label:    "",
			Color:    "",
			Position: 0,
		}

		// Act & Assert
		assert.NotNil(t, selectOption, "SelectOption should not be nil")
		assert.Equal(t, int64(0), selectOption.ID, "ID can be zero")
		assert.Equal(t, int64(0), selectOption.FieldID, "FieldID can be zero")
		assert.Equal(t, "", selectOption.Label, "Label can be empty")
		assert.Equal(t, "", selectOption.Color, "Color can be empty")
		assert.Equal(t, 0, selectOption.Position, "Position can be zero")
	})

	t.Run("SelectOption with nil field relation", func(t *testing.T) {
		// Arrange
		selectOption := &domain.SelectOption{
			ID:       4,
			FieldID:  30,
			Label:    "Option with Nil Field",
			Color:    "#0000FF",
			Position: 4,
			Field:    nil, // Nil field
		}

		// Act & Assert
		assert.NotNil(t, selectOption, "SelectOption should not be nil")
		assert.Nil(t, selectOption.Field, "Field can be nil")
	})

	t.Run("SelectOption with negative values", func(t *testing.T) {
		// Arrange
		selectOption := &domain.SelectOption{
			ID:       -1,
			FieldID:  -10,
			Label:    "Negative Values Option",
			Color:    "#FFFF00",
			Position: -1,
		}

		// Act & Assert
		assert.NotNil(t, selectOption, "SelectOption should not be nil")
		assert.Equal(t, int64(-1), selectOption.ID, "Negative ID should be supported")
		assert.Equal(t, int64(-10), selectOption.FieldID, "Negative FieldID should be supported")
		assert.Equal(t, -1, selectOption.Position, "Negative position should be supported")
	})

	t.Run("SelectOption with very long label", func(t *testing.T) {
		// Arrange
		longLabel := "VeryLongLabelThatExceedsNormalLengthButShouldStillWorkAccordingToGORMConstraintsAndDatabaseLimits"
		
		selectOption := &domain.SelectOption{
			ID:       5,
			FieldID:  40,
			Label:    longLabel,
			Color:    "#FF00FF",
			Position: 5,
		}

		// Act & Assert
		assert.NotNil(t, selectOption, "SelectOption should not be nil")
		assert.Equal(t, longLabel, selectOption.Label, "Long label should be supported")
	})
}

// ГРАНИЧНЫЕ ТЕСТЫ - тестирование крайних значений
func TestSelectOptionStruct_Boundary(t *testing.T) {
	t.Run("SelectOption with maximum int64 values", func(t *testing.T) {
		// Arrange
		selectOption := &domain.SelectOption{
			ID:       9223372036854775807,        // math.MaxInt64
			FieldID:  9223372036854775807,       // math.MaxInt64
			Label:    "Max Values Option",
			Color:    "#FFFFFF",
			Position: 2147483647,                // math.MaxInt32
		}

		// Act & Assert
		assert.NotNil(t, selectOption, "SelectOption should not be nil")
		assert.Equal(t, int64(9223372036854775807), selectOption.ID, "Max ID should be supported")
		assert.Equal(t, int64(9223372036854775807), selectOption.FieldID, "Max FieldID should be supported")
		assert.Equal(t, 2147483647, selectOption.Position, "Max position should be supported")
	})

	t.Run("SelectOption with minimum int64 values", func(t *testing.T) {
		// Arrange
		selectOption := &domain.SelectOption{
			ID:       -9223372036854775808,       // math.MinInt64
			FieldID:  -9223372036854775808,      // math.MinInt64
			Label:    "Min Values Option",
			Color:    "#000000",
			Position: -2147483648,               // math.MinInt32
		}

		// Act & Assert
		assert.NotNil(t, selectOption, "SelectOption should not be nil")
		assert.Equal(t, int64(-9223372036854775808), selectOption.ID, "Min ID should be supported")
		assert.Equal(t, int64(-9223372036854775808), selectOption.FieldID, "Min FieldID should be supported")
		assert.Equal(t, -2147483648, selectOption.Position, "Min position should be supported")
	})

	t.Run("SelectOption with empty strings", func(t *testing.T) {
		// Arrange
		selectOption := &domain.SelectOption{
			ID:       6,
			FieldID:  50,
			Label:    "", // Empty label
			Color:    "", // Empty color
			Position: 6,
		}

		// Act & Assert
		assert.NotNil(t, selectOption, "SelectOption should not be nil")
		assert.Equal(t, "", selectOption.Label, "Empty label should be supported")
		assert.Equal(t, "", selectOption.Color, "Empty color should be supported")
	})
}

// ТЕСТЫ ДЛЯ ОСОБЫХ СЛУЧАЕВ
func TestSelectOptionStruct_SpecialCases(t *testing.T) {
	t.Run("SelectOption with same ID and FieldID", func(t *testing.T) {
		// Arrange
		selectOption := &domain.SelectOption{
			ID:       100,
			FieldID:  100, // Same as ID
			Label:    "Same ID Option",
			Color:    "#123456",
			Position: 7,
		}

		// Act & Assert
		assert.NotNil(t, selectOption, "SelectOption should not be nil")
		assert.Equal(t, selectOption.ID, selectOption.FieldID, "ID and FieldID can be equal")
	})

	t.Run("SelectOption with special characters in label", func(t *testing.T) {
		// Arrange
		specialLabel := "Option_With-Special.Characters@123#$%&*()"
		
		selectOption := &domain.SelectOption{
			ID:       7,
			FieldID:  60,
			Label:    specialLabel,
			Color:    "#ABCDEF",
			Position: 8,
		}

		// Act & Assert
		assert.NotNil(t, selectOption, "SelectOption should not be nil")
		assert.Equal(t, specialLabel, selectOption.Label, "Special characters should be supported")
	})

	t.Run("SelectOption with zero position", func(t *testing.T) {
		// Arrange
		selectOption := &domain.SelectOption{
			ID:       8,
			FieldID:  70,
			Label:    "Zero Position Option",
			Color:    "#789012",
			Position: 0, // Zero position
		}

		// Act & Assert
		assert.NotNil(t, selectOption, "SelectOption should not be nil")
		assert.Equal(t, 0, selectOption.Position, "Zero position should be allowed")
	})

	t.Run("SelectOption with negative position", func(t *testing.T) {
		// Arrange
		selectOption := &domain.SelectOption{
			ID:       9,
			FieldID:  80,
			Label:    "Negative Position Option",
			Color:    "#345678",
			Position: -5, // Negative position
		}

		// Act & Assert
		assert.NotNil(t, selectOption, "SelectOption should not be nil")
		assert.Equal(t, -5, selectOption.Position, "Negative position should be supported")
	})
}

// ТЕСТЫ СВЯЗЕЙ И ОГРАНИЧЕНИЙ
func TestSelectOptionStruct_Relations(t *testing.T) {
	t.Run("SelectOption with field relation", func(t *testing.T) {
		// Arrange
		field := &domain.Field{
			ID:   200,
			Name: "Related Field",
			Type: domain.FieldTypeSingleSelect,
		}

		selectOption := &domain.SelectOption{
			ID:       10,
			FieldID:  200,
			Label:    "Related Option",
			Color:    "#901234",
			Position: 10,
			Field:    field,
		}

		// Act & Assert
		assert.NotNil(t, selectOption, "SelectOption should not be nil")
		assert.NotNil(t, selectOption.Field, "Field should not be nil")
		assert.Equal(t, int64(200), selectOption.FieldID, "FieldID should match")
		assert.Equal(t, "Related Field", selectOption.Field.Name, "Field name should match")
		assert.Equal(t, domain.FieldTypeSingleSelect, selectOption.Field.Type, "Field type should match")
	})

	t.Run("SelectOption without field relation", func(t *testing.T) {
		// Arrange
		selectOption := &domain.SelectOption{
			ID:       11,
			FieldID:  90,
			Label:    "Unrelated Option",
			Color:    "#567890",
			Position: 11,
			Field:    nil, // No field relation
		}

		// Act & Assert
		assert.NotNil(t, selectOption, "SelectOption should not be nil")
		assert.Nil(t, selectOption.Field, "Field can be nil")
		assert.Equal(t, int64(90), selectOption.FieldID, "FieldID should still be set")
	})

	t.Run("SelectOption with field ID but no field object", func(t *testing.T) {
		// Arrange
		selectOption := &domain.SelectOption{
			ID:       12,
			FieldID:  300, // Valid FieldID
			Label:    "Field ID Only Option",
			Color:    "#234567",
			Position: 12,
			Field:    nil, // But no field object
		}

		// Act & Assert
		assert.NotNil(t, selectOption, "SelectOption should not be nil")
		assert.Equal(t, int64(300), selectOption.FieldID, "FieldID should be set")
		assert.Nil(t, selectOption.Field, "Field object can be nil")
	})
}