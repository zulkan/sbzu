package usecase

import (
	"context"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gozu/domain"
	"gozu/domain/mocks"
	"gozu/utils"
	"testing"
)

type StockTestSuite struct {
	suite.Suite
}

func (s *StockTestSuite) TestInitiateStockSummary() {
	record := &domain.StockRecord{
		StockCode: "BBAC",
		Type:      "P",
		Quantity:  23,
		Price:     300,
	}
	summary := initiateStockSummary(record)

	s.Equalf(summary.Value, record.Price*record.Quantity, "the value should be price * quantity")
	s.Equalf(summary.Open, record.Price, "the open value should be price")
	s.Equalf(summary.Low, record.Price, "the low value should be price")
	s.Equalf(summary.High, record.Price, "the high value should be price")
	s.Equalf(summary.Close, record.Price, "the close value should be price")
	s.Equalf(summary.Volume, record.Quantity, "the volume price value should be quantity")
	s.Equalf(summary.AvgPrice, float64(record.Price), "the avg price value should be price")
	s.Equalf(summary.StockCode, "BBAC", "stock code should be BBAC")
}
func (s *StockTestSuite) TestInitiateStockSummary_Type_A() {
	record := &domain.StockRecord{
		StockCode: "BBAC",
		Type:      "A",
		Quantity:  23,
		Price:     300,
	}
	summary := initiateStockSummary(record)

	s.Equalf(summary.Value, int64(0), "the value should be 0")
	s.Equalf(summary.Open, int64(0), "the open value should be 0")
	s.Equalf(summary.Prev, record.Price, "the prev value should be equal price")
	s.Equalf(summary.Close, int64(0), "the close value should be 0")
	s.Equalf(summary.Volume, int64(0), "the volume value should be 0")
	s.Equalf(summary.AvgPrice, float64(0), "the avg price value should be 0")
	s.Equalf(summary.StockCode, "BBAC", "stock code should be BBAC")
}
func (s *StockTestSuite) TestUpdateStockSummary() {
	summary := &domain.StockSummary{
		StockCode: "ABC",
		Open:      3000,
		High:      5000,
		Low:       3200,
		Close:     4000,
		Prev:      1500,
		Volume:    600,
		Value:     1000000,
		AvgPrice:  1667,
	}
	record := &domain.StockRecord{
		StockCode: "ABC",
		Type:      "P",
		Quantity:  23,
		Price:     300,
	}
	updateStockSummary(summary, record)
	s.Equalf(summary.Value, int64(1006900), "the value should be 0")
	s.Equalf(summary.Open, int64(3000), "the open value should be 0")
	s.Equalf(summary.Prev, int64(1500), "the prev value should be equal price")
	s.Equalf(summary.Close, int64(300), "the close value should be 0")
	s.Equalf(summary.Volume, int64(623), "the volume value should be 0")
	s.Greaterf(summary.AvgPrice, float64(1616), "the avg price value should above 1616")
	s.Lessf(summary.AvgPrice, float64(1617), "the avg price value should less than 1617")
	s.Equalf(summary.StockCode, "ABC", "stock code should be BBAC")
}

func (s *StockTestSuite) TestProcessFileData_Type_A_NoPublish() {
	lineRawData :=
		`
{"type":"A","order_book":"35","price":"4540","stock_code":"UNVR"}
		`
	mockQueue := mocks.NewQueueUseCase(s.T())

	useCase := NewStockUseCase(mockQueue, nil)

	err := useCase.ProcessFileData(context.Background(), lineRawData)

	mockQueue.AssertNotCalled(s.T(), "PublishMessage", mock.Anything, mock.Anything)
	s.Nil(err)
}
func (s *StockTestSuite) TestProcessFileData_Type_A_Publish() {
	lineRawData :=
		`
{"type":"A","order_number":"202211100000000701","order_verb":"S","quantity":"0","order_book":"552","price":"4240","stock_code":"TLKM"}
		`

	expectedRecord := domain.StockRecord{
		StockCode: "TLKM",
		Type:      "A",
		Quantity:  0,
		Price:     4240,
	}
	mockQueue := mocks.NewQueueUseCase(s.T())

	useCase := NewStockUseCase(mockQueue, nil)
	mockQueue.On("PublishMessage", context.Background(), "TLKM", utils.ToJSON(expectedRecord)).Return(nil)

	err := useCase.ProcessFileData(context.Background(), lineRawData)

	mockQueue.AssertNumberOfCalls(s.T(), "PublishMessage", 1)
	s.Nil(err)
}

