package httpclient

import (
	"encoding/json"
	"github.com/hashicorp/go-retryablehttp"
	"io"
	"net/http"
	"sync"
	"time"
)

var (
	clientInstance *retryablehttp.Client
	once           sync.Once
)

type Response struct {
	data   map[string]bool `json:"response_data"`
	status []string        `json:"status"`
}

// func init() {
// }

func GetHttpClient() *retryablehttp.Client {
	once.Do(func() {
		clientInstance = retryablehttp.NewClient()
		clientInstance.RetryMax = 2
		clientInstance.RetryWaitMin = 1 * time.Second
		clientInstance.RetryWaitMax = 3 * time.Second
	})
	return clientInstance
}

func SentRequest(url string, params string) ([]byte, error) {
	if params == "" {
		return nil, nil
	}

	client := GetHttpClient()
	req, err := retryablehttp.NewRequest("GET", url+params, nil)
	if err != nil {
		log.Printf("Failed to prepare http request", slog.Any("error", err))
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send http request", slog.Any("error", err))
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body", slog.Any("error", err))
		return nil, err
	}

	return body, nil
}

func DecodeResponse(respBody []byte) (Response, error) {
	response := Response{}
	err := json.Unmarshal(respBody, &response)

	if err != nil {
		log.Printf("Failed to parse JSON", slog.Any("error", err))
		return response, err
	}

	return response, nil
}