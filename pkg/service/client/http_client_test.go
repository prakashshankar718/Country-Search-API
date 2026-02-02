package http_client

import (
	"context"
	"country-search-api/pkg/models"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tidwall/gjson"
)

type MockRoundTripper struct {
	mock.Mock
}

// // func (m *MockRoundTripper) RoundTrip(
// // 	req *http.Request,
// // ) (*http.Response, error) {
// // 	args := m.Called(req)
// // 	return args.Get(0).(*http.Response), args.Error(1)
// // }

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	args := m.Called(req)

	// Check if the first return argument is a function
	if rf, ok := args.Get(0).(func(*http.Request) *http.Response); ok {
		return rf(req), args.Error(1)
	}

	return args.Get(0).(*http.Response), args.Error(1)
}

func TestClient_Get_Success(t *testing.T) {
	body := `[
		{
			"name": {"common": "India"},
			"capital": ["New Delhi"],
			"population": 1400000000,
			"currencies": {"INR": {"symbol": "₹"}}
		}
	]`

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}

	rt := new(MockRoundTripper)
	rt.On("RoundTrip", mock.Anything).Return(resp, nil)

	client := NewHTTPClient(5*time.Second, rt)

	// countryBytes, err := client.Get(context.Background(), "https://restcountries.com/v3.1/name/India?fields=name,capital,currencies,population&fullText=true")
	countryBytes, err := client.Get(context.Background(), "india")
	assert.NoError(t, err)

	country := models.Country{
		Name:       gjson.GetBytes(countryBytes, "0.name.common").String(),
		Currency:   gjson.GetBytes(countryBytes, "0.currencies.*.symbol").String(),
		Capital:    gjson.GetBytes(countryBytes, "0.capital.0").String(),
		Population: gjson.GetBytes(countryBytes, "0.population").Int(),
	}
	assert.NotNil(t, country)
	assert.Equal(t, "New Delhi", country.Capital)

	rt.AssertExpectations(t)
}

func TestClient_Get_NotFound(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       io.NopCloser(strings.NewReader("")),
	}

	rt := new(MockRoundTripper)
	rt.On("RoundTrip", mock.Anything).Return(resp, nil)

	client := NewHTTPClient(5*time.Second, rt)
	// _, err := client.Get(context.Background(), "Atlantis")
	_, err := client.Get(context.Background(), "Mumbai")

	assert.ErrorIs(t, err, ErrNotFound)
}

func TestGet_UpstreamError(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(strings.NewReader("boom")),
	}

	rt := new(MockRoundTripper)
	rt.On("RoundTrip", mock.Anything).Return(resp, nil)

	client := NewHTTPClient(2*time.Second, rt)

	_, err := client.Get(context.Background(), "India")

	assert.ErrorIs(t, err, ErrUpstream)
}

func TestGet_InvalidData(t *testing.T) {
	// body := `[{ "name": {} }]` // missing required fields
	// body := ""

	resp := &http.Response{
		// StatusCode: http.StatusOK,
		// Body:       io.NopCloser(strings.NewReader(body)),
	}

	rt := new(MockRoundTripper)
	rt.On("RoundTrip", mock.Anything).Return(resp, nil)

	client := NewHTTPClient(2*time.Second, rt)

	_, err := client.Get(context.Background(), "")

	assert.ErrorIs(t, err, ErrInvalidData)
}

func TestGet_ContextTimeout(t *testing.T) {
	rt := new(MockRoundTripper)

	rt.On("RoundTrip", mock.Anything).
		Run(func(args mock.Arguments) {
			time.Sleep(50 * time.Millisecond)
		}).
		Return(&http.Response{}, context.DeadlineExceeded)

	client := NewHTTPClient(100*time.Millisecond, rt)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := client.Get(ctx, "India")

	assert.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestGet_NetworkFailure(t *testing.T) {
	rt := new(MockRoundTripper)
	rt.On("RoundTrip", mock.Anything).
		Return(&http.Response{}, errors.New("connection refused"))

	client := NewHTTPClient(2*time.Second, rt)

	_, err := client.Get(context.Background(), "India")

	assert.ErrorIs(t, err, ErrUpstream)
}

func TestGet_Concurrent(t *testing.T) {
	body := `[{"name": {"common": "India"}, "capital": ["New Delhi"], "population": 1400000000, "currencies": {"INR": {"symbol": "₹"}}}]`
	rt := new(MockRoundTripper)

	rt.On("RoundTrip", mock.Anything).
		Return(func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(body)),
				Header:     make(http.Header),
			}
		}, nil)
	client := NewHTTPClient(5*time.Second, rt)

	wg := sync.WaitGroup{}
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			countryBytes, err := client.Get(context.Background(), "India")
			country := models.Country{
				Name:       gjson.GetBytes(countryBytes, "0.name.common").String(),
				Currency:   gjson.GetBytes(countryBytes, "0.currencies.*.symbol").String(),
				Capital:    gjson.GetBytes(countryBytes, "0.capital.0").String(),
				Population: gjson.GetBytes(countryBytes, "0.population").Int(),
			}
			assert.NoError(t, err)
			assert.NotNil(t, country)
			assert.Equal(t, "New Delhi", country.Capital)
		}()
	}

	wg.Wait()
	rt.AssertExpectations(t)
}

func BenchmarkGet(b *testing.B) {
	body := `[{"name": {"common": "India"}, "capital": ["New Delhi"], "population": 1400000000, "currencies": {"INR": {"symbol": "₹"}}}]`
	rt := new(MockRoundTripper)

	rt.On("RoundTrip", mock.Anything).
		Return(func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(body)),
				Header:     make(http.Header),
			}
		}, nil)
	client := NewHTTPClient(5*time.Second, rt)

	for b.Loop() {
		client.Get(context.Background(), "India")
	}
}
