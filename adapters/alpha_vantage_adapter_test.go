package adapters

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Visharad18/feedadapter/config"
	"github.com/go-playground/assert/v2"
)

var (
	server           *httptest.Server
	cfg              *config.Config
	alphaVantageResp = `{
		"Time Series (Daily)": {
			"2024-03-22 10:00:00": {
				"1. open": "2897.05",
				"2. high": "2920.0",
				"3. low": "2895.35",
				"4. close": "2909.9",
				"5. adjusted close": "2909.9",
				"6. volume": "562484",
				"7. dividend amount": "0.0000",
				"8. split coefficient": "1.0"
			},
			"2024-03-22 10:05:00": {
				"1. open": "2891.4",
				"2. high": "2915.0",
				"3. low": "2889.65",
				"4. close": "2901.3",
				"5. adjusted close": "2901.3",
				"6. volume": "77091",
				"7. dividend amount": "0.0000",
				"8. split coefficient": "1.0"
			}
		}
	}`

	polygonResponse = `{
		"request_id": "6a7e466379af0a71039d60cc78e72282",
		"results": [
		  {
			"c": 75.0875,
			"h": 75.15,
			"l": 73.7975,
			"n": 1,
			"o": 74.06,
			"t": 1577941200000,
			"v": 135647456,
			"vw": 74.6099
		  },
		  {
			"c": 74.3575,
			"h": 75.145,
			"l": 74.125,
			"n": 1,
			"o": 74.2875,
			"t": 1578027600000,
			"v": 146535512,
			"vw": 74.7026
		  }
		]
	}`
)

func TestMain(m *testing.M) {

	// mock server to return mock responses for adapters respectively on the basis of endpoint
	server = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.URL.Path, "alphavantage") {
				w.WriteHeader(http.StatusAccepted)
				w.Write([]byte(alphaVantageResp))
				return
			}
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte(polygonResponse))
		},
	))

	cfg = &config.Config{
		AlphaVantageConfig: config.AlphaVantageConfig{
			BaseURL: server.URL + "/alphavantage/%s/%s",
			APIKey:  "1",
		},
		PolygonConfig: config.PolygonConfig{
			BaseURL: server.URL + "/polygon/%s/%d/%d/%s",
			APIKey:  "2",
		},
	}

	m.Run()
}

// Test for Alpha Vantage adapter verfying if there is no error parsing a dummy alpha vantage response/
// and matching OHLCV values
func TestAlphaVantageAdapter(t *testing.T) {
	a := NewAlphaVantageAdapter(&cfg.AlphaVantageConfig)
	data, err := a.Get(context.TODO(), "AAPL", time.Minute)
	if err != nil {
		t.Errorf("error in AlphaVantageAdapter.Get: %s", err)
	}
	assert.Equal(t, 2, len(data))
	assert.Equal(t, 2897.05, data[1711101600].Open)
	assert.Equal(t, 2920.0, data[1711101600].High)
	assert.Equal(t, 2895.35, data[1711101600].Low)
	assert.Equal(t, 2909.9, data[1711101600].Close)
	assert.Equal(t, 562484, data[1711101600].Volume)
}
