package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/digstack/go-eda/pkg/consumer"
	"github.com/labstack/echo/v4"
	emw "github.com/labstack/echo/v4/middleware"
	"github.com/nats-io/nats.go"

	"app/core/health"
	"app/core/observability"
	"app/core/pkg/db"
	"app/core/project"
	projectevents "app/core/project/events"
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

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}
	nc, err := nats.Connect(natsURL,
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2*time.Second),
	)
	if err != nil {
		log.Fatalf("nats connect: %v", err)
	}
	defer nc.Drain()

	// Durable JetStream consumers. consumer.Run owns all redelivery semantics
	// (Term on permanent errors, bounded MaxDeliver, staged BackOff, the DLQ
	// stream, ack heartbeats, reconnect). Register one goroutine per handler;
	// add yours under apps/core/<context>/events/ — never hand-roll a
	// JetStream subscription.
	consumerCtx, stopConsumers := context.WithCancel(ctx)
	defer stopConsumers()
	go consumer.Run(consumerCtx, nc, projectevents.NewProjectCreatedHandler(), consumer.Config{})

	e := echo.New()
	e.HideBanner = true
	e.Use(emw.Recover(), emw.RequestID(), observability.EchoMiddleware("core"))

	health.Register(e)
	project.NewService(pool).RegisterRoutes(e.Group("/api/v1"))

	port := os.Getenv("PORT")
	if port == "" {
		port = "4100"
	}
	log.Printf("core listening on :%s (nats=%s)", port, natsURL)
	if err := e.Start(":" + port); err != nil {
		log.Fatal(err)
	}
}
