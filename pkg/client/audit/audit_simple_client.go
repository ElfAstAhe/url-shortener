package audit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/ElfAstAhe/url-shortener/internal/utils"
	"github.com/ElfAstAhe/url-shortener/pkg/client/audit/dto"
)

type SimpleClient struct {
	client  *http.Client
	baseURL string
}

func NewSimpleClient(baseURL string, timeOut time.Duration) *SimpleClient {
	return &SimpleClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: timeOut},
	}
}

func (client *SimpleClient) AuditIncome(data *dto.IncomeAuditDto) error {
	if data == nil {
		return nil
	}

	auditURL, err := url.Parse(client.baseURL)
	if err != nil {
		return NewClientError("Error parsing base URL", err)
	}
	buffer, err := json.Marshal(data)
	if err != nil {
		return NewClientError("Error marshalling audit data", err)
	}

	req, err := http.NewRequest(http.MethodPost, auditURL.String(), bytes.NewBuffer(buffer))
	if err != nil {
		return NewClientError("Error creating request", err)
	}

	resp, err := client.client.Do(req)
	if err != nil {
		return NewClientError("Error sending request", err)
	}

	if !utils.IsSuccess(resp.StatusCode) {
		return NewClientError(fmt.Sprintf("Error in response with status code [%v]", resp.StatusCode), nil)
	}

	return nil
}
