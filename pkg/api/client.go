package api

import (
	"context"
	"country-search-api/pkg/models"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/tidwall/gjson"
)

const defaultBaseURL = "https://restcountries.com/v3.1"

var (
	ErrNotFound = errors.New("country not found")
	ErrUpstream = errors.New("upstream service error")
)

type Client struct {
	baseURL string
	http    *http.Client
}

func NewClient(timeout time.Duration) *Client {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	return &Client{
		baseURL: defaultBaseURL,
		http: &http.Client{
			Timeout:   timeout,
			Transport: transport,
		},
	}
}

func (c *Client) GetByName(ctx context.Context, name string) (country *models.Country, err error) {
	url := fmt.Sprintf("%s/name/%s?fields=name,capital,currencies,population&fullText=true", c.baseURL, name)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, ErrNotFound
	default:
		return nil, ErrUpstream
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read details from %s: %w", url, err)
	}

	// fmt.Printf("Body: %s\n", body)

	country = &models.Country{}
	country.Name = gjson.GetBytes(body, "0.name.common").String()
	country.Currency = gjson.GetBytes(body, "0.currencies.*.symbol").String()
	country.Capital = gjson.GetBytes(body, "0.capital.0").String()
	country.Population = gjson.GetBytes(body, "0.population").Int()
	return
}
