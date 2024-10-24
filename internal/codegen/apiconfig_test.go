package codegen

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAPIConfig_CreateConfigEntity_CreatesType(t *testing.T) {
	config := APIConfig{
		BaseURL: "baseURLValue",
	}

	entity, err := config.CreateConfigEntity()
	assert.NoError(t, err)

	assert.Equal(t, "APIConfig", entity.Name)

	assert.Contains(t, entity.Properties, "baseURL")
	assert.Equal(t, TypeID_STRING, entity.Properties["baseURL"].Type.TypeID)
}
