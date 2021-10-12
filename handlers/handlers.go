package handlers

import (
	"context"
	"samwang0723/jarvis/dto"
	"samwang0723/jarvis/services"
)

type IHandler interface {
	BatchingDownload(ctx context.Context, req *dto.DownloadRequest)
	ListDailyClose(ctx context.Context, req *dto.ListDailyCloseRequest) (*dto.ListDailyCloseResponse, error)
}

type handlerImpl struct {
	dataService services.IService
}

func New(dataService services.IService) IHandler {
	res := &handlerImpl{
		dataService: dataService,
	}
	return res
}
