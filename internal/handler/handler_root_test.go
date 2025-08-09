package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	_srv "github.com/ElfAstAhe/url-shortener/internal/service"
	"github.com/ElfAstAhe/url-shortener/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestRootHandler_unacceptableMethods_shouldReturnBadRequest(t *testing.T) {
	cases := []test.HttpTestCase{
		{Name: "http method CONNECT", Method: http.MethodConnect, ExpectedStatusCode: http.StatusBadRequest},
		{Name: "http method HEAD", Method: http.MethodHead, ExpectedStatusCode: http.StatusBadRequest},
		{Name: "http method DELETE", Method: http.MethodDelete, ExpectedStatusCode: http.StatusBadRequest},
		{Name: "http method PATCH", Method: http.MethodPatch, ExpectedStatusCode: http.StatusBadRequest},
		{Name: "http method PUT", Method: http.MethodPut, ExpectedStatusCode: http.StatusBadRequest},
	}

	for _, testCase := range cases {
		t.Run(testCase.Name, func(t *testing.T) {
			req := httptest.NewRequest(testCase.Method, RootHandlePath, nil)
			recorder := httptest.NewRecorder()
			RootHandler(recorder, req)

			assert.Equal(t, http.StatusBadRequest, recorder.Code)
		})
	}
}

func TestRootHandler_getMethod_emptyKey_shouldReturnBadRequest(t *testing.T) {
	cases := []test.HttpTestCase{
		{Name: "method GET path /", Method: http.MethodGet, Path: "/", ExpectedStatusCode: http.StatusBadRequest},
		{Name: "method GET path //", Method: http.MethodGet, Path: "//", ExpectedStatusCode: http.StatusBadRequest},
	}

	for _, testCase := range cases {
		t.Run(testCase.Name, func(t *testing.T) {
			req := httptest.NewRequest(testCase.Method, testCase.Path, nil)
			recorder := httptest.NewRecorder()
			RootHandler(recorder, req)

			assert.Equal(t, testCase.ExpectedStatusCode, recorder.Code)
		})
	}
}

func TestRootHandler_getMethod_success(t *testing.T) {
	// prepare
	expectedURL := "http://localhost/test/data"
	_, _ = _srv.NewShorterService().Store(expectedURL)
	//
	testCases := []test.HttpTestCase{
		{Name: "method GET -> not found", Method: http.MethodGet, Path: "/123", ExpectedStatusCode: http.StatusNotFound, ExpectedStrValue: ""},
		{Name: "method GET -> redirect", Method: http.MethodGet, Path: "/9fae3719b30b1794910a21beefc7f375", ExpectedStatusCode: http.StatusTemporaryRedirect, ExpectedStrValue: expectedURL},
	}
	// action
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			req := httptest.NewRequest(testCase.Method, testCase.Path, nil)
			recorder := httptest.NewRecorder()
			RootHandler(recorder, req)

			// assert
			assert.Equal(t, testCase.ExpectedStatusCode, recorder.Code)
			assert.Equal(t, testCase.ExpectedStrValue, recorder.Header().Get("Location"))
		})
	}

}

func TestRootHandler_getMethod_dataNotExists_shouldReturn404(t *testing.T) {}
