package utils

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

// CreateJSONFromString - вспомогательная функция для создания datatypes.JSON
func CreateJSONFromString(t *testing.T, jsonString string) datatypes.JSON {
	var jsonData any
	err := json.Unmarshal([]byte(jsonString), &jsonData)
	assert.NoError(t, err, "Should unmarshal JSON string without error")

	jsonBytes, err := json.Marshal(jsonData)
	assert.NoError(t, err, "Should marshal JSON data without error")

	return datatypes.JSON(jsonBytes)
}
