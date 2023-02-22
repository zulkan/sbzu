package repo

import (
	"context"
	"gozu/domain"
	"gozu/utils"
	"time"

	"github.com/go-redis/redis"
)

type stock struct {
	redisClient *redis.Client
}

func (s *stock) WriteStockSummary(_ context.Context, summary *domain.StockSummary) error {
	return utils.WriteRedis(s.redisClient, summary.StockCode, utils.ToJSON(summary), time.Hour) // expired just 1 hour only for development
}

func (s *stock) GetStockSummary(_ context.Context, stockCode string) (*domain.StockSummary, error) {
	dataStr, err := utils.ReadRedis(s.redisClient, stockCode)
	if err != nil {
		return nil, err
	}
	return utils.ReadJSON[domain.StockSummary](dataStr), nil
}

func NewStockRepo(redisClient *redis.Client) domain.StockRepo {
	return &stock{redisClient: redisClient}
}
