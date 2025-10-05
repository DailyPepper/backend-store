package app

import (
	"backend-store/config"
	"backend-store/internal/handlers"
	"backend-store/internal/service"
	"backend-store/internal/storage"
	"log"
)

type App struct {
	Config   *config.Config
	Storage  storage.Storage
	Services *Services
	Handlers *Handlers
}

type Services struct {
	ProductService service.ProductService
	OrderService   service.OrderService
}

type Handlers struct {
	ProductHandler *handlers.ProductHandler
	OrderHandler   *handlers.OrderHandler
}

func New(cfg *config.Config) (*App, error) {
	app := &App{Config: cfg}

	store, err := initStorage(cfg)
	if err != nil {
		return nil, err
	}
	app.Storage = store
	app.Services = initServices(store)

	app.Handlers = initHandlers(app.Services)

	return app, nil
}

func initStorage(cfg *config.Config) (storage.Storage, error) {
	if cfg.DatabaseURL != "" {
		log.Println("Using PostgreSQL storage")
		postgresStore, err := storage.NewPostgresStorage(cfg.DatabaseURL)
		if err != nil {
			return nil, err
		}

		postgresStore.SetMaxOpenConns(cfg.MaxOpenConns)
		postgresStore.SetMaxIdleConns(cfg.MaxIdleConns)
		postgresStore.SetConnMaxLifetime(cfg.ConnMaxLifetime)

		if err := postgresStore.Init(); err != nil {
			return nil, err
		}

		return postgresStore, nil
	}

	log.Println("Using in-memory storage")
	memoryStore := storage.NewMemoryStorage()
	return memoryStore, nil
}

func initServices(store storage.Storage) *Services {
	return &Services{
		ProductService: service.NewProductService(store),
		OrderService:   service.NewOrderService(store),
	}
}

func initHandlers(services *Services) *Handlers {
	return &Handlers{
		ProductHandler: handlers.NewProductHandler(services.ProductService),
		OrderHandler:   handlers.NewOrderHandler(services.OrderService),
	}
}

func (a *App) Close() error {
	if a.Storage != nil {
		return a.Storage.Close()
	}
	return nil
}
