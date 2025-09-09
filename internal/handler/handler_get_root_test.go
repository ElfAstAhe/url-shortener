package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ElfAstAhe/url-shortener/internal/config"
	"github.com/ElfAstAhe/url-shortener/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestRootHandler_getMethod_emptyKey_shouldReturnBadRequest(t *testing.T) {
	// prepare
	if config.AppConfig == nil {
		config.AppConfig = config.NewConfig()
		config.AppConfig.LoadConfig()
	}
	router := NewRouter(config.AppConfig, zap.NewNop().Sugar())
	// test cases
	cases := []test.HTTPTestCase{
		{Name: "method GET path /", Method: http.MethodGet, Path: "/", ExpectedStatusCode: http.StatusMethodNotAllowed},
		{Name: "method GET path //", Method: http.MethodGet, Path: "//", ExpectedStatusCode: http.StatusNotFound},
	}
	// act
	for _, testCase := range cases {
		t.Run(testCase.Name, func(t *testing.T) {
			req := httptest.NewRequest(testCase.Method, testCase.Path, nil)
			recorder := httptest.NewRecorder()
			router.GetRouter().ServeHTTP(recorder, req)

			assert.Equal(t, testCase.ExpectedStatusCode, recorder.Code)
		})
	}
}

func TestRootHandler_getMethod_success(t *testing.T) {
	// prepare
	if config.AppConfig == nil {
		config.AppConfig = config.NewConfig()
		config.AppConfig.LoadConfig()
	}
	expectedURL := "http://localhost/test/data"
	router := NewRouter(config.AppConfig, zap.NewNop().Sugar())
	chiRouter, ok := router.(*chiRouter)
	require.True(t, ok)
	var service, _ = chiRouter.createShortenService()
	_, err := service.Store(expectedURL)
	require.NoError(t, err)
	// test cases
	testCases := []test.HTTPTestCase{
		{Name: "method GET -> not found", Method: http.MethodGet, Path: "http://localhost:8080/123", ExpectedStatusCode: http.StatusNotFound, ExpectedStrValue: ""},
		{Name: "method GET -> redirect", Method: http.MethodGet, Path: "http://localhost:8080/9fae3719b30b1794910a21beefc7f375", ExpectedStatusCode: http.StatusTemporaryRedirect, ExpectedStrValue: expectedURL},
	}
	// act
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			req := httptest.NewRequest(testCase.Method, testCase.Path, nil)
			recorder := httptest.NewRecorder()
			router.GetRouter().ServeHTTP(recorder, req)

			// assert
			assert.Equal(t, testCase.ExpectedStatusCode, recorder.Code)
			assert.Equal(t, testCase.ExpectedStrValue, recorder.Header().Get("Location"))
		})
	}

}
