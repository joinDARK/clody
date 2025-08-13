package domain_test

import (
	"testing"

	"clody/core/domain"

	"github.com/stretchr/testify/assert"
)

// ПОЗИТИВНЫЕ ТЕСТЫ - тестирование нормальных сценариев
func TestRecordRelationStruct_Positive(t *testing.T) {
	t.Run("RecordRelation struct creation with valid data", func(t *testing.T) {
		// Arrange
		recordRelation := &domain.RecordRelation{
			ID:             1,
			FieldID:        10,
			SourceRecordID: 20,
			TargetRecordID: 30,
		}

		// Act & Assert
		assert.NotNil(t, recordRelation, "RecordRelation struct should not be nil")
		assert.Equal(t, int64(1), recordRelation.ID, "ID should be 1")
		assert.Equal(t, int64(10), recordRelation.FieldID, "FieldID should be 10")
		assert.Equal(t, int64(20), recordRelation.SourceRecordID, "SourceRecordID should be 20")
		assert.Equal(t, int64(30), recordRelation.TargetRecordID, "TargetRecordID should be 30")
	})

	t.Run("RecordRelation with relations", func(t *testing.T) {
		// Arrange
		field := &domain.Field{ID: 100, Name: "Relation Field"}
		sourceRecord := &domain.Record{ID: 200, TableID: 1}
		targetRecord := &domain.Record{ID: 300, TableID: 2}

		recordRelation := &domain.RecordRelation{
			ID:             2,
			FieldID:        100,
			SourceRecordID: 200,
			TargetRecordID: 300,
			Field:          field,
			SourceRecord:   sourceRecord,
			TargetRecord:   targetRecord,
		}

		// Act & Assert
		assert.NotNil(t, recordRelation, "RecordRelation should not be nil")
		assert.NotNil(t, recordRelation.Field, "Field relation should not be nil")
		assert.NotNil(t, recordRelation.SourceRecord, "SourceRecord relation should not be nil")
		assert.NotNil(t, recordRelation.TargetRecord, "TargetRecord relation should not be nil")

		// Проверка связей
		assert.Equal(t, int64(100), recordRelation.Field.ID, "Field ID should match")
		assert.Equal(t, "Relation Field", recordRelation.Field.Name, "Field name should match")
		assert.Equal(t, int64(200), recordRelation.SourceRecord.ID, "SourceRecord ID should match")
		assert.Equal(t, int64(300), recordRelation.TargetRecord.ID, "TargetRecord ID should match")
	})

	t.Run("RecordRelation with same SourceRecordID and TargetRecordID", func(t *testing.T) {
		// Arrange
		recordRelation := &domain.RecordRelation{
			ID:             3,
			FieldID:        15,
			SourceRecordID: 100,
			TargetRecordID: 100, // Same as SourceRecordID
		}

		// Act & Assert
		assert.NotNil(t, recordRelation, "RecordRelation should not be nil")
		assert.Equal(t, recordRelation.SourceRecordID, recordRelation.TargetRecordID,
			"SourceRecordID and TargetRecordID can be equal")
	})
}

