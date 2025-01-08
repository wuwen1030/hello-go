package main

import (
	"log"

	"github.com/wuwen/hello-go/internal/app"
)

func main() {
	app := app.New()

	if err := app.Initialize(); err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("Error running app: %v", err)
	}
}
