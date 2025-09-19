package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ElfAstAhe/url-shortener/internal/config"
	"github.com/ElfAstAhe/url-shortener/internal/handler/dto"
	"github.com/ElfAstAhe/url-shortener/internal/service/auth"
	"github.com/ElfAstAhe/url-shortener/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestShortenPostHandler_DataCorrect_ShouldSuccess(t *testing.T) {
	// prepare
	if config.AppConfig == nil {
		config.AppConfig = config.NewConfig()
		config.AppConfig.LoadConfig()
	}
	router := NewRouter(config.AppConfig, zap.NewNop().Sugar())
	income := dto.ShortenCreateRequest{
		URL: "http://localhost/test/data/post",
	}
	userInfo := auth.BuildRandomUserInfo()
	jwtString, err := auth.NewJWTStringFromUserInfo(userInfo)
	require.NoError(t, err)
	incomeJSON, _ := json.Marshal(income)
	t.Run("", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(incomeJSON))
		req.AddCookie(&http.Cookie{
			Name:  auth.CookieName,
			Value: jwtString,
		})
		recorder := httptest.NewRecorder()
		router.GetRouter().ServeHTTP(recorder, req)

		// assert
		assert.Equal(t, http.StatusCreated, recorder.Code)
	})
}

func TestShortenPostHandler_DataIncorrect_ShouldFail(t *testing.T) {
	// prepare
	if config.AppConfig == nil {
		config.AppConfig = config.NewConfig()
		config.AppConfig.LoadConfig()
	}
	router := NewRouter(config.AppConfig, zap.NewNop().Sugar())
	income := dto.ShortenCreateRequest{
		URL: "",
	}
	userInfo := auth.BuildRandomUserInfo()
	jwtString, err := auth.NewJWTStringFromUserInfo(userInfo)
	require.NoError(t, err)
	incomeJSON, _ := json.Marshal(income)
	testCases := []test.HTTPTestCase{
		{Name: "POST bad data 1", Method: "POST", Path: "/api/shorten", Body: incomeJSON, ExpectedStatusCode: 400},
		{Name: "POST bad data 2", Method: "POST", Path: "/api/shorten", Body: []byte(""), ExpectedStatusCode: 500},
	}

	// act
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			buffer := bytes.NewBuffer(tc.Body)
			req := httptest.NewRequest(tc.Method, tc.Path, buffer)
			req.AddCookie(&http.Cookie{
				Name:  auth.CookieName,
				Value: jwtString,
			})
			recorder := httptest.NewRecorder()
			router.GetRouter().ServeHTTP(recorder, req)

			// assert
			assert.Equal(t, tc.ExpectedStatusCode, recorder.Code)
		})
	}
}
