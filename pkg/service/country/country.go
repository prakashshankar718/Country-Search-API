package country

import (
	"context"
	"country-search-api/pkg/logger"
	"country-search-api/pkg/models"
	"country-search-api/pkg/service/cache"
	http_client "country-search-api/pkg/service/client"
	"fmt"
	"net/url"

	"github.com/tidwall/gjson"
)

// const defaultBaseURL = "https://restcountries.com/v3.1"

type CountryService interface {
	GetCountryByName(ctx context.Context, name string) (models.Country, error)
}

type countryService struct {
	httpClient http_client.ClientInf
	baseURL    string
}

func NewCountryService(httpClient http_client.ClientInf, baseURL string) CountryService {
	return &countryService{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

func (cs *countryService) GetCountryByName(ctx context.Context, name string) (models.Country, error) {
	// fmt.Println("GetCountryByName")
	logger.Log().Info("searching country details in local cache:", "country", name)
	if country, ok := cache.Cache.Get(name); ok {
		logger.Log().Info("country details present in local cache:", "country", name)
		return country.(models.Country), nil
	}

	logger.Log().Info("country details does not exist in local cache:", "country", name)
	logger.Log().Info("searching in 3rd party API:", "country", name)

	escaped := url.PathEscape(name)
	endpoint := fmt.Sprintf(
		"%s/name/%s?fields=name,capital,currencies,population&fullText=true",
		cs.baseURL,
		escaped,
	)

	countryBytes, err := cs.httpClient.Get(ctx, endpoint)
	if err != nil {
		logger.Log().Error("unable to get country details from 3rd party API:", "country", name)
		return models.Country{}, err
	}

	country := models.Country{
		Name:       gjson.GetBytes(countryBytes, "0.name.common").String(),
		Currency:   gjson.GetBytes(countryBytes, "0.currencies.*.symbol").String(),
		Capital:    gjson.GetBytes(countryBytes, "0.capital.0").String(),
		Population: gjson.GetBytes(countryBytes, "0.population").Int(),
	}

	if !country.Validate() {
		// fmt.Println(country)
		return models.Country{}, http_client.ErrInvalidData
	}

	logger.Log().Info("storing country details in local cache:", "country", name)
	go cache.Cache.Set(name, country)
	return country, nil
}
