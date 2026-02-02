package handler

import (
	"context"
	http_client "country-search-api/pkg/service/client"
	"country-search-api/pkg/service/country"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CountryHandler struct {
	cs country.CountryService
}

func NewCountryHandler(cs country.CountryService) *CountryHandler {
	return &CountryHandler{cs: cs}
}

func (ch *CountryHandler) GetCountry(c *gin.Context) {
	// fmt.Println("GetCountry")

	countryName := c.DefaultQuery("name", "India")

	country, err := ch.cs.GetCountryByName(c.Request.Context(), countryName)
	if err != nil {
		switch {
		case errors.Is(err, http_client.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "country not found"})

		case errors.Is(err, context.DeadlineExceeded):
			c.JSON(http.StatusRequestTimeout, gin.H{"error": "request timeout"})

		case errors.Is(err, http_client.ErrInvalidData):
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "could not validate country details"})

		case errors.Is(err, http_client.ErrUpstream):
			c.JSON(http.StatusBadGateway, gin.H{"error": "upstream service error"})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to get country details"})
		}
		return
	}

	c.JSON(http.StatusOK, country)
}
