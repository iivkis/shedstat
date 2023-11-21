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
	"golang.org/x/exp/slog"
)

/*
	/profile/:id
	/profile/:id/metrics
	/top?filter=[subscriptions,subscibers,likes]&amount[0-100]
*/

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	pgdb, err := repository.NewDBPostgres(repository.DBPostgresSource{
		Name:     config.Get().DB.Name,
		Host:     config.Get().DB.Host,
		Port:     config.Get().DB.Port,
		User:     config.Get().DB.User,
		Password: config.Get().DB.Password,
	})
	if err != nil {
		panic(err)
	}

	clickdb, err := repository.NewDBClickHouse(&clickhouse.Options{
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

	repoProfile := repository.NewProfilePostgresRepository(pgdb)
	repoProfileCollector := repository.NewProfileCollectorPostgresRepository(pgdb)
	repoMetrics := repository.NewMetricsClickHouseRepository(clickdb)
	repoMetricsCollector := repository.NewMetricsCollectorPostgresRepository(pgdb)

	svcProfile := services.NewProfileService(
		logger,
		repoProfile,
		repoProfileCollector,
		repoMetrics,
		shedevrumapi.NewShedevrumAPI(shedevrumapi.Config{}),
	)

	services.NewMetricsService(
		logger,
		repoMetrics,
		repoMetricsCollector,
		repoProfile,
		shedevrumapi.NewShedevrumAPI(shedevrumapi.Config{}),
	)

	router := chi.NewRouter()
	handlers.NewProfileHTTPHandler(svcProfile).Setup(router)

	fmt.Println("server is up")
	if err := http.ListenAndServe(":80", router); err != nil {
		panic(err)
	}
}
