package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ElfAstAhe/url-shortener/internal/config"
	"github.com/ElfAstAhe/url-shortener/internal/handler/helper"
	"github.com/ElfAstAhe/url-shortener/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestRootHandler_getMethod_emptyKey_shouldReturnBadRequest(t *testing.T) {
	// prepare
	if config.AppConfig == nil {
		config.AppConfig = config.NewConfig()
		config.AppConfig.LoadConfig()
	}
	router := BuildRouter()
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
			router.ServeHTTP(recorder, req)

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
	_, _ = helper.CreateService().Store(expectedURL)
	router := BuildRouter()
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
			router.ServeHTTP(recorder, req)

			// assert
			assert.Equal(t, testCase.ExpectedStatusCode, recorder.Code)
			assert.Equal(t, testCase.ExpectedStrValue, recorder.Header().Get("Location"))
		})
	}

}
