package main

import (
	"fmt"
	"net/http"
	"os"

	_cfg "github.com/ElfAstAhe/url-shortener/internal/config"
	_hnd "github.com/ElfAstAhe/url-shortener/internal/handler"
)

func main() {
	fmt.Println("Loading config...")
	_cfg.AppConfig = _cfg.NewConfig()
	if err := _cfg.AppConfig.LoadConfig(); err != nil {
		fmt.Println("Error loading config:", err)

		os.Exit(1)
	}

	fmt.Println("Initializing http server router...")
	var router = _hnd.BuildRouter()

	fmt.Println("Starting server...")
	if err := http.ListenAndServe(_cfg.AppConfig.HTTP.GetListenerAddr(), router); err != nil {
		fmt.Printf("Error starting server with error [%v]", err)

		os.Exit(1)
	}

	fmt.Println("Server stopped")
}
