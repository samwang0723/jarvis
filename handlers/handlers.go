package handlers

import (
	"context"
	"samwang0723/jarvis/services"
)

type IHandler interface {
	BatchingDownload(ctx context.Context, rewindLimit int, rateLimit int)
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
