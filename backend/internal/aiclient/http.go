package aiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// HTTPAIClient calls the real AI microservice's POST /ai/v1/predict.
type HTTPAIClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewHTTPAIClient(baseURL string) *HTTPAIClient {
	return &HTTPAIClient{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *HTTPAIClient) Predict(ctx context.Context, req PredictRequest) (PredictResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return PredictResponse{}, fmt.Errorf("marshal request: %w", err) // Handle the 
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/ai/v1/predict", bytes.NewReader(body))
	if err != nil {
		return PredictResponse{}, fmt.Errorf("build request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return PredictResponse{}, fmt.Errorf("ai service unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return PredictResponse{}, fmt.Errorf("ai service returned status %d", resp.StatusCode)
	}

	var out PredictResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return PredictResponse{}, fmt.Errorf("decode response: %w", err)
	}
	return out, nil
}