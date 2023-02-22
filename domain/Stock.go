package domain

import "context"

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
	WriteStockSummary(ctx context.Context, summary *StockSummary) error
	GetStockSummary(ctx context.Context, stockCode string) (*StockSummary, error)
}

type StockUseCase interface {
	QueueProcessor // shows this interface can process kafka message
	ProcessFileData(ctx context.Context, rawData string) error
	WriteStockSummary(ctx context.Context, record *StockRecord) error
	GetStockSummary(ctx context.Context, stockCode string) (*StockSummary, error)
}
