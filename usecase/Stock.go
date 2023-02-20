package usecase

import (
	"github.com/go-redis/redis"
	"gozu/domain"
	"gozu/utils"
)

type stockUseCase struct {
	queueUseCase domain.QueueUseCase
	stockRepo    domain.StockRepo
}

func (s *stockUseCase) ProcessQueueMessage(rawData string) error {
	stockRecord := utils.ReadJSON[domain.StockRecord](rawData)

	return s.WriteStockSummary(stockRecord)
}

// ProcessFileData publish to kafka from input from file
func (s *stockUseCase) ProcessFileData(rawData string) error {
	stockRecord := utils.ReadJSON[domain.StockRecord](rawData)

	return s.queueUseCase.PublishMessage(stockRecord.StockCode, rawData)
}

func initiateStockSummary(record *domain.StockRecord) *domain.StockSummary {
	var stockSummary *domain.StockSummary
	if record.Type == "A" {
		stockSummary = &domain.StockSummary{
			StockCode: record.StockCode,
			Open:      0,
			High:      0,
			Low:       0,
			Close:     0,
			Prev:      record.Price,
		}
	} else {
		stockSummary = &domain.StockSummary{
			StockCode: record.StockCode,
			Open:      record.Price,
			High:      record.Price,
			Low:       record.Price,
			Close:     record.Price,
			Volume:    record.Quantity,
			Value:     record.Quantity * record.Price,
			AvgPrice:  float64(record.Price),
		}
	}
	return stockSummary
}

func updateStockSummary(stockSummary *domain.StockSummary, record *domain.StockRecord) {
	if record.Type == "A" {
		stockSummary.Prev = record.Price
	} else {
		if stockSummary.Open == 0 {
			stockSummary.Open = record.Price
		}
		stockSummary.Close = record.Price
		if record.Price < stockSummary.Low {
			stockSummary.Low = record.Price
		}
		if record.Price > stockSummary.High {
			stockSummary.High = record.Price
		}
		stockSummary.Volume += record.Quantity
		stockSummary.Value += record.Quantity * record.Price
		stockSummary.AvgPrice = float64(stockSummary.Value) / float64(stockSummary.Volume)
	}
}

func (s *stockUseCase) WriteStockSummary(record *domain.StockRecord) error {
	stockSummary, err := s.GetStockSummary(record.StockCode)

	if err != nil && err != redis.Nil {
		return err
	}
	if err == redis.Nil {
		stockSummary = initiateStockSummary(record)
	} else {
		updateStockSummary(stockSummary, record)
	}
	return s.stockRepo.WriteStockSummary(stockSummary)
}

func (s *stockUseCase) GetStockSummary(stockCode string) (*domain.StockSummary, error) {
	return s.stockRepo.GetStockSummary(stockCode)
}

func NewStockUseCase(queueUseCase domain.QueueUseCase, stockRepo domain.StockRepo) domain.StockUseCase {
	return &stockUseCase{queueUseCase: queueUseCase, stockRepo: stockRepo}
}
