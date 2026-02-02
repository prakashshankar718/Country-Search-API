package intergration

import (
	"country-search-api/pkg/handler"
	http_client "country-search-api/pkg/service/client"
	"country-search-api/pkg/service/country"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCountryAPI_GetCountry_Integration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	httpClient := http_client.NewHTTPClient(5*time.Second, nil)
	endpoint := fmt.Sprintf("https://restcountries.com/v3.1/name/India?fields=name,capital,currencies,population&fullText=true")

	ncs := country.NewCountryService(httpClient, endpoint)
	nch := handler.NewCountryHandler(ncs)

	r := gin.New()

	r.GET("/api/countries/search", nch.GetCountry)

	server := httptest.NewServer(r)
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/countries/search?name=India")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
