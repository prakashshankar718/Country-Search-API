package main

import (
	"country-search-api/pkg/api"
	"country-search-api/pkg/logger"
	"log/slog"
)

func main() {

	logger.Init(slog.LevelInfo)
	api.RegisterRoutes()

}
