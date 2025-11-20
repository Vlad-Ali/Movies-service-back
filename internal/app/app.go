package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Vlad-Ali/Movies-service-back/internal/application/usecase/transactionmanager"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var (
	ConfigPath = "config.yml"
)

type App struct {
	repositories *Repositories
	services     *Services
	handlers     *Handlers
	server       *http.Server
	db           *sql.DB
	config       *Config
}

func NewApp() (*App, error) {
	setLogger()
	cfg, err := LoadConfig(ConfigPath)
	if err != nil {
		slog.Error("Error loading config: ", "error", err)
		return nil, err
	}
	db, err := sql.Open("postgres", cfg.PostgresConfig.DSN)
	if err != nil {
		slog.Error("Error opening connection to database", "error", err)
		return nil, err
	}

	db.SetMaxIdleConns(cfg.PostgresConfig.MaxIdleConnections)
	db.SetMaxOpenConns(cfg.PostgresConfig.MaxConnections)
	db.SetConnMaxLifetime(cfg.PostgresConfig.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.PostgresConfig.ConnMaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		slog.Error("Error connecting to database", "error", err)
		return nil, err
	}

	txUser := transactionmanager.NewTransactionUser(db)
	repos := NewRepositories(db)
	services := NewServices(db, repos, cfg.SecretKey, txUser, cfg.ModelConfig)
	handlers := NewHandlers(services)
	handler := handlers.registerRoutes()

	slog.Info("Successfully connected to PostgreSQL")
	return &App{db: db, config: cfg, handlers: handlers, services: services,
		repositories: repos, server: &http.Server{Addr: cfg.Address,
			Handler: handler, WriteTimeout: cfg.WriteTimeout, ReadTimeout: cfg.ReadTimeout}}, nil
}

func (a *App) Start() error {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	slog.Info("Starting application...")

	err := a.RunMigrations()
	if err != nil {
		slog.Error("Error running migrations", "error", err)
		a.db.Close()
		return err
	}

	go func() {
		slog.Info(fmt.Sprintf("Server started at %s", a.server.Addr))
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error(fmt.Sprintf("Server failed: %v", err))
			os.Exit(1)
		}
	}()

	<-signalChan
	slog.Info("Shutting down...")

	slog.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	servErr := a.server.Shutdown(ctx)
	if servErr != nil {
		slog.Error("Error shutting down server", "error", servErr)
	}

	dbErr := a.db.Close()
	if dbErr != nil {
		slog.Error("Error closing database", "error", dbErr)
	}

	if dbErr != nil || servErr != nil {
		return fmt.Errorf("error during application shutdown: %v, %v", servErr, dbErr)
	}

	slog.Info("Successfully closed database")
	return nil
}

func (a *App) RunMigrations() error {
	driver, err := postgres.WithInstance(a.db, &postgres.Config{})
	if err != nil {
		slog.Error("Error creating database driver", "error", err)
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		slog.Error("Error creating migration instance", "error", err)
		return err
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error("Error running migration", "error", err)
		return err
	}

	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		slog.Error("Error getting version", "error", err)
		return err
	}

	slog.Info("Migrations complete", "version", version, "dirty", dirty)
	return nil
}

func (a *App) Migrate(version uint) error {
	driver, err := postgres.WithInstance(a.db, &postgres.Config{})
	if err != nil {
		slog.Error("Error creating database driver", "error", err)
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		slog.Error("Error creating migration instance", "error", err)
		return err
	}

	err = m.Migrate(version)
	if err != nil {
		slog.Error("Error running migration", "error", err)
		return err
	}

	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		slog.Error("Error getting version", "error", err)
		return err
	}
	slog.Info("Migration complete", "version", version, "dirty", dirty)
	return nil
}

func (a *App) ForceMigration(versionToMigrate int) error {
	driver, err := postgres.WithInstance(a.db, &postgres.Config{})
	if err != nil {
		slog.Error("Error creating database driver", "error", err)
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		slog.Error("Error creating migration instance", "error", err)
		return err
	}

	err = m.Force(versionToMigrate)
	if err != nil {
		slog.Error("Error running migration", "error", err)
		return err
	}

	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		slog.Error("Error getting version", "error", err)
		return err
	}
	slog.Info("Migration complete", "version", version, "dirty", dirty)
	return nil
}

func setLogger() {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
