// @title           Core API
// @version         0.1.0
// @description     Stack application core REST API.
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Pass `Bearer <token>`. JWT is verified by the gateway; the user id arrives as the X-User-Id header.
package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	emw "github.com/labstack/echo/v4/middleware"

	"app/core/health"
	"app/core/middleware"
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

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	pool, err := db.Open(ctx, dbURL)
	if err != nil {
		log.Fatalf("db open: %v", err)
	}
	defer pool.Close()
	if err := db.Migrate(ctx, pool, "./migrations/postgres"); err != nil {
		log.Fatalf("db migrate: %v", err)
	}

	e := echo.New()
	e.HideBanner = true
	e.Use(emw.Recover(), emw.RequestID(), observability.EchoMiddleware("core"))

	health.Register(e)

	// Protected API. The web app proxies browser requests to /api/core/* which
	// map 1:1 onto these routes (see apps/web routes/api/core/$.ts). Every
	// request carries the `token` cookie verified by middleware.RequireAuth.
	protected := e.Group("/v1", middleware.RequireAuth())
	protected.GET("/me", func(c echo.Context) error {
		u, _ := middleware.GetUser(c)
		return c.JSON(http.StatusOK, map[string]string{
			"user_id": u.ID,
			"email":   u.Email,
			"name":    u.Name,
		})
	})
	project.NewService(pool).RegisterRoutes(protected)

	port := os.Getenv("PORT")
	if port == "" {
		port = "4100"
	}
	log.Printf("core listening on :%s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatal(err)
	}
}
