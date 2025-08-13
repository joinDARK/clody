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
func TestFieldStruct_Positive(t *testing.T) {
	t.Run("Field struct creation with valid data", func(t *testing.T) {
		// Arrange
		description := "Test field description"
		options := datatypes.JSON(`{"required": true, "maxLength": 100}`)
		
		field := &domain.Field{
			ID:          1,
			TableID:     10,
			Name:        "Test Field",
			Type:        domain.FieldTypeText,
			Position:    1,
			Description: &description,
			Options:     options,
			CreatedAt:   time.Now(),
		}

		// Act & Assert
		assert.NotNil(t, field, "Field struct should not be nil")
		assert.Equal(t, int64(1), field.ID, "ID should be 1")
		assert.Equal(t, int64(10), field.TableID, "TableID should be 10")
		assert.Equal(t, "Test Field", field.Name, "Name should be 'Test Field'")
		assert.Equal(t, domain.FieldTypeText, field.Type, "Type should be text")
		assert.Equal(t, 1, field.Position, "Position should be 1")
		assert.Equal(t, description, *field.Description, "Description should match")
		assert.Equal(t, options, field.Options, "Options should match")
		assert.NotZero(t, field.CreatedAt, "CreatedAt should not be zero")
	})

	t.Run("Field with all field types", func(t *testing.T) {
		// Arrange
		fieldTypes := []domain.FieldType{
			domain.FieldTypeText,
			domain.FieldTypeNumber,
			domain.FieldTypeCheckbox,
			domain.FieldTypeDate,
			domain.FieldTypeSingleSelect,
			domain.FieldTypeMultiSelect,
			domain.FieldTypeMany2One,
			domain.FieldTypeOne2Many,
			domain.FieldTypeMany2Many,
			domain.FieldTypeFormula,
		}

		for i, fieldType := range fieldTypes {
			t.Run(string(fieldType), func(t *testing.T) {
				field := &domain.Field{
					ID:        int64(i + 1),
					TableID:   20,
					Name:      string(fieldType) + " Field",
					Type:      fieldType,
					Position:  i,
					CreatedAt: time.Now(),
				}

				assert.NotNil(t, field, "Field should not be nil")
				assert.Equal(t, fieldType, field.Type, "Field type should match")
			})
		}
	})

	t.Run("Field with relations", func(t *testing.T) {
		// Arrange
		table := &domain.Table{ID: 100, Name: "Test Table"}
		selectOption1 := &domain.SelectOption{ID: 1, Label: "Option 1"}
		selectOption2 := &domain.SelectOption{ID: 2, Label: "Option 2"}
		options := utils.CreateJSONFromString(t, `{"minValue": 0, "maxValue": 100}`)

		field := &domain.Field{
			ID:            3,
			TableID:       100,
			Name:          "Field with Relations",
			Type:          domain.FieldTypeSingleSelect,
			Position:      2,
			Options:       options,
			CreatedAt:     time.Now(),
			Table:         table,
			SelectOptions: []*domain.SelectOption{selectOption1, selectOption2},
		}

		// Act & Assert
		assert.NotNil(t, field, "Field should not be nil")
		assert.NotNil(t, field.Table, "Table relation should not be nil")
		assert.Equal(t, 2, len(field.SelectOptions), "Should have 2 select options")
		assert.Equal(t, int64(100), field.Table.ID, "Table ID should match")
		assert.Equal(t, "Test Table", field.Table.Name, "Table name should match")
		assert.Equal(t, int64(1), field.SelectOptions[0].ID, "First option ID should be 1")
		assert.Equal(t, "Option 1", field.SelectOptions[0].Label, "First option name should match")
		assert.Equal(t, int64(2), field.SelectOptions[1].ID, "Second option ID should be 2")
		assert.Equal(t, "Option 2", field.SelectOptions[1].Label, "Second option name should match")
	})

	t.Run("Field with empty slices and nil options", func(t *testing.T) {
		// Arrange
		field := &domain.Field{
			ID:            4,
			TableID:       30,
			Name:          "Field with Empty Relations",
			Type:          domain.FieldTypeText,
			Position:      3,
			Options:       datatypes.JSON{}, // Empty JSON
			CreatedAt:     time.Now(),
			SelectOptions: []*domain.SelectOption{}, // Empty slice
		}

		// Act & Assert
		assert.NotNil(t, field, "Field should not be nil")
		assert.Empty(t, field.SelectOptions, "SelectOptions should be empty")
		assert.NotNil(t, field.SelectOptions, "SelectOptions slice should not be nil")
	})
}

