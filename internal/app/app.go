package app

import (
	"backend-store/config"
	"backend-store/internal/handlers"
	"backend-store/internal/service"
	"backend-store/internal/storage"
	"backend-store/pkg/logger"
)

type App struct {
	Config   *config.Config
	Storage  storage.Storage
	Services *Services
	Handlers *Handlers
	log      logger.Log
}

type Services struct {
	ProductService service.ProductService
	OrderService   service.OrderService
}

type Handlers struct {
	ProductHandler *handlers.ProductHandler
	OrderHandler   *handlers.OrderHandler
}

func New(cfg *config.Config, log logger.Log) (*App, error) {
	app := &App{
		Config: cfg,
		log:    log,
	}

	store, err := app.initStorage()
	if err != nil {
		return nil, err
	}
	app.Storage = store
	app.Services = app.initServices()
	app.Handlers = app.initHandlers()

	app.log.Info("Application initialized successfully")
	return app, nil
}

func (a *App) initStorage() (storage.Storage, error) {
	if a.Config.DatabaseURL != "" {
		a.log.Info("Using PostgreSQL storage")
		postgresStore, err := storage.NewPostgresStorage(a.Config.DatabaseURL)
		if err != nil {
			return nil, err
		}

		postgresStore.SetMaxOpenConns(a.Config.MaxOpenConns)
		postgresStore.SetMaxIdleConns(a.Config.MaxIdleConns)
		postgresStore.SetConnMaxLifetime(a.Config.ConnMaxLifetime)

		if err := postgresStore.Init(); err != nil {
			return nil, err
		}

		return postgresStore, nil
	}

	a.log.Info("Using in-memory storage")
	memoryStore := storage.NewMemoryStorage()
	return memoryStore, nil
}

func (a *App) initServices() *Services {
	return &Services{
		ProductService: service.NewProductService(a.Storage),
		OrderService:   service.NewOrderService(a.Storage),
	}
}

func (a *App) initHandlers() *Handlers {
	return &Handlers{
		ProductHandler: handlers.NewProductHandler(a.Services.ProductService),
		OrderHandler:   handlers.NewOrderHandler(a.Services.OrderService),
	}
}

func (a *App) Close() error {
	if a.Storage != nil {
		a.log.Info("Closing application resources")
		return a.Storage.Close()
	}
	return nil
}
