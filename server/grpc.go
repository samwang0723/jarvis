package server

import (
	"context"
	"samwang0723/jarvis/dto"
	pb "samwang0723/jarvis/pb"
)

func (s *server) ListDailyClose(ctx context.Context, req *pb.ListDailyCloseRequest) (*pb.ListDailyCloseResponse, error) {
	res, err := s.Handler().ListDailyClose(ctx, dto.ListDailyCloseRequestFromPB(req))
	if err != nil {
		return nil, err
	}

	return dto.ListDailyCloseResponseToPB(res), nil
}
