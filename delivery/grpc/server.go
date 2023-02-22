package grpc

import (
	"context"
	"fmt"
	"gozu/domain"
	"log"
	"net"

	"github.com/go-redis/redis"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	UnimplementedStockServiceServer
	stockUseCase domain.StockUseCase
}

func StockSummaryToResponse(stockSummary *domain.StockSummary) *StockSummaryResponse {
	return &StockSummaryResponse{
		StockCode: stockSummary.StockCode,
		Open:      stockSummary.Open,
		High:      stockSummary.High,
		Low:       stockSummary.Low,
		Close:     stockSummary.Close,
		Prev:      stockSummary.Prev,
		Volume:    stockSummary.Volume,
		Value:     stockSummary.Value,
		AvgPrice:  int64(stockSummary.AvgPrice),
	}
}

func (s *server) GetStockSummary(ctx context.Context, req *StockSummaryRequest) (*StockSummaryResponse, error) {
	stockSummary, err := s.stockUseCase.GetStockSummary(ctx, req.GetStockCode())

	if err == redis.Nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("%s not found", req.StockCode))
	}
	if err != nil {
		return nil, err
	}

	return StockSummaryToResponse(stockSummary), nil
}
func NewGrpcServer(stockUseCase domain.StockUseCase) {
	port := 50051
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("grpc server started at", port)

	grpcServer := grpc.NewServer()
	RegisterStockServiceServer(grpcServer, &server{stockUseCase: stockUseCase})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
