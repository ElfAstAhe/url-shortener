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
	err := _cfg.GlobalConfig.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)

		os.Exit(1)
	}

	fmt.Println("Initializing server...")
	var router = _hnd.BuildRouter()

	fmt.Println("Starting server...")
	if err := http.ListenAndServe(_cfg.GlobalConfig.HTTP.GetHost(), router); err != nil {
		fmt.Println("Error starting server:", err)
		panic(err)
	}

	fmt.Println("Server stopped")
}
