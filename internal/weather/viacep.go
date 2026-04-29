package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type ViaCEPClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewViaCEPClient(httpClient *http.Client, baseURL string) *ViaCEPClient {
	return &ViaCEPClient{
		httpClient: httpClient,
		baseURL:    strings.TrimRight(baseURL, "/"),
	}
}

func (c *ViaCEPClient) ResolveCity(ctx context.Context, zipcode string) (string, error) {
	endpoint := fmt.Sprintf("%s/%s/json/", c.baseURL, url.PathEscape(zipcode))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return "", ErrZipcodeNotFound
	}
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return "", fmt.Errorf("viacep returned status %d", res.StatusCode)
	}

	var payload struct {
		Localidade string `json:"localidade"`
		Erro       bool   `json:"erro"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return "", err
	}
	if payload.Erro || strings.TrimSpace(payload.Localidade) == "" {
		return "", ErrZipcodeNotFound
	}

	return payload.Localidade, nil
}
