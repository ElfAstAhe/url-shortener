package main

import (
	"fmt"
	"os"

	"github.com/ElfAstAhe/url-shortener/internal/bootstrap"
)

func main() {
	app, err := bootstrap.NewApp()
	if err != nil {
		fmt.Println("Error instantiate app", err)

		os.Exit(1)
	}

	if err := app.Init(); err != nil {
		fmt.Println("Error initialize app", err)

		os.Exit(1)
	}

	fmt.Println("Starting server...")
	if err := app.Run(); err != nil {
		fmt.Printf("Error starting server with error [%v]", err)

		os.Exit(1)
	}

	os.Exit(0)
}
