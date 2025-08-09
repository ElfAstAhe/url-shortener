package utils

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeURIStr_data_shouldReturnValue(t *testing.T) {
	expected := "9fae3719b30b1794910a21beefc7f375"
	data := "http://localhost/test/data"
	t.Run("with data", func(t *testing.T) {
		actual := EncodeURIStr(data)

		assert.Equal(t, expected, actual)
	})
}

func TestEncodeURIStr_empty_shouldReturnEmpty(t *testing.T) {
	expected := ""
	data := ""
	t.Run("with empty", func(t *testing.T) {
		actual := EncodeURIStr(data)

		assert.Equal(t, expected, actual)
	})
}

func TestEncodeURI_empty_shouldReturnNil(t *testing.T) {
	var data = make([]byte, 0)
	t.Run("with nil", func(t *testing.T) {
		actual := EncodeURI(data)

		assert.Nil(t, actual)
	})
}

func TestEncodeURI_data_shouldReturn(t *testing.T) {
	data := []byte("http://localhost/test/data")
	expected, _ := hex.DecodeString("9fae3719b30b1794910a21beefc7f375")
	t.Run("with data", func(t *testing.T) {
		actual := EncodeURI(data)

		assert.Equal(t, expected, actual)
	})
}
