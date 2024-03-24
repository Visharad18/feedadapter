package app

import (
	"context"
	"time"

	"github.com/Visharad18/feedadapter/adapters"
	"github.com/Visharad18/feedadapter/cache"
	"github.com/Visharad18/feedadapter/config"
	"github.com/sirupsen/logrus"
)

type App struct {
	cfg      *config.Config
	lg       *logrus.Logger
	cache    cache.Cache
	adapters []adapters.Adapter
}

func NewApp(cfg *config.Config, lg *logrus.Logger) *App {
	return &App{
		cfg:   cfg,
		lg:    lg,
		cache: cache.NewInMeInMemoryCache(),
		adapters: []adapters.Adapter{
			adapters.NewAlphaVantageAdapter(&cfg.AlphaVantageConfig),
			adapters.NewPolygonAdapter(&cfg.PolygonConfig),
		},
	}
}

// Run infinitely polls for data after every config.FetchInterval
func (a *App) Run(ctx context.Context) error {
	for {
		go a.sync(ctx)
		time.Sleep(a.cfg.FetchInterval)
	}
}

// sync fetches data for all symbols from all adapters
func (a *App) sync(ctx context.Context) {
	for _, ad := range a.adapters {
		for _, s := range a.cfg.Symbols {
			go func(adapter adapters.Adapter, symbol string) {
				_ctx, cancel := context.WithTimeout(ctx, time.Minute)
				defer cancel()
				res, err := adapter.Get(_ctx, symbol, a.cfg.FetchInterval)
				if err != nil {
					a.lg.Errorf("error in fetching data in adapter: %s", err)
				}

				if err := a.cache.Store(symbol, res); err != nil {
					a.lg.Errorf("error in storing data in cache: %s", err)
				}
			}(ad, s)
		}
	}
}

// GetData returns all the data in cache
func (a *App) GetData() map[string]any {
	return a.cache.GetAll()
}
