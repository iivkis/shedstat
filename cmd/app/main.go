package main

import (
	"fmt"
	"net/http"
	"os"
	"shedstat/internal/adapters/handlers"
	"shedstat/internal/adapters/repository"
	"shedstat/internal/core/config"
	"shedstat/internal/core/services"
	shedevrumapi "shedstat/pkg/shedevrum-api"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"golang.org/x/exp/slog"
)

/*
	/profile/:id
	/profile/:id/metrics
	/top/profiles?filter=[subscriptions,subscibers,likes]&amount[0-100]
*/

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	postgresDB, err := repository.NewDBPostgres(repository.DBPostgresSource{
		Name:     config.Get().DB.Name,
		Host:     config.Get().DB.Host,
		Port:     config.Get().DB.Port,
		User:     config.Get().DB.User,
		Password: config.Get().DB.Password,
	})
	if err != nil {
		panic(err)
	}

	clickhouseDB, err := repository.NewDBClickHouse(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%s", config.Get().ClickHouse.Host, config.Get().ClickHouse.Port)},
		Auth: clickhouse.Auth{
			Database: config.Get().DB.Name,
			Username: config.Get().DB.User,
			Password: config.Get().DB.Password,
		},
	})
	if err != nil {
		panic(err)
	}

	repoProfile := repository.NewProfilePostgresRepository(postgresDB)
	repoProfileCollector := repository.NewProfileCollectorPostgresRepository(postgresDB)
	repoProfileMetrics := repository.NewProfileMetricsClickHouseRepository(clickhouseDB)
	repoProfileMetricsCollector := repository.NewProfileMetricsCollectorPostgresRepository(postgresDB)

	svcProfile := services.NewProfileService(
		logger,
		repoProfile,
		repoProfileCollector,
		repoProfileMetrics,
		repoProfileMetricsCollector,
		shedevrumapi.NewShedevrumAPI(shedevrumapi.Config{}),
	)

	router := chi.NewRouter()
	router.Use(newCorsHandler())
	router.Route("/api/v1/", func(r chi.Router) {
		handlers.NewProfileHTTPHandler(svcProfile).Setup(r)
	})

	fmt.Println("server is up")
	if err := http.ListenAndServe(":80", router); err != nil {
		panic(err)
	}
}

func newCorsHandler() func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	})
}
