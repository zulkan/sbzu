package domain

//go:generate mockery --name StockRepo
//go:generate mockery --name StockUseCase

type StockRecord struct {
	StockCode string `json:"stock_code"`
	Type      string `json:"type"`
	Quantity  int64  `json:"quantity"`
	Price     int64  `json:"price"`
}
type StockSummary struct {
	StockCode string  `json:"stock_code"`
	Open      int64   `json:"open"`
	High      int64   `json:"high"`
	Low       int64   `json:"low"`
	Close     int64   `json:"close"`
	Prev      int64   `json:"prev"`
	Volume    int64   `json:"volume"`
	Value     int64   `json:"value"`
	AvgPrice  float64 `json:"avg_price"`
}

type StockRepo interface {
	WriteStockSummary(summary *StockSummary) error
	GetStockSummary(stockCode string) (*StockSummary, error)
}

type StockUseCase interface {
	QueueProcessor // shows this interface can process kafka message
	ProcessFileData(rawData string) error
	WriteStockSummary(record *StockRecord) error
	GetStockSummary(stockCode string) (*StockSummary, error)
}
