package service

import (
	"context"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dao/speedtestdao"
	"github.com/muchlist/risa_restfull/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func NewSpeedTestService(speedDao speedtestdao.SpedTestDaoAssumer) SpeedTestServiceAssumer {
	return &speedService{
		sd: speedDao,
	}
}

type speedService struct {
	sd speedtestdao.SpedTestDaoAssumer
}
type SpeedTestServiceAssumer interface {
	InsertSpeed(ctx context.Context, input dto.SpeedTest) (*string, rest_err.APIError)
	RetrieveSpeed(ctx context.Context) (dto.SpeedTestList, rest_err.APIError)
}

func (s *speedService) InsertSpeed(ctx context.Context, input dto.SpeedTest) (*string, rest_err.APIError) {
	timeNow := time.Now().Unix()

	data := dto.SpeedTest{
		ID:        primitive.ObjectID{},
		Time:      timeNow,
		LatencyMs: input.LatencyMs,
		Upload:    input.Upload,
		Download:  input.Download,
	}
	// DB
	insertedID, err := s.sd.InsertSpeed(ctx, data)
	if err != nil {
		return nil, err
	}
	return insertedID, nil
}

func (s *speedService) RetrieveSpeed(ctx context.Context) (dto.SpeedTestList, rest_err.APIError) {
	otherList, err := s.sd.RetrieveSpeed(ctx)
	return otherList, err
}
