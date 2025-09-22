package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_cfg "github.com/ElfAstAhe/url-shortener/internal/config"
	_db "github.com/ElfAstAhe/url-shortener/internal/config/db"
	_hnd "github.com/ElfAstAhe/url-shortener/internal/handler"
	_log "github.com/ElfAstAhe/url-shortener/internal/logger"
	_storage "github.com/ElfAstAhe/url-shortener/internal/storage"
	_utl "github.com/ElfAstAhe/url-shortener/internal/utils"
	_migr "github.com/ElfAstAhe/url-shortener/migrations"
	"go.uber.org/zap"
)

type App struct {
	AppRouter _hnd.AppRouter
	DB        _db.DB
	log       *zap.SugaredLogger
}

func NewApp() (*App, error) {
	return &App{}, nil
}

func (app *App) Init() error {
	fmt.Println("Loading config...")
	_cfg.AppConfig = _cfg.NewConfig()
	if err := _cfg.AppConfig.LoadConfig(); err != nil {
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
	app.AppRouter = _hnd.NewRouter(_cfg.AppConfig, _log.Log.Sugar())

	return nil
}

func (app *App) Run() error {
	app.log.Info("Starting graceful shutdown go routine...")
	go app.gracefulShutdown()

	app.log.Info("Starting server...")
	if err := http.ListenAndServe(_cfg.AppConfig.HTTP.GetListenerAddr(), app.AppRouter.GetRouter()); err != nil {
		app.log.Errorf("Error starting server with error [%v]", err)

		os.Exit(1)
	}

	return nil
}

func (app *App) initLogger() error {
	if err := _log.Initialize(_cfg.AppConfig.LogLevel, _cfg.AppConfig.ProjectStage); err != nil {
		return err
	}

	app.log = _log.Log.Sugar()

	return nil
}

func (app *App) initDatabase() error {
	db, err := _db.NewDB(_cfg.AppConfig.DBKind, _cfg.AppConfig.DBDsn)
	if err != nil {
		app.log.Errorf("Failed to initialize database: [%v]", err)
		return err
	}
	app.DB = db

	return nil
}

func (app *App) initCacheData() error {
	if cache, ok := app.DB.(_db.InMemoryCache); ok {
		app.log.Info("Load data from storage...")
		if err := app.loadShortURIData(_cfg.AppConfig.StoragePath, cache); err != nil {
			app.log.Errorf("Error loading data: [%v]", err)
			app.log.Warn("Using empty data storage")
		}
	}

	return nil
}

func (app *App) migrateDatabase() error {
	if app.DB.GetDBKind() == _cfg.DBKindPostgres {
		app.log.Info("DB migrations postgres...")

		migrator, err := _migr.NewGooseDBMigrator(context.Background(), app.DB.GetDB(), _log.Log)
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

	if cache, ok := app.DB.(_db.InMemoryCache); ok {
		if err := app.saveShortURIData(_cfg.AppConfig.StoragePath, cache); err != nil {
			app.log.Errorf("Error save shortURI data: [%v]", err)
		}
	}

	if err := _db.CloseDB(app.DB); err != nil {
		app.log.Errorf("Error closing database: [%v]", err)
	}

	app.log.Info("Shutting down server done")

	os.Exit(0)
}

func (app *App) loadShortURIData(storagePath string, cache _db.InMemoryCache) error {
	storageReader, err := _storage.NewShortURLStorageReader(storagePath)
	if err != nil {
		return err
	}
	defer _utl.CloseOnly(storageReader)

	return storageReader.LoadData(cache.GetShortURICache())
}

func (app *App) loadShortURIUserData(storagePath string, cache _db.InMemoryCache) error {
	storageReader, err := _storage.NewShortURIUserStorageReader(storagePath)
	if err != nil {
		return err
	}
	defer _utl.CloseOnly(storageReader)

	return storageReader.LoadData(cache.GetShortURIUserCache())
}

func (app *App) saveShortURIData(storagePath string, cache _db.InMemoryCache) error {
	storageWriter, err := _storage.NewShortURLStorageWriter(storagePath)
	if err != nil {
		return err
	}
	defer _utl.CloseOnly(storageWriter)

	return storageWriter.SaveData(cache.GetShortURICache())
}

func (app *App) saveShortURIUserData(storagePath string, cache _db.InMemoryCache) error {
	storageWriter, err := _storage.NewShortURIUserStorageWriter(storagePath)
	if err != nil {
		return err
	}
	defer _utl.CloseOnly(storageWriter)

	return storageWriter.SaveData(cache.GetShortURIUserCache())
}