// НЕГАТИВНЫЕ ТЕСТЫ - тестирование граничных и "плохих" сценариев
func TestFieldStruct_Negative(t *testing.T) {
	t.Run("Field with zero values", func(t *testing.T) {
		// Arrange
		field := &domain.Field{
			ID:        0,
			TableID:   0,
			Name:      "",
			Type:      "", // Invalid field type
			Position:  0,
			CreatedAt: time.Time{}, // Zero time
		}

		// Act & Assert
		assert.NotNil(t, field, "Field should not be nil")
		assert.Equal(t, int64(0), field.ID, "ID can be zero")
		assert.Equal(t, int64(0), field.TableID, "TableID can be zero")
		assert.Equal(t, "", field.Name, "Name can be empty")
		assert.Equal(t, domain.FieldType(""), field.Type, "Type can be empty")
		assert.Equal(t, time.Time{}, field.CreatedAt, "CreatedAt can be zero time")
	})

	t.Run("Field with nil description and options", func(t *testing.T) {
		// Arrange
		field := &domain.Field{
			ID:          5,
			TableID:     40,
			Name:        "Field with Nil Values",
			Type:        domain.FieldTypeText,
			Position:    4,
			Description: nil, // Nil description
			Options:     datatypes.JSON{}, // Zero JSON
			CreatedAt:   time.Now(),
		}

		// Act & Assert
		assert.NotNil(t, field, "Field should not be nil")
		assert.Nil(t, field.Description, "Description should be nil")
	})

	t.Run("Field with nil relations", func(t *testing.T) {
		// Arrange
		field := &domain.Field{
			ID:            6,
			TableID:       50,
			Name:          "Field with Nil Relations",
			Type:          domain.FieldTypeNumber,
			Position:      5,
			CreatedAt:     time.Now(),
			Table:         nil, // Nil table
			SelectOptions: nil, // Nil select options
		}

		// Act & Assert
		assert.NotNil(t, field, "Field should not be nil")
		assert.Nil(t, field.Table, "Table can be nil")
		assert.Nil(t, field.SelectOptions, "SelectOptions can be nil")
	})

	t.Run("Field with negative position", func(t *testing.T) {
		// Arrange
		field := &domain.Field{
			ID:        7,
			TableID:   60,
			Name:      "Field with Negative Position",
			Type:      domain.FieldTypeText,
			Position:  -1, // Negative position
			CreatedAt: time.Now(),
		}

		// Act & Assert
		assert.NotNil(t, field, "Field should not be nil")
		assert.Equal(t, -1, field.Position, "Negative position should be allowed")
	})
}

// ТЕСТЫ ВАЛИДАЦИИ ТИПОВ ПОЛЕЙ
func TestFieldTypeValidation(t *testing.T) {
	t.Run("Valid field types", func(t *testing.T) {
		validTypes := []domain.FieldType{
			domain.FieldTypeText,
			domain.FieldTypeNumber,
			domain.FieldTypeCheckbox,
			domain.FieldTypeDate,
			domain.FieldTypeSingleSelect,
			domain.FieldTypeMultiSelect,
			domain.FieldTypeMany2One,
			domain.FieldTypeOne2Many,
			domain.FieldTypeMany2Many,
			domain.FieldTypeFormula,
		}

		for _, fieldType := range validTypes {
			t.Run(string(fieldType), func(t *testing.T) {
				isValid := fieldType.IsValid()
				assert.True(t, isValid, "Field type %s should be valid", fieldType)
			})
		}
	})

	t.Run("Invalid field types", func(t *testing.T) {
		invalidTypes := []domain.FieldType{
			"",
			"invalid_type",
			"unknown",
			"text_field", // Хотя похож на text, но не точное совпадение
			"number_field",
		}

		for _, fieldType := range invalidTypes {
			t.Run(string(fieldType), func(t *testing.T) {
				isValid := fieldType.IsValid()
				assert.False(t, isValid, "Field type %s should be invalid", fieldType)
			})
		}
	})

	t.Run("Case sensitivity test", func(t *testing.T) {
		// Arrange
		mixedCaseTypes := []struct {
			input    domain.FieldType
			expected bool
		}{
			{domain.FieldType("TEXT"), false},           // Верхний регистр
			{domain.FieldType("Text"), false},           // Смешанный регистр
			{domain.FieldType("text"), true},            // Правильный регистр
			{domain.FieldType("NUMBER"), false},         // Верхний регистр
			{domain.FieldType("number"), true},          // Правильный регистр
		}

		// Act & Assert
		for _, testCase := range mixedCaseTypes {
			result := testCase.input.IsValid()
			assert.Equal(t, testCase.expected, result, 
				"Field type %s validation should be %v", testCase.input, testCase.expected)
		}
	})
}

