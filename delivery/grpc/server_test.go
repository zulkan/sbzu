package grpc

import (
	"context"
	"gozu/domain"
	"gozu/domain/mocks"
	"testing"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/suite"
)

type ServerTestSuite struct {
	suite.Suite
}

func (s *ServerTestSuite) TestGetStockSummary_NotFound() {
	stockUseCase := mocks.NewStockUseCase(s.T())

	server := &server{stockUseCase: stockUseCase}

	stockUseCase.On("GetStockSummary", context.Background(), "ABC").Return(nil, redis.Nil)
	resp, err := server.GetStockSummary(context.Background(), &StockSummaryRequest{StockCode: "ABC"})

	s.Nil(resp, "response is nil")
	s.Errorf(err, "ABC not found")
}
func (s *ServerTestSuite) TestGetStockSummary_Found() {
	stockUseCase := mocks.NewStockUseCase(s.T())

	stockSummary := &domain.StockSummary{
		StockCode: "ABC",
		Open:      123,
		High:      422,
		Low:       21,
		Close:     21,
		Prev:      11,
		Volume:    321,
		Value:     1333,
		AvgPrice:  1333 / 321,
	}

	server := &server{stockUseCase: stockUseCase}

	stockUseCase.On("GetStockSummary", context.Background(), "ABC").Return(stockSummary, nil)
	resp, err := server.GetStockSummary(context.Background(), &StockSummaryRequest{StockCode: "ABC"})

	s.Equal(StockSummaryToResponse(stockSummary), resp, "response is nil")
	s.Nil(err, "error is nill")
}
func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}
