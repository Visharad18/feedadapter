package adapters

import (
	"context"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
)

func TestPolygonAdapter(t *testing.T) {
	a := NewPolygonAdapter(&cfg.PolygonConfig)
	data, err := a.Get(context.TODO(), "AAPL", time.Minute)
	if err != nil {
		t.Errorf("error in PolygonAdapter.Get: %s", err)
	}
	assert.Equal(t, 2, len(data))
	assert.Equal(t, 75.0875, data[1577941200000].Close)
	assert.Equal(t, 75.15, data[1577941200000].High)
	assert.Equal(t, 73.7975, data[1577941200000].Low)
	assert.Equal(t, 74.06, data[1577941200000].Open)
	assert.Equal(t, 135647456, data[1577941200000].Volume)
}
