package api

import (
	"context"
	"country-search-api/pkg/handler"
	"country-search-api/pkg/logger"
	http_client "country-search-api/pkg/service/client"
	"country-search-api/pkg/service/country"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// curl http://localhost:8080/api/countries/search?name=India
// var restcountries = "https://restcountries.com/v3.1/name/{name}?fields=name,capital,currencies,population&fullText=true"

func RegisterRoutes() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	router := gin.Default()
	router.Use(TimeoutMiddleware(20 * time.Second))

	defaultBaseURL := "https://restcountries.com/v3.1"
	httpClient := http_client.NewHTTPClient(5*time.Second, nil)
	counryService := country.NewCountryService(
		httpClient,
		defaultBaseURL,
	)
	countryHandler := handler.NewCountryHandler(counryService)

	router.GET("/api/countries/search", countryHandler.GetCountry)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log().Error("listen: %s\n", "exit", err)
			os.Exit(1)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	logger.Log().Warn("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Log().Warn("Server forced to shutdown: ", "Shutdown", err)
	}

	logger.Log().Warn("Server exiting")
}

func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// func countriesHandler(c *gin.Context) {
// 	countryName := c.DefaultQuery("name", "India")

// 	country, err := country.NewCountryService().GetCountryByName(c.Request.Context(), countryName)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, http_client.ErrNotFound):
// 			c.JSON(http.StatusNotFound, gin.H{"error": "country not found"})

// 		case errors.Is(err, context.DeadlineExceeded):
// 			c.JSON(http.StatusRequestTimeout, gin.H{"error": "request timeout"})

// 		case errors.Is(err, http_client.ErrInvalidData):
// 			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "could not validate country details"})

// 		case errors.Is(err, http_client.ErrUpstream):
// 			c.JSON(http.StatusBadGateway, gin.H{"error": "upstream service error"})

// 		default:
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to get country details"})
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, country)
// }

// func getCountryFromAPI(countryName string) (country *models.Country, err error) {
// 	url := strings.ReplaceAll(restcountries, "{name}", countryName)

// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to get details from %s: %w", url, err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("unable to get details from %s: %s", url, resp.Status)
// 	}

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to read details from %s: %w", url, err)
// 	}

// 	// fmt.Printf("Body: %s\n", body)

// 	country = &models.Country{}
// 	country.Name = gjson.GetBytes(body, "0.name.common").String()
// 	country.Currency = gjson.GetBytes(body, "0.currencies.*.symbol").String()
// 	country.Capital = gjson.GetBytes(body, "0.capital.0").String()
// 	country.Population = gjson.GetBytes(body, "0.population").Int()
// 	return
// }