// НЕГАТИВНЫЕ ТЕСТЫ - тестирование граничных и "плохих" сценариев
func TestRecordRelationStruct_Negative(t *testing.T) {
	t.Run("RecordRelation with zero values", func(t *testing.T) {
		// Arrange
		recordRelation := &domain.RecordRelation{
			ID:             0,
			FieldID:        0,
			SourceRecordID: 0,
			TargetRecordID: 0,
		}

		// Act & Assert
		assert.NotNil(t, recordRelation, "RecordRelation should not be nil")
		assert.Equal(t, int64(0), recordRelation.ID, "ID can be zero")
		assert.Equal(t, int64(0), recordRelation.FieldID, "FieldID can be zero")
		assert.Equal(t, int64(0), recordRelation.SourceRecordID, "SourceRecordID can be zero")
		assert.Equal(t, int64(0), recordRelation.TargetRecordID, "TargetRecordID can be zero")
	})

	t.Run("RecordRelation with nil relations", func(t *testing.T) {
		// Arrange
		recordRelation := &domain.RecordRelation{
			ID:             4,
			FieldID:        20,
			SourceRecordID: 30,
			TargetRecordID: 40,
			Field:          nil, // Nil field
			SourceRecord:   nil, // Nil source record
			TargetRecord:   nil, // Nil target record
		}

		// Act & Assert
		assert.NotNil(t, recordRelation, "RecordRelation should not be nil")
		assert.Nil(t, recordRelation.Field, "Field can be nil")
		assert.Nil(t, recordRelation.SourceRecord, "SourceRecord can be nil")
		assert.Nil(t, recordRelation.TargetRecord, "TargetRecord can be nil")
	})

	t.Run("RecordRelation with negative IDs", func(t *testing.T) {
		// Arrange
		recordRelation := &domain.RecordRelation{
			ID:             -1,
			FieldID:        -10,
			SourceRecordID: -20,
			TargetRecordID: -30,
		}

		// Act & Assert
		assert.NotNil(t, recordRelation, "RecordRelation should not be nil")
		assert.Equal(t, int64(-1), recordRelation.ID, "Negative ID should be supported")
		assert.Equal(t, int64(-10), recordRelation.FieldID, "Negative FieldID should be supported")
		assert.Equal(t, int64(-20), recordRelation.SourceRecordID, "Negative SourceRecordID should be supported")
		assert.Equal(t, int64(-30), recordRelation.TargetRecordID, "Negative TargetRecordID should be supported")
	})
}

// ГРАНИЧНЫЕ ТЕСТЫ - тестирование крайних значений
func TestRecordRelationStruct_Boundary(t *testing.T) {
	t.Run("RecordRelation with maximum int64 values", func(t *testing.T) {
		// Arrange
		recordRelation := &domain.RecordRelation{
			ID:             9223372036854775807, // math.MaxInt64
			FieldID:        9223372036854775807, // math.MaxInt64
			SourceRecordID: 9223372036854775807, // math.MaxInt64
			TargetRecordID: 9223372036854775807, // math.MaxInt64
		}

		// Act & Assert
		assert.NotNil(t, recordRelation, "RecordRelation should not be nil")
		assert.Equal(t, int64(9223372036854775807), recordRelation.ID, "Max ID should be supported")
		assert.Equal(t, int64(9223372036854775807), recordRelation.FieldID, "Max FieldID should be supported")
		assert.Equal(t, int64(9223372036854775807), recordRelation.SourceRecordID, "Max SourceRecordID should be supported")
		assert.Equal(t, int64(9223372036854775807), recordRelation.TargetRecordID, "Max TargetRecordID should be supported")
	})

	t.Run("RecordRelation with minimum int64 values", func(t *testing.T) {
		// Arrange
		recordRelation := &domain.RecordRelation{
			ID:             -9223372036854775808, // math.MinInt64
			FieldID:        -9223372036854775808, // math.MinInt64
			SourceRecordID: -9223372036854775808, // math.MinInt64
			TargetRecordID: -9223372036854775808, // math.MinInt64
		}

		// Act & Assert
		assert.NotNil(t, recordRelation, "RecordRelation should not be nil")
		assert.Equal(t, int64(-9223372036854775808), recordRelation.ID, "Min ID should be supported")
		assert.Equal(t, int64(-9223372036854775808), recordRelation.FieldID, "Min FieldID should be supported")
		assert.Equal(t, int64(-9223372036854775808), recordRelation.SourceRecordID, "Min SourceRecordID should be supported")
		assert.Equal(t, int64(-9223372036854775808), recordRelation.TargetRecordID, "Min TargetRecordID should be supported")
	})

	t.Run("RecordRelation with zero FieldID but valid record IDs", func(t *testing.T) {
		// Arrange
		recordRelation := &domain.RecordRelation{
			ID:             5,
			FieldID:        0,   // Zero FieldID
			SourceRecordID: 100, // Valid SourceRecordID
			TargetRecordID: 200, // Valid TargetRecordID
		}

		// Act & Assert
		assert.NotNil(t, recordRelation, "RecordRelation should not be nil")
		assert.Equal(t, int64(0), recordRelation.FieldID, "FieldID can be zero")
		assert.Equal(t, int64(100), recordRelation.SourceRecordID, "SourceRecordID should be 100")
		assert.Equal(t, int64(200), recordRelation.TargetRecordID, "TargetRecordID should be 200")
	})
}

