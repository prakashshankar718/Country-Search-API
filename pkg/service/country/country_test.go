package country

import (
	"context"
	mock_http_client "country-search-api/mock/ClientInf"
	"country-search-api/pkg/models"
	http_client "country-search-api/pkg/service/client"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetCountryByName_Success(t *testing.T) {
	mockClient := new(mock_http_client.MockClientInf)

	body := `[
		{
			"name": {"common": "India"},
			"capital": ["New Delhi"],
			"population": 1400000000,
			"currencies": {"INR": {"symbol": "₹"}}
		}
	]`

	expected := models.Country{
		Name:       "India",
		Capital:    "New Delhi",
		Population: 1400000000,
		Currency:   "₹",
	}

	mockClient.On("Get", mock.Anything, mock.Anything).Return([]byte(body), nil).Once()
	ncs := NewCountryService(mockClient, "defaultBaseURL")

	country, err := ncs.GetCountryByName(context.Background(), "India")

	// fmt.Println(country, err)

	assert.NoError(t, err)
	assert.Equal(t, expected, country)
	mockClient.AssertExpectations(t)
}

func TestGetCountryByName_InvalidData(t *testing.T) {
	mockClient := new(mock_http_client.MockClientInf)

	body := `[{ "name": {} }]` // missing required fields

	mockClient.On("Get", mock.Anything, mock.Anything).Return([]byte(body), nil).Once()
	ncs := NewCountryService(mockClient, "defaultBaseURL")

	country, err := ncs.GetCountryByName(context.Background(), "India")

	assert.ErrorIs(t, err, http_client.ErrInvalidData)
	assert.Equal(t, country, models.Country{})
	mockClient.AssertExpectations(t)
}
