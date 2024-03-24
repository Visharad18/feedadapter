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

// AlphaVantageAdapter fetches OHLCV data for a symbol from AlphaVantage Adapter
type AlphaVantageAdapter struct {
	cfg    *config.AlphaVantageConfig
	client *http.Client
}

type alphaVantageResponse struct {
	TimeSeries map[string]TimeSeriesValue `json:"Time Series (5min)"`
}
type TimeSeriesValue struct {
	Open   float64 `json:"1. open"`
	High   float64 `json:"2. high"`
	Low    float64 `json:"3. low"`
	Close  float64 `json:"4. close"`
	Volume float64 `json:"5. volume"`
}

func NewAlphaVantageAdapter(cfg *config.AlphaVantageConfig) *AlphaVantageAdapter {
	return &AlphaVantageAdapter{
		cfg:    cfg,
		client: http.DefaultClient,
	}
}

// Get fetches the OHLCV data for the symbol
// TODO could take interval as enum arg supported by the AlphaVantage API and change interface signatue accorindingly
func (a *AlphaVantageAdapter) Get(_ context.Context, symbol string, _ time.Duration) (map[int64]*entity.HistoricalData, error) {
	url := fmt.Sprintf(a.cfg.BaseURL, symbol, a.cfg.APIKey)
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error in getting data from Alpha Vantage, %s", err)
	}

	defer resp.Body.Close()
	var data alphaVantageResponse
	// err = json.NewDecoder(resp.Body).Decode(&data)
	b, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(b, &data)
	if err != nil {
		return nil, fmt.Errorf("invalid data received from alpha vantage, error: %s, %s", err,
			func() []byte {

				return b
			}(),
		)
	}

	return data.toHistoricalDataEntity()
}

func (resp *alphaVantageResponse) toHistoricalDataEntity() (map[int64]*entity.HistoricalData, error) {
	res := make(map[int64]*entity.HistoricalData)

	for t, hd := range resp.TimeSeries {
		_time, err := time.Parse(time.DateTime, t)
		if err != nil {
			return nil, fmt.Errorf("error parsing timestamp in alpha vantage response:%s", err)
		}
		res[_time.Unix()] = &entity.HistoricalData{
			Open:   hd.Open,
			High:   hd.High,
			Low:    hd.Low,
			Close:  hd.Close,
			Volume: hd.Volume,
		}
	}

	return res, nil
}
