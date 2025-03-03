package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	eventDelivery "github.com/Inspirate789/grpc-template/internal/event/delivery"
	eventRepository "github.com/Inspirate789/grpc-template/internal/event/repository"
	eventUsecase "github.com/Inspirate789/grpc-template/internal/event/usecase"
	"github.com/Inspirate789/grpc-template/internal/pkg/app"
	userDelivery "github.com/Inspirate789/grpc-template/internal/user/delivery"
	userRepository "github.com/Inspirate789/grpc-template/internal/user/repository"
	userUsecase "github.com/Inspirate789/grpc-template/internal/user/usecase"
	"github.com/Inspirate789/grpc-template/pkg/migrations"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/lmittmann/tint"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/pflag"
	"golang.org/x/sync/errgroup"
)

func startApp(webApp *app.WebApp, grpcApp *app.GrpcApp, config app.Config, logger *slog.Logger) {
	logger.Debug(fmt.Sprintf("app starts with configuration: %+v", config))

	go func() {
		err := webApp.Start()
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		err := grpcApp.Start()
		if err != nil {
			panic(err)
		}
	}()
}

func shutdownApp(webApp *app.WebApp, grpcApp *app.GrpcApp, logger *slog.Logger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	const shutdownTimeout = time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)

	eg := errgroup.Group{}

	eg.Go(func() error {
		logger.Debug("shutdown web app ...")

		err := webApp.Shutdown(ctx)
		if err != nil {
			return err
		}

		logger.Debug("web app exited")

		return nil
	})

	eg.Go(func() error {
		logger.Debug("shutdown grpc app ...")
		grpcApp.Shutdown()
		logger.Debug("grpc app exited")

		return nil
	})

	err := eg.Wait()
	if err != nil {
		panic(err)
	}

	cancel()
}

func main() {
	var configPath, migrationsPath string
	pflag.StringVarP(&configPath, "config", "c", "configs/app.yaml", "Config file path")
	pflag.StringVarP(&migrationsPath, "migrations", "m", "migrations", "Migrations directory path")
	pflag.Parse()

	config, err := app.ReadLocalConfig(configPath)
	if err != nil {
		panic(err)
	}

	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.Level(config.Logging.Level)}))

	db, err := sqlx.Connect(config.DB.DriverName, config.DB.ConnectionString)
	if err != nil {
		panic(err)
	}

	defer func(db *sqlx.DB) {
		err = db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	dbInstance, err := sqlite3.WithInstance(db.DB, &sqlite3.Config{})
	if err != nil {
		panic(err)
	}

	err = migrations.Do(config.DB.DriverName, migrationsPath, dbInstance, logger)
	if err != nil {
		panic(err)
	}

	users := userDelivery.New(userUsecase.New(userRepository.NewSqlx(db, logger), logger), logger)
	events := eventDelivery.New(eventUsecase.New(eventRepository.NewSqlx(db, logger), logger), logger)

	webApp := app.NewWebApp(config.Web, nil, nil, logger)
	grpcApp := app.NewGrpcApp(config.GRPC, logger, users, events)

	startApp(webApp, grpcApp, config, logger)
	shutdownApp(webApp, grpcApp, logger)
}
