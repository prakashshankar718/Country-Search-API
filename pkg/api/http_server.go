package api

import (
	"context"
	"country-search-api/pkg/cache"
	"country-search-api/pkg/logger"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// curl http://localhost:8080/api/countries/search?name=India
// var restcountries = "https://restcountries.com/v3.1/name/{name}?fields=name,capital,currencies,population&fullText=true"
var countryClient = NewClient(7 * time.Second)

func RegisterRoutes() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	router := gin.Default()
	router.Use(TimeoutMiddleware(5 * time.Second))
	router.GET("/api/countries/search", countriesSearch)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
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

func countriesSearch(c *gin.Context) {
	countryName := c.DefaultQuery("name", "India")

	logger.Log().Info("searching country details in local cache:", "country", countryName)
	if country, ok := cache.Cache.Get(countryName); ok {
		logger.Log().Info("country details present in local cache:", "country", countryName)
		c.JSON(http.StatusOK, country)
		return
	}

	logger.Log().Info("country details does not exist in local cache:", "country", countryName)
	logger.Log().Info("searching in 3rd party API:", "country", countryName)
	// country, err := getCountryFromAPI(countryName)
	// if err != nil {
	// 	logger.Log().Error("falied getCountryFromRC()", "error", err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }

	// ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	// defer cancel()

	country, err := countryClient.GetByName(c.Request.Context(), countryName)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			c.JSON(404, gin.H{"error": "country not found"})
		case errors.Is(err, context.DeadlineExceeded):
			c.JSON(504, gin.H{"error": "request timeout"})
		default:
			c.JSON(502, gin.H{"error": "upstream service error"})
		}
		return
	}

	logger.Log().Info("storing country details in local cache:", "country", countryName)
	go cache.Cache.Set(countryName, country)

	c.JSON(http.StatusOK, country)
}

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
