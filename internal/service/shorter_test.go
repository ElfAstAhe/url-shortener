package service

import (
	"testing"
	"time"

	"github.com/ElfAstAhe/url-shortener/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const ExpectedOriginalURL = "http://localhost:8080/test/data"
const ExpectedKey = "8fe59a11923ca3ea1b7118818e3a7b3c"
const ExpectedID = "123"

type repoMock struct {
}

func (r repoMock) ListAllByUser(userID string) ([]*_model.ShortURI, error) {
	//TODO implement me
	panic("implement me")
}

func (r repoMock) BatchCreate(batch map[string]*model.ShortURI) (map[string]*model.ShortURI, error) {
	return map[string]*model.ShortURI{}, nil
}

func (r repoMock) Get(id string) (*model.ShortURI, error) {
	data := buildModel()
	if data.ID == id {
		return data, nil
	}

	return nil, nil
}

func (r repoMock) GetByKey(key string) (*model.ShortURI, error) {
	data := buildModel()
	if data.Key == key {
		return data, nil
	}
	return nil, nil
}

func (r repoMock) Create(shortURI *model.ShortURI) (*model.ShortURI, error) {
	if shortURI == nil {
		return nil, nil
	}

	return model.NewShortURI(shortURI.OriginalURL.URL.String(), shortURI.Key)
}

func buildModel() *model.ShortURI {
	techData := model.TechData{
		CreateUser: "unknown",
		Created:    time.Now(),
		UpdateUser: "unknown",
		Updated:    time.Now(),
	}
	data, _ := model.NewShortURIFull(ExpectedID, ExpectedOriginalURL, ExpectedKey, &techData)

	return data
}

func TestShorterService_store_shouldReturnKey(t *testing.T) {
	t.Run("should return key", func(t *testing.T) {
		service, err := NewShorterService(&repoMock{})
		require.NoError(t, err)
		actual, err := service.Store(ExpectedOriginalURL)

		assert.NoError(t, err)
		assert.Equal(t, ExpectedKey, actual)
	})
}

func TestNewShorterService_getDataExists_shouldReturnURL(t *testing.T) {
	t.Run("should return URL", func(t *testing.T) {
		service, err := NewShorterService(&repoMock{})
		require.NoError(t, err)
		actual, err := service.GetURL(ExpectedKey)

		assert.NoError(t, err)
		assert.Equal(t, ExpectedOriginalURL, actual)
	})

}

func TestNewShorterService_getDataNotExists_shouldReturnEmpty(t *testing.T) {
	t.Run("should return empty", func(t *testing.T) {
		service, err := NewShorterService(&repoMock{})
		require.NoError(t, err)
		actual, err := service.GetURL("22")

		assert.NoError(t, err)
		assert.Equal(t, "", actual)
	})
}
