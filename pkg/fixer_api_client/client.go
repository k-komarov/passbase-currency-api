package fixer_api_client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/k-komarov/passbase-currency-api/pkg/types"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client interface {
	GetLatestEURToUSDRate(ctx context.Context) (*LatestRateResponse, error)
}

type client struct {
	BaseURL    string
	AccessKey  string
	HttpClient HTTPClient
	Logger     *logrus.Logger
}
type LatestRateResponse struct {
	Success bool `json:"success"`
	Error   struct {
		Info string `json:"info"`
	} `json:"error,omitempty"`
	Timestamp types.Timestamp    `json:"timestamp,omitempty"`
	Base      string             `json:"base,omitempty"`
	Rates     map[string]float64 `json:"rates,omitempty"`
}
type Option func(*client)

func NewClient(baseUrl, accessKey string, opts ...Option) Client {
	c := &client{
		BaseURL:    baseUrl,
		AccessKey:  accessKey,
		HttpClient: &http.Client{},
		Logger:     logrus.StandardLogger(),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}
func WithHttpClient(httpClient HTTPClient) Option {
	return func(c *client) {
		c.HttpClient = httpClient
	}
}

func (c *client) GetLatestEURToUSDRate(ctx context.Context) (*LatestRateResponse, error) {
	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, err
	}
	u.Path = fmt.Sprintf("%s/latest", u.Path)
	q := u.Query()
	q.Set("access_key", c.AccessKey)
	q.Set("symbols", "USD")
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.WithContext(ctx)
	resp, err := c.HttpClient.Do(req)
	defer resp.Body.Close()
	c.Logger.WithContext(ctx).Debugf("Request to Fixer.io: %s", u.String())
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(resp.Body)
	if err != nil {
		return nil, err
	}
	latestResponse := &LatestRateResponse{}
	err = decoder.Decode(latestResponse)
	if err != nil {
		return nil, err
	}

	return latestResponse, nil
}