func (s *StockTestSuite) TestProcessFileData_Type_P_E_Publish() {
	lineRawData :=
		`
{"type":"E","order_number":"202211100000019752","executed_quantity":"153","execution_price":"4560","order_verb":"B","stock_code":"BBRI","order_book":"911"}
		`

	expectedRecord := domain.StockRecord{
		StockCode: "BBRI",
		Type:      "E",
		Quantity:  153,
		Price:     4560,
	}
	mockQueue := mocks.NewQueueUseCase(s.T())

	useCase := NewStockUseCase(mockQueue, nil)
	mockQueue.On("PublishMessage", context.Background(), expectedRecord.StockCode, utils.ToJSON(expectedRecord)).Return(nil)

	err := useCase.ProcessFileData(context.Background(), lineRawData)

	mockQueue.AssertNumberOfCalls(s.T(), "PublishMessage", 1)
	s.Nil(err)
}

func (s *StockTestSuite) TestWriteStockSummary_init() {
	record := &domain.StockRecord{
		StockCode: "BBRI",
		Type:      "E",
		Quantity:  153,
		Price:     4560,
	}
	mockQueue := mocks.NewQueueUseCase(s.T())
	mockRepo := mocks.NewStockRepo(s.T())
	summary := initiateStockSummary(record)

	useCase := NewStockUseCase(mockQueue, mockRepo)

	mockRepo.On("GetStockSummary", context.Background(), record.StockCode).Return(nil, redis.Nil)
	mockRepo.On("WriteStockSummary", context.Background(), summary).Return(nil)

	err := useCase.WriteStockSummary(context.Background(), record)

	mockRepo.AssertNumberOfCalls(s.T(), "GetStockSummary", 1)
	mockRepo.AssertNumberOfCalls(s.T(), "WriteStockSummary", 1)
	s.Nil(err)
}

func (s *StockTestSuite) TestWriteStockSummary_updateStock() {
	record := &domain.StockRecord{
		StockCode: "BBRI",
		Type:      "E",
		Quantity:  153,
		Price:     4560,
	}
	mockQueue := mocks.NewQueueUseCase(s.T())
	mockRepo := mocks.NewStockRepo(s.T())
	oldSummary := &domain.StockSummary{
		StockCode: "ABC",
		Open:      3000,
		High:      5000,
		Low:       3200,
		Close:     4000,
		Prev:      1500,
		Volume:    600,
		Value:     1000000,
		AvgPrice:  1667,
	}
	newSummary := &domain.StockSummary{
		StockCode: "ABC",
		Open:      3000,
		High:      5000,
		Low:       3200,
		Close:     4000,
		Prev:      1500,
		Volume:    600,
		Value:     1000000,
		AvgPrice:  1667,
	}
	updateStockSummary(newSummary, record)

	useCase := NewStockUseCase(mockQueue, mockRepo)

	mockRepo.On("GetStockSummary", context.Background(), record.StockCode).Return(oldSummary, nil)
	mockRepo.On("WriteStockSummary", context.Background(), newSummary).Return(nil)

	err := useCase.WriteStockSummary(context.Background(), record)

	mockRepo.AssertNumberOfCalls(s.T(), "GetStockSummary", 1)
	mockRepo.AssertNumberOfCalls(s.T(), "WriteStockSummary", 1)
	s.Nil(err)
}

func TestStockTestSuite(t *testing.T) {
	suite.Run(t, new(StockTestSuite))
}
