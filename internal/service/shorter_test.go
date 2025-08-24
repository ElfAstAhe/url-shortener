package service

import (
	"testing"

	_model "github.com/ElfAstAhe/url-shortener/internal/model"
	"github.com/stretchr/testify/assert"
)

type repoMock struct {
}

const ExpectedOriginalURL = "http://localhost:8080/test/data"
const ExpectedKey = "8fe59a11923ca3ea1b7118818e3a7b3c"
const ExpectedID = "123"

func (r repoMock) GetByID(id string) (*_model.ShortURI, error) {
	data := buildModel()
	if data.ID == id {
		return data, nil
	}

	return nil, nil
}

func (r repoMock) GetByKey(key string) (*_model.ShortURI, error) {
	data := buildModel()
	if data.Key == key {
		return data, nil
	}
	return nil, nil
}

func (r repoMock) Create(shortURI *_model.ShortURI) (*_model.ShortURI, error) {
	if shortURI == nil {
		return nil, nil
	}

	return _model.NewShortURI(shortURI.OriginalURL.String(), shortURI.Key)
}

func buildModel() *_model.ShortURI {
	data, _ := _model.NewShortURIFull(ExpectedID, ExpectedOriginalURL, ExpectedKey)

	return data
}

func TestShorterService_store_shouldReturnKey(t *testing.T) {
	t.Run("should return key", func(t *testing.T) {
		actual, err := NewShorterService(&repoMock{}).Store(ExpectedOriginalURL)

		assert.NoError(t, err)
		assert.Equal(t, ExpectedKey, actual)
	})
}

func TestNewShorterService_getDataExists_shouldReturnURL(t *testing.T) {
	t.Run("should return URL", func(t *testing.T) {
		actual, err := NewShorterService(&repoMock{}).GetURL(ExpectedKey)

		assert.NoError(t, err)
		assert.Equal(t, ExpectedOriginalURL, actual)
	})

}

func TestNewShorterService_getDataNotExists_shouldReturnEmpty(t *testing.T) {
	t.Run("should return empty", func(t *testing.T) {
		actual, err := NewShorterService(&repoMock{}).GetURL("22")

		assert.NoError(t, err)
		assert.Equal(t, "", actual)
	})
}
