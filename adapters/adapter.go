package adapters

import (
	"context"
	"time"

	"github.com/Visharad18/feedadapter/entity"
)

// Adapter: type for fetching stock prices
// method Get: takes symbol name and time duration for which prices have to be fetched
// and returns stock price struct
type Adapter interface {
	Get(context.Context, string, time.Duration) (map[int64]*entity.HistoricalData, error)
}
