package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ElfAstAhe/url-shortener/internal/config"
	"github.com/ElfAstAhe/url-shortener/internal/config/db"
	"github.com/ElfAstAhe/url-shortener/internal/handler"
	"github.com/ElfAstAhe/url-shortener/internal/logger"
	"github.com/ElfAstAhe/url-shortener/internal/storage"
	"github.com/ElfAstAhe/url-shortener/internal/utils"
	"github.com/ElfAstAhe/url-shortener/migrations"
	"go.uber.org/zap"
)

type App struct {
	AppRouter handler.AppRouter
	DB        db.DB
	log       *zap.SugaredLogger
}

func NewApp() (*App, error) {
	return &App{}, nil
}

func (app *App) Init() error {
	fmt.Println("Loading config...")
	config.AppConfig = config.NewConfig()
	if err := config.AppConfig.LoadConfig(); err != nil {
		return err
	}

	fmt.Println("Initializing logger...")
	if err := app.initLogger(); err != nil {
		return err
	}

	app.log.Info("Initializing database...")
	if err := app.initDatabase(); err != nil {
		return err
	}

	// Load cache data
	if err := app.initCacheData(); err != nil {
		return err
	}

	// DB migrations
	if err := app.migrateDatabase(); err != nil {
		return err
	}

	app.log.Info("Initializing http server router...")
	app.AppRouter = handler.NewRouter(config.AppConfig, logger.Log.Sugar())

	return nil
}

func (app *App) Run() error {
	app.log.Info("Starting graceful shutdown go routine...")
	go app.gracefulShutdown()

	app.log.Info("Starting server...")
	if err := http.ListenAndServe(config.AppConfig.HTTP.GetListenerAddr(), app.AppRouter.GetRouter()); err != nil {
		app.log.Errorf("Error starting server with error [%v]", err)

		os.Exit(1)
	}

	return nil
}

func (app *App) initLogger() error {
	if err := logger.Initialize(config.AppConfig.LogLevel, config.AppConfig.ProjectStage); err != nil {
		return err
	}

	app.log = logger.Log.Sugar()

	return nil
}

func (app *App) initDatabase() error {
	db, err := db.NewDB(config.AppConfig.DBKind, config.AppConfig.DBDsn)
	if err != nil {
		app.log.Errorf("Failed to initialize database: [%v]", err)
		return err
	}
	app.DB = db

	return nil
}

func (app *App) initCacheData() error {
	if cache, ok := app.DB.(db.InMemoryCache); ok {
		app.log.Info("Load data from storage...")
		if err := app.loadShortURIData(config.AppConfig.StoragePath, cache); err != nil {
			app.log.Errorf("Error loading data: [%v]", err)
			app.log.Warn("Using empty data storage")
		}

		app.log.Info("Load data from storage user...")
		if err := app.loadShortURIUserData(config.AppConfig.StorageUserPath, cache); err != nil {
			app.log.Errorf("Error loading data: [%v]", err)
			app.log.Warn("Using empty data storage")
		}
	}

	return nil
}

func (app *App) migrateDatabase() error {
	if app.DB.GetDBKind() == config.DBKindPostgres {
		app.log.Info("DB migrations postgres...")

		migrator, err := migrations.NewGooseDBMigrator(context.Background(), app.DB.GetDB(), logger.Log)
		if err != nil {
			app.log.Errorf("Error instantiate DB migrator: [%v]", err)

			return err
		}

		if err := migrator.Initialize(); err != nil {
			app.log.Errorf("Error initializing DB migrator: [%v]", err)

			return err
		}

		if err := migrator.Up(); err != nil {
			app.log.Errorf("Error DB migrate up: [%v]", err)

			return err
		}
	}

	return nil
}

func (app *App) gracefulShutdown() {
	// channel
	sig := make(chan os.Signal, 1)
	// register channel signals
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	// awaiting signal
	<-sig

	if cache, ok := app.DB.(db.InMemoryCache); ok {
		if err := app.saveShortURIData(config.AppConfig.StoragePath, cache); err != nil {
			app.log.Errorf("Error save shortURI data: [%v]", err)
		}
		if err := app.saveShortURIUserData(config.AppConfig.StorageUserPath, cache); err != nil {
			app.log.Errorf("Error save shortURI user data: [%v]", err)
		}
	}

	if err := db.CloseDB(app.DB); err != nil {
		app.log.Errorf("Error closing database: [%v]", err)
	}

	app.log.Info("Shutting down server done")

	os.Exit(0)
}

func (app *App) loadShortURIData(storagePath string, cache db.InMemoryCache) error {
	storageReader, err := storage.NewShortURLStorageReader(storagePath)
	if err != nil {
		return err
	}
	defer utils.CloseOnly(storageReader)

	return storageReader.LoadData(cache.GetShortURICache())
}

func (app *App) loadShortURIUserData(storagePath string, cache db.InMemoryCache) error {
	storageReader, err := storage.NewShortURIUserStorageReader(storagePath)
	if err != nil {
		return err
	}
	defer utils.CloseOnly(storageReader)

	return storageReader.LoadData(cache.GetShortURIUserCache())
}

func (app *App) saveShortURIData(storagePath string, cache db.InMemoryCache) error {
	storageWriter, err := storage.NewShortURLStorageWriter(storagePath)
	if err != nil {
		return err
	}
	defer utils.CloseOnly(storageWriter)

	return storageWriter.SaveData(cache.GetShortURICache())
}

func (app *App) saveShortURIUserData(storagePath string, cache db.InMemoryCache) error {
	storageWriter, err := storage.NewShortURIUserStorageWriter(storagePath)
	if err != nil {
		return err
	}
	defer utils.CloseOnly(storageWriter)

	return storageWriter.SaveData(cache.GetShortURIUserCache())
}