// ТЕСТЫ ДЛЯ ОСОБЫХ СЛУЧАЕВ
func TestRecordRelationStruct_SpecialCases(t *testing.T) {
	t.Run("RecordRelation with same IDs for all fields", func(t *testing.T) {
		// Arrange
		sameID := int64(100)

		recordRelation := &domain.RecordRelation{
			ID:             sameID,
			FieldID:        sameID,
			SourceRecordID: sameID,
			TargetRecordID: sameID,
		}

		// Act & Assert
		assert.NotNil(t, recordRelation, "RecordRelation should not be nil")
		assert.Equal(t, sameID, recordRelation.ID, "ID should match")
		assert.Equal(t, sameID, recordRelation.FieldID, "FieldID should match")
		assert.Equal(t, sameID, recordRelation.SourceRecordID, "SourceRecordID should match")
		assert.Equal(t, sameID, recordRelation.TargetRecordID, "TargetRecordID should match")
	})

	t.Run("RecordRelation with FieldID equal to SourceRecordID", func(t *testing.T) {
		// Arrange
		recordRelation := &domain.RecordRelation{
			ID:             6,
			FieldID:        200, // Same as SourceRecordID
			SourceRecordID: 200, // Same as FieldID
			TargetRecordID: 300,
		}

		// Act & Assert
		assert.NotNil(t, recordRelation, "RecordRelation should not be nil")
		assert.Equal(t, recordRelation.FieldID, recordRelation.SourceRecordID,
			"FieldID and SourceRecordID can be equal")
	})

	t.Run("RecordRelation with very large ID differences", func(t *testing.T) {
		// Arrange
		recordRelation := &domain.RecordRelation{
			ID:             1,
			FieldID:        1000000,
			SourceRecordID: 2000000,
			TargetRecordID: 3000000,
		}

		// Act & Assert
		assert.NotNil(t, recordRelation, "RecordRelation should not be nil")
		assert.Equal(t, int64(1000000), recordRelation.FieldID, "Large FieldID should be supported")
		assert.Equal(t, int64(2000000), recordRelation.SourceRecordID, "Large SourceRecordID should be supported")
		assert.Equal(t, int64(3000000), recordRelation.TargetRecordID, "Large TargetRecordID should be supported")
	})
}

// ТЕСТЫ СВЯЗЕЙ И ОГРАНИЧЕНИЙ
func TestRecordRelationStruct_Relations(t *testing.T) {
	t.Run("RecordRelation with partial relations", func(t *testing.T) {
		// Arrange
		field := &domain.Field{ID: 300, Name: "Partial Field"}
		sourceRecord := &domain.Record{ID: 400, TableID: 3}

		recordRelation := &domain.RecordRelation{
			ID:             7,
			FieldID:        300,
			SourceRecordID: 400,
			TargetRecordID: 500,
			Field:          field,        // Only Field relation
			SourceRecord:   sourceRecord, // Only SourceRecord relation
			TargetRecord:   nil,          // No TargetRecord relation
		}

		// Act & Assert
		assert.NotNil(t, recordRelation, "RecordRelation should not be nil")
		assert.NotNil(t, recordRelation.Field, "Field should not be nil")
		assert.NotNil(t, recordRelation.SourceRecord, "SourceRecord should not be nil")
		assert.Nil(t, recordRelation.TargetRecord, "TargetRecord can be nil")

		assert.Equal(t, "Partial Field", recordRelation.Field.Name, "Field name should match")
		assert.Equal(t, int64(400), recordRelation.SourceRecord.ID, "SourceRecord ID should match")
	})

	t.Run("RecordRelation with no relations", func(t *testing.T) {
		// Arrange
		recordRelation := &domain.RecordRelation{
			ID:             8,
			FieldID:        350,
			SourceRecordID: 450,
			TargetRecordID: 550,
			Field:          nil,
			SourceRecord:   nil,
			TargetRecord:   nil,
		}

		// Act & Assert
		assert.NotNil(t, recordRelation, "RecordRelation should not be nil")
		assert.Nil(t, recordRelation.Field, "Field should be nil")
		assert.Nil(t, recordRelation.SourceRecord, "SourceRecord should be nil")
		assert.Nil(t, recordRelation.TargetRecord, "TargetRecord should be nil")
	})
}
