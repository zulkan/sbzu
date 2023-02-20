package repo

import (
	"github.com/go-redis/redis"
	"gozu/domain"
	"gozu/utils"
	"time"
)

type stock struct {
	redisClient *redis.Client
}

func (s *stock) WriteStockSummary(summary *domain.StockSummary) error {
	return utils.WriteRedis(s.redisClient, summary.StockCode, utils.ToJSON(summary), time.Hour) // expired just 1 hour only for development
}

func (s *stock) GetStockSummary(stockCode string) (*domain.StockSummary, error) {
	dataStr, err := utils.ReadRedis(s.redisClient, stockCode)
	if err != nil {
		return nil, err
	}
	return utils.ReadJSON[domain.StockSummary](dataStr), nil
}

func NewStockRepo(redisClient *redis.Client) domain.StockRepo {
	return &stock{redisClient: redisClient}
}