// ГРАНИЧНЫЕ ТЕСТЫ - тестирование крайних значений
func TestFieldStruct_Boundary(t *testing.T) {
	t.Run("Field with maximum int64 values", func(t *testing.T) {
		// Arrange
		description := "Maximum values test"
		options := utils.CreateJSONFromString(t, `{"maxTest": true}`)
		
		field := &domain.Field{
			ID:          9223372036854775807,        // math.MaxInt64
			TableID:     9223372036854775807,       // math.MaxInt64
			Name:        "Max Values Field",
			Type:        domain.FieldTypeText,
			Position:    2147483647,                // math.MaxInt32
			Description: &description,
			Options:     options,
			CreatedAt:   time.Now(),
		}

		// Act & Assert
		assert.NotNil(t, field, "Field should not be nil")
		assert.Equal(t, int64(9223372036854775807), field.ID, "Max ID should be supported")
		assert.Equal(t, int64(9223372036854775807), field.TableID, "Max TableID should be supported")
		assert.Equal(t, 2147483647, field.Position, "Max position should be supported")
	})

	t.Run("Field with minimum int64 values", func(t *testing.T) {
		// Arrange
		field := &domain.Field{
			ID:        -9223372036854775808,       // math.MinInt64
			TableID:   -9223372036854775808,      // math.MinInt64
			Name:      "Min Values Field",
			Type:      domain.FieldTypeText,
			Position:  -2147483648,               // math.MinInt32
			CreatedAt: time.Time{},
		}

		// Act & Assert
		assert.NotNil(t, field, "Field should not be nil")
		assert.Equal(t, int64(-9223372036854775808), field.ID, "Min ID should be supported")
		assert.Equal(t, int64(-9223372036854775808), field.TableID, "Min TableID should be supported")
		assert.Equal(t, -2147483648, field.Position, "Min position should be supported")
	})

	t.Run("Field with complex JSON options", func(t *testing.T) {
		// Arrange
		complexOptions := `{
			"validation": {
				"required": true,
				"minLength": 1,
				"maxLength": 255
			},
			"ui": {
				"placeholder": "Enter text here",
				"helpText": "This is a helper text"
			},
			"default": "Default value",
			"advanced": {
				"regex": "^[a-zA-Z0-9]+$",
				"transform": "uppercase"
			}
		}`
		
		options := utils.CreateJSONFromString(t, complexOptions)

		field := &domain.Field{
			ID:        8,
			TableID:   70,
			Name:      "Field with Complex Options",
			Type:      domain.FieldTypeText,
			Position:  8,
			Options:   options,
			CreatedAt: time.Now(),
		}

		// Act & Assert
		assert.NotNil(t, field, "Field should not be nil")
		assert.Equal(t, options, field.Options, "Complex JSON options should be stored correctly")
	})
}

// ТЕСТЫ ДЛЯ ОСОБЫХ СЛУЧАЕВ
func TestFieldStruct_SpecialCases(t *testing.T) {
	t.Run("Field with same ID and TableID", func(t *testing.T) {
		// Arrange
		description := "Same ID test"
		
		field := &domain.Field{
			ID:          100,
			TableID:     100, // Same as ID
			Name:        "Same ID Field",
			Type:        domain.FieldTypeText,
			Position:    9,
			Description: &description,
			CreatedAt:   time.Now(),
		}

		// Act & Assert
		assert.NotNil(t, field, "Field should not be nil")
		assert.Equal(t, field.ID, field.TableID, "ID and TableID can be equal")
	})

	t.Run("Field with empty name (should be allowed by struct)", func(t *testing.T) {
		// Arrange
		field := &domain.Field{
			ID:        9,
			TableID:   80,
			Name:      "", // Empty name
			Type:      domain.FieldTypeText,
			Position:  10,
			CreatedAt: time.Now(),
		}

		// Act & Assert
		assert.NotNil(t, field, "Field should not be nil")
		assert.Equal(t, "", field.Name, "Empty name should be allowed at struct level")
	})

	t.Run("Field with special characters in name", func(t *testing.T) {
		// Arrange
		specialName := "Field_With-Special.Characters@123#$%"
		
		field := &domain.Field{
			ID:        10,
			TableID:   90,
			Name:      specialName,
			Type:      domain.FieldTypeText,
			Position:  11,
			CreatedAt: time.Now(),
		}

		// Act & Assert
		assert.NotNil(t, field, "Field should not be nil")
		assert.Equal(t, specialName, field.Name, "Special characters should be supported")
	})

	t.Run("Field with default type from GORM tag", func(t *testing.T) {
		// Arrange
		field := &domain.Field{
			ID:        11,
			TableID:   100,
			Name:      "Field with Default Type",
			// Type не указан, должен быть default из GORM тега
			Position:  12,
			CreatedAt: time.Now(),
		}

		// Act & Assert
		assert.NotNil(t, field, "Field should not be nil")
		// Note: В структуре default тип 'text', но в коде это не устанавливается автоматически
	})
}