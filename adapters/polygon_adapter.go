package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Visharad18/feedadapter/config"
	"github.com/Visharad18/feedadapter/entity"
)

type PolygonAdapter struct {
	cfg    *config.PolygonConfig
	client *http.Client
}

type PolygonResponse struct {
	Results []Result `json:"results"`
}

type Result struct {
	Open   float64 `json:"o"`
	High   float64 `json:"h"`
	Low    float64 `json:"l"`
	Close  float64 `json:"c"`
	Volume float64 `json:"v"`
	Time   int64   `json:"t"`
}

func NewPolygonAdapter(cfg *config.PolygonConfig) *PolygonAdapter {
	return &PolygonAdapter{
		cfg:    cfg,
		client: http.DefaultClient,
	}
}

// Get fetches the OHLCV data for the symbol
// TODO could take interval as enum arg supported by the Polygon API and change interface signatue accorindingly
func (a *PolygonAdapter) Get(_ context.Context, symbol string, duration time.Duration) (map[int64]*entity.HistoricalData, error) {
	start := time.Now().Add(-24 * time.Hour)
	end := start.Add(duration)
	url := fmt.Sprintf(a.cfg.BaseURL, symbol, start.UnixMilli(), end.UnixMilli(), a.cfg.APIKey)
	fmt.Println(url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error in getting data from Polygon, %s", err)
	}

	defer resp.Body.Close()
	var data PolygonResponse
	// err = json.NewDecoder(resp.Body).Decode(&data)
	b, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(b, &data)
	if err != nil {
		return nil, fmt.Errorf("invalid data received from Polygon, error: %s %s", err, b)
	}

	return data.toHistoricalDataEntity()
}

func (resp *PolygonResponse) toHistoricalDataEntity() (map[int64]*entity.HistoricalData, error) {
	res := make(map[int64]*entity.HistoricalData)

	for _, r := range resp.Results {
		res[r.Time/1000] = &entity.HistoricalData{
			Open:   r.Open,
			High:   r.High,
			Low:    r.Low,
			Close:  r.Close,
			Volume: r.Volume,
		}
	}

	return res, nil
}
