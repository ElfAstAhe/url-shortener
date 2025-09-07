package bootstrap

import (
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
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type App struct {
	Router chi.Router
	log    *zap.SugaredLogger
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
	if err := _log.Initialize(_cfg.AppConfig.LogLevel, _cfg.AppConfig.ProjectStage); err != nil {
		return err
	}

	app.log = _log.Log.Sugar()

	app.log.Info("Load data from storage...")
	if err := loadShortURIData(_cfg.AppConfig.StoragePath); err != nil {
		app.log.Errorf("Error loading data: [%v]", err)
		app.log.Warn("Using empty data storage")
	}

	app.log.Info("Initializing http server router...")
	app.Router = _hnd.BuildRouter()

	return nil
}

func (app *App) Run() error {
	app.log.Info("Starting graceful shutdown go routine...")
	go app.gracefulShutdown()

	app.log.Info("Starting server...")
	if err := http.ListenAndServe(_cfg.AppConfig.HTTP.GetListenerAddr(), app.Router); err != nil {
		app.log.Errorf("Error starting server with error [%v]", err)

		os.Exit(1)
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

	if err := saveShortURIData(_cfg.AppConfig.StoragePath); err != nil {
		app.log.Errorf("Error save shortURI data: [%v]", err)
	}
}

func loadShortURIData(storagePath string) error {
	storageReader, err := _storage.NewShortURLStorageReader(storagePath)
	if err != nil {
		return err
	}
	defer storageReader.Close()

	return storageReader.LoadData(_db.InMemoryDBInstance.ShortURI)
}

func saveShortURIData(storagePath string) error {
	storageWriter, err := _storage.NewShortURLStorageWriter(storagePath)
	if err != nil {
		return err
	}
	defer storageWriter.Close()

	return storageWriter.SaveData(_db.InMemoryDBInstance.ShortURI)
}
