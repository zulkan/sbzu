package usecase

import (
	"context"
	"gozu/domain"
	"gozu/utils"
	"log"
	"strings"

	"github.com/go-redis/redis"
)

type stockUseCase struct {
	queueUseCase domain.QueueUseCase
	stockRepo    domain.StockRepo
}

func (s *stockUseCase) ProcessQueueMessage(ctx context.Context, rawData string) error {
	stockRecord := utils.ReadJSON[domain.StockRecord](rawData)
	// To debug process queue
	// log.Println("Process from", rawData, "to", utils.ToJSON(stockRecord))

	return s.WriteStockSummary(ctx, stockRecord)
}

// ProcessFileData publish to kafka from input from file
func (s *stockUseCase) ProcessFileData(ctx context.Context, rawData string) error {
	dataMap := utils.ReadJSON[map[string]interface{}](rawData)
	dataType := (*dataMap)["type"]
	record := domain.StockRecord{}
	record.Type = dataType.(string)

	if dataType == "A" {
		if (*dataMap)["quantity"] == "0" {
			record.Quantity = 0
			record.Price = utils.ToInt64((*dataMap)["price"].(string))
			record.StockCode = (*dataMap)["stock_code"].(string)
		}
	}
	if dataType == "E" || dataType == "P" {
		record.Quantity = utils.ToInt64((*dataMap)["executed_quantity"].(string))
		record.Price = utils.ToInt64((*dataMap)["execution_price"].(string))
		record.StockCode = (*dataMap)["stock_code"].(string)
	}

	// the data is valid, initiated from one of condition above
	if record.Price != 0 {
		return s.queueUseCase.PublishMessage(ctx, record.StockCode, utils.ToJSON(record))
	}
	return nil
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

func (s *stockUseCase) WriteStockSummary(ctx context.Context, record *domain.StockRecord) error {
	if record == nil || strings.TrimSpace(record.StockCode) == "" {
		return nil
	}
	stockSummary, err := s.GetStockSummary(ctx, record.StockCode)

	if err != nil && err != redis.Nil {
		log.Println("Failed when process", record.StockCode, err)
		return err
	}
	if err == redis.Nil {
		stockSummary = initiateStockSummary(record)
	} else {
		updateStockSummary(stockSummary, record)
	}
	return s.stockRepo.WriteStockSummary(ctx, stockSummary)
}

func (s *stockUseCase) GetStockSummary(ctx context.Context, stockCode string) (*domain.StockSummary, error) {
	return s.stockRepo.GetStockSummary(ctx, stockCode)
}

func NewStockUseCase(queueUseCase domain.QueueUseCase, stockRepo domain.StockRepo) domain.StockUseCase {
	return &stockUseCase{queueUseCase: queueUseCase, stockRepo: stockRepo}
}
