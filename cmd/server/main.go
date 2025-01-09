package main

import (
	"log"

	_ "github.com/wuwen/hello-go/docs"
	"github.com/wuwen/hello-go/internal/app"
)

// @title        CMS API
// @version      1.0
// @description  A simple CMS system API
// @host        localhost:8080
// @BasePath    /api/v1
// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
// @type apiKey
func main() {
	app := app.New()

	if err := app.Initialize(); err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("Error running app: %v", err)
	}
}
