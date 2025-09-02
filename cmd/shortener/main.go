package main

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
)

func main() {
	fmt.Println("Loading config...")
	_cfg.AppConfig = _cfg.NewConfig()
	if err := _cfg.AppConfig.LoadConfig(); err != nil {
		fmt.Println("Error loading config:", err)

		os.Exit(1)
	}

	fmt.Println("Initializing logger...")
	if err := _log.Initialize(_cfg.AppConfig.LogLevel, _cfg.AppConfig.ProjectStage); err != nil {
		fmt.Println("Error initializing logger:", err)

		os.Exit(1)
	}

	fmt.Println("Load data from storage...")
	if err := _db.InMemoryDBInstance.LoadData(_cfg.AppConfig.StoragePath); err != nil {
		fmt.Println("Error loading data:", err)
	}

	fmt.Println("Initializing http server router...")
	var router = _hnd.BuildRouter()

	fmt.Println("Starting post action go routine...")
	go postAction()

	fmt.Println("Starting server...")
	if err := http.ListenAndServe(_cfg.AppConfig.HTTP.GetListenerAddr(), router); err != nil {
		fmt.Printf("Error starting server with error [%v]", err)

		os.Exit(1)
	}
}

func postAction() {
	// channel
	sig := make(chan os.Signal, 1)
	// register channel signals
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	// awaiting signal
	<-sig

	// save data
	fmt.Println("Save data into storage...")
	if err := _db.InMemoryDBInstance.SaveData(_cfg.AppConfig.StoragePath); err != nil {
		fmt.Println("Error saving data into storage:", err)

		os.Exit(1)
	}

	// shutdown
	fmt.Println("Server stopped")
	os.Exit(0)
}
