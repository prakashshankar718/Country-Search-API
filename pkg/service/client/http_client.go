package http_client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

var (
	ErrNotFound    = errors.New("not found")
	ErrUpstream    = errors.New("upstream service error")
	ErrInvalidData = errors.New("invalid data")
)

type ClientInf interface {
	Get(ctx context.Context, url string) ([]byte, error)
}

type client struct {
	http *http.Client
}

// var RestClient ClientInf = NewHTTPClient(5*time.Second, nil)

func NewHTTPClient(timeout time.Duration, transport http.RoundTripper) ClientInf {
	if transport == nil {
		transport = &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 20,
			MaxConnsPerHost:     50,
			IdleConnTimeout:     30 * time.Second,
			TLSHandshakeTimeout: 10 * time.Second,
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
		}
	}

	return &client{
		http: &http.Client{
			Timeout:   timeout,
			Transport: transport,
		},
	}
}

func (c *client) Get(ctx context.Context, url_str string) ([]byte, error) {
	if url_str == "" {
		return nil, ErrInvalidData
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url_str, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		if errors.Is(err, context.Canceled) ||
			errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}
		// fmt.Println("ErrUpstream")
		return nil, fmt.Errorf("%w: %v", ErrUpstream, err)
	}
	defer resp.Body.Close()
	// fmt.Println(resp.StatusCode)
	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, ErrNotFound
	default:
		return nil, fmt.Errorf("%w: status %d", ErrUpstream, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read details from %s: %w", url_str, err)
	}

	// fmt.Printf("Body: %s\n", body)
	return body, nil
}

// func (c *client) Get(ctx context.Context, endpoint string) (any, error) {
// 	if endpoint == "" {
// 		return nil, ErrInvalid
// 	}

// 	// escaped := url.PathEscape(url_str)
// 	// endpoint := fmt.Sprintf(
// 	// 	"%s/name/%s?fields=name,capital,currencies,population&fullText=true",
// 	// 	c.baseURL,
// 	// 	escaped,
// 	// )

// 	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	resp, err := c.http.Do(req)
// 	if err != nil {
// 		if errors.Is(err, context.Canceled) ||
// 			errors.Is(err, context.DeadlineExceeded) {
// 			return nil, err
// 		}
// 		fmt.Println("ErrUpstream")
// 		return nil, fmt.Errorf("%w: %v", ErrUpstream, err)
// 	}
// 	defer resp.Body.Close()

// 	switch resp.StatusCode {
// 	case http.StatusOK:
// 	case http.StatusNotFound:
// 		return nil, ErrNotFound
// 	default:
// 		return nil, fmt.Errorf("%w: status %d", ErrUpstream, resp.StatusCode)
// 	}

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to read details from %s: %w", endpoint, err)
// 	}

// 	// fmt.Printf("Body: %s\n", body)

// 	country := &models.Country{
// 		Name:       gjson.GetBytes(body, "0.name.common").String(),
// 		Currency:   gjson.GetBytes(body, "0.currencies.*.symbol").String(),
// 		Capital:    gjson.GetBytes(body, "0.capital.0").String(),
// 		Population: gjson.GetBytes(body, "0.population").Int(),
// 	}

// 	if !country.Validate() {
// 		// fmt.Println(country)
// 		return nil, ErrInvalid
// 	}

// 	return country, nil
// }
