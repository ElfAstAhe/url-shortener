package service

import (
	"context"
	"testing"

	"github.com/ElfAstAhe/url-shortener/internal/model"
	"github.com/ElfAstAhe/url-shortener/internal/service/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const ExpectedOriginalURL = "http://localhost:8080/test/data"
const ExpectedKey = "8fe59a11923ca3ea1b7118818e3a7b3c"
const ExpectedID = "123"

var testRoles auth.Roles = auth.Roles{
	"userRole1", "userRole2", "userRole3",
}
var testAdminRoles auth.Roles = auth.Roles{
	"adminRole1", "adminRole2", "adminRole3", "userRole1",
}

type repoMock struct {
}

func (rm *repoMock) GetByKeyUser(ctx context.Context, userID string, key string) (*model.ShortURI, error) {
	if key == "" {
		return nil, nil
	}
	if userID == "" {
		return nil, nil
	}
	data := rm.buildModel()
	if data.Key == key {
		return data, nil
	}
	return nil, nil
}

func (rm *repoMock) ListAllByUser(ctx context.Context, userID string) ([]*model.ShortURI, error) {
	//TODO implement me
	panic("implement me")
}

func (rm *repoMock) BatchCreate(ctx context.Context, batch map[string]*model.ShortURI) (map[string]*model.ShortURI, error) {
	return map[string]*model.ShortURI{}, nil
}

func (rm *repoMock) Get(ctx context.Context, id string) (*model.ShortURI, error) {
	data := rm.buildModel()
	if data.ID == id {
		return data, nil
	}

	return nil, nil
}

func (rm *repoMock) GetByKey(ctx context.Context, key string) (*model.ShortURI, error) {
	data := rm.buildModel()
	if data.Key == key {
		return data, nil
	}
	return nil, nil
}

func (rm *repoMock) Create(ctx context.Context, shortURI *model.ShortURI) (*model.ShortURI, error) {
	if shortURI == nil {
		return nil, nil
	}

	return model.NewShortURI(shortURI.OriginalURL.URL.String(), shortURI.Key)
}

func (rm *repoMock) buildModel() *model.ShortURI {
	data, _ := model.NewShortURIFull(ExpectedID, ExpectedOriginalURL, ExpectedKey)

	return data
}

func TestShorterService_store_shouldReturnKey(t *testing.T) {
	t.Run("should return key", func(t *testing.T) {
		service, err := NewShorterService(&repoMock{})
		require.NoError(t, err)
		actual, err := service.Store(context.Background(), ExpectedOriginalURL)

		assert.NoError(t, err)
		assert.Equal(t, ExpectedKey, actual)
	})
}

func TestNewShorterService_getDataExists_shouldReturnURL(t *testing.T) {
	t.Run("should return URL", func(t *testing.T) {
		repo := &repoMock{}
		ctx := context.WithValue(context.Background(), auth.ContextUserInfo, auth.BuildUnknownUserInfo())
		service, err := NewShorterService(repo)
		require.NoError(t, err)
		actual, err := service.GetURL(ctx, ExpectedKey)

		assert.NoError(t, err)
		assert.Equal(t, ExpectedOriginalURL, actual)
	})

}

func TestNewShorterService_getDataNotExists_shouldReturnEmpty(t *testing.T) {
	t.Run("should return empty", func(t *testing.T) {
		repo := &repoMock{}
		ctx := context.WithValue(context.Background(), auth.ContextUserInfo, auth.BuildUnknownUserInfo())
		service, err := NewShorterService(repo)
		require.NoError(t, err)
		actual, err := service.GetURL(ctx, "22")

		assert.NoError(t, err)
		assert.Equal(t, "", actual)
	})
}
