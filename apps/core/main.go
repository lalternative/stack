package main

// @title           Core API
// @version         0.1.0
// @description     Stack application core REST API.
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Pass `Bearer <token>`. JWT is verified by the gateway; the user id arrives as the X-User-Id header.

import (
	"context"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	emw "github.com/labstack/echo/v4/middleware"

	"app/core/health"
	"app/core/observability"
	"app/core/pkg/db"
	"app/core/project"
)

func main() {
	ctx := context.Background()

	shutdown, err := observability.Init(ctx, observability.Config{
		ServiceName: "core",
		Endpoint:    os.Getenv("SKALPAI_ENDPOINT"),
		APIKey:      os.Getenv("SKALPAI_API_KEY"),
		ProjectID:   os.Getenv("SKALPAI_PROJECT_ID"),
	})
	if err != nil {
		log.Printf("observability init: %v (continuing without telemetry)", err)
	}
	defer shutdown(ctx)

	dbPath := os.Getenv("DUCKDB_PATH")
	if dbPath == "" {
		dbPath = "./data/app.duckdb"
	}
	exec, err := db.Open(dbPath)
	if err != nil {
		log.Fatalf("db open: %v", err)
	}
	if err := db.Migrate(ctx, exec, "./migrations/duckdb"); err != nil {
		log.Fatalf("db migrate: %v", err)
	}

	e := echo.New()
	e.HideBanner = true
	e.Use(emw.Recover(), emw.RequestID(), observability.EchoMiddleware("core"))

	health.Register(e)
	project.NewService(exec).RegisterRoutes(e.Group("/api/v1"))

	port := os.Getenv("PORT")
	if port == "" {
		port = "4100"
	}
	log.Printf("core listening on :%s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatal(err)
	}
}
