package main

import (
	"fmt"
	_cfg "github.com/ElfAstAhe/url-shortener/internal/config"
	"github.com/ElfAstAhe/url-shortener/internal/handler"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Loading config...")
	err := _cfg.GlobalConfig.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)

		os.Exit(1)
	}

	fmt.Println("Initializing server...")

	mux := http.NewServeMux()
	mux.HandleFunc(handler.RootHandlePath, handler.RootHandler)

	fmt.Println("Starting server...")
	if err := http.ListenAndServe(_cfg.GlobalConfig.Http[0].GetHost(), mux); err != nil {
		fmt.Println("Error starting server:", err)
		panic(err)
	}

	fmt.Println("Server stopped")
}
