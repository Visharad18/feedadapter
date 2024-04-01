package entity

// HistoricalData: entity layer type, acts as common base for all adapters
type HistoricalData struct {
	Open   float64
	Close  float64
	Low    float64
	High   float64
	Volume float64
}
