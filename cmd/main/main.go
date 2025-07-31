package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"

	"github.com/hibiken/asynq"

	"github.com/fedya-eremin/crypto-api/api"
	"github.com/fedya-eremin/crypto-api/api/impl"
	"github.com/fedya-eremin/crypto-api/clients/cmc"
	db_gen "github.com/fedya-eremin/crypto-api/database/gen"
	currency_repo "github.com/fedya-eremin/crypto-api/repo/currency"
	pricelog "github.com/fedya-eremin/crypto-api/repo/price-log"
	"github.com/fedya-eremin/crypto-api/service/currency"
	"github.com/fedya-eremin/crypto-api/tasks"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/golang-migrate/migrate/v4"
	migratePgx "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	DATABASE_URL string `env:"DATABASE_URL,required"`
	REDIS_URL    string `env:"REDIS_URL,required"`
	CMC_API_KEY  string `env:"CMC_API_KEY,required"`
}

func main() {
	ctx := context.Background()

	var cfg Config
	if err := envconfig.Process(ctx, &cfg); err != nil {
		slog.Error("Config parse failed:", "error", err)
		os.Exit(1)
	}

	pool, err := pgxpool.New(ctx, cfg.DATABASE_URL)
	if err != nil {
		slog.Error("Error connecting to DB:", "error", err)
		os.Exit(1)
	}
	defer pool.Close()
	database, err := sql.Open("postgres", cfg.DATABASE_URL)
	if err != nil {
		slog.Error("Error connecting to DB:", "error", err)
		pool.Close()
		os.Exit(1)
	}
	defer database.Close()

	driver, err := migratePgx.WithInstance(database, &migratePgx.Config{})
	if err != nil {
		slog.Error("Failed to init migrator", "error", err)
		pool.Close()
		database.Close()
		os.Exit(1)
	}

	migrator, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		slog.Error("Failed to init migrator", "error", err)
		pool.Close()
		database.Close()
		os.Exit(1)
	}
	err = migrator.Up()
	if err != nil && err != migrate.ErrNoChange {
		slog.Error("Failed to apply migrations", "error", err)
		pool.Close()
		database.Close()
		os.Exit(1)
	}
	slog.Info("migrations applied")

	redisOpt := asynq.RedisClientOpt{Addr: cfg.REDIS_URL}
	scheduler := asynq.NewScheduler(redisOpt, nil)
	defer scheduler.Shutdown()

	db := db_gen.New(pool)
	currencyRepo := currency_repo.New(db)
	pricelogRepo := pricelog.New(db)
	cmcClient := cmc.New(cfg.CMC_API_KEY, "USD")
	currencyService := currency.New(currencyRepo, cmcClient, scheduler)
	server := impl.New(currencyService)

	taskHandler := tasks.NewHandler(pricelogRepo, cmcClient)

	workerServer := asynq.NewServer(redisOpt, asynq.Config{Concurrency: 10})
	workerMux := asynq.NewServeMux()
	workerMux.HandleFunc(currency.TypeUpdateCurrencyTask, taskHandler.HandleCurrencyUpdateTask)
	defer workerServer.Shutdown()

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
	}))
	app.Use(logger.New())
	swaggerCfg := swagger.Config{
		BasePath: "/",
		FilePath: "./openapi.yaml",
		Path:     "docs",
		Title:    "Crypto Api Docs",
	}
	app.Use(swagger.New(swaggerCfg))
	api.RegisterHandlers(app, api.NewStrictHandler(server, nil))

	go func() {
		slog.Error("error processing tasks", "error", scheduler.Run())
		scheduler.Shutdown()
		pool.Close()
		database.Close()
		os.Exit(1)
	}()
	go func() {
		slog.Error("error processing tasks", "error", workerServer.Run(workerMux))
		workerServer.Shutdown()
		pool.Close()
		database.Close()
		os.Exit(1)
	}()
	err = currencyService.BootstrapTasks(ctx)
	if err != nil {
		slog.Error("error bootstrapping tasks", "error", err)
	}

	slog.Error("Server error", "error", app.Listen(":8001"))
}
