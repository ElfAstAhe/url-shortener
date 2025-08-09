package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildNewURI_data_shouldReturnValue(t *testing.T) {
	dataBaseURL := "http://localhost:8080"
	dataKey := "123"
	expected := dataBaseURL + "/" + dataKey

	t.Run("with data", func(t *testing.T) {
		actual := BuildNewURI(dataBaseURL, dataKey)

		assert.Equal(t, expected, actual)
	})
}

func TestBuildNewURI_keyEmpty_shouldReturnEmpty(t *testing.T) {
	dataBaseURL := "http://localhost:8080"
	dataKey := ""
	expected := ""

	t.Run("with empty key", func(t *testing.T) {
		actual := BuildNewURI(dataBaseURL, dataKey)

		assert.Equal(t, expected, actual)
	})
}

func TestBuildNewURI_baseURLEmpty_shouldReturnEmpty(t *testing.T) {
	dataBaseURL := ""
	dataKey := "123"
	expected := ""

	t.Run("with empty baseURL", func(t *testing.T) {
		actual := BuildNewURI(dataBaseURL, dataKey)

		assert.Equal(t, expected, actual)
	})
}
