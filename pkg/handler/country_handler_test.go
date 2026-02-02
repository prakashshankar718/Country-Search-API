package handler

import (
	mock_http_client "country-search-api/mock/ClientInf"
	"country-search-api/pkg/service/country"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := `[
		{
			"name": {"common": "India"},
			"capital": ["New Delhi"],
			"population": 1400000000,
			"currencies": {"INR": {"symbol": "₹"}}
		}
	]`

	mockClient := new(mock_http_client.MockClientInf)
	mockClient.On("Get", mock.Anything, mock.Anything).Return([]byte(body), nil).Once()

	ncs := country.NewCountryService(mockClient, "")
	ch := NewCountryHandler(ncs)

	r := gin.New()
	r.GET("/api/countries/search", ch.GetCountry)

	req := httptest.NewRequest(http.MethodGet, "/api/countries/search?name=India", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "New Delhi")
}
