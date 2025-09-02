package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ElfAstAhe/url-shortener/internal/config"
	"github.com/ElfAstAhe/url-shortener/internal/handler/dto"
	"github.com/ElfAstAhe/url-shortener/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestShortenPostHandler_DataCorrect_ShouldSuccess(t *testing.T) {
	// prepare
	if config.AppConfig == nil {
		config.AppConfig = config.NewConfig()
		config.AppConfig.LoadConfig()
	}
	router := BuildRouter()
	income := dto.ShortenCreateRequest{
		URL: "http://localhost/test/data",
	}
	incomeJSON, _ := json.Marshal(income)
	t.Run("", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(incomeJSON))
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

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
	router := BuildRouter()
	income := dto.ShortenCreateRequest{
		URL: "",
	}
	incomeJSON, _ := json.Marshal(income)
	testCases := []test.HTTPTestCase{
		{Name: "POST bad data 1", Method: "POST", Path: "/api/shorten", Body: incomeJSON, ExpectedStatusCode: 400},
		{Name: "POST bad data 2", Method: "POST", Path: "/api/shorten", Body: []byte(""), ExpectedStatusCode: 500},
	}

	// act
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			buffer := bytes.NewBuffer(tc.Body)
			req, err := http.NewRequest(tc.Method, tc.Path, buffer)
			if err != nil {
				t.Fatal(err)
			}
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			// assert
			assert.Equal(t, tc.ExpectedStatusCode, recorder.Code)
		})
	}
}
