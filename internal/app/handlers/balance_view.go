// Copyright 2021 Wei (Sam) Wang <sam.wang.0723@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package handlers

import (
	"context"

	"github.com/samwang0723/jarvis/internal/app/dto"
	"github.com/samwang0723/jarvis/internal/app/entity"
)

func (h *handlerImpl) UpdateBalanceView(
	ctx context.Context,
	req *dto.UpdateBalanceViewRequest,
) (*dto.UpdateBalanceViewResponse, error) {
	balanceView := &entity.BalanceView{
		UserID:        req.UserID,
		CurrentAmount: req.CurrentAmount,
	}

	err := h.dataService.UpdateBalanceView(ctx, balanceView)
	if err != nil {
		return &dto.UpdateBalanceViewResponse{
			Status:       dto.StatusError,
			ErrorCode:    "",
			ErrorMessage: err.Error(),
			Success:      false,
		}, err
	}

	return &dto.UpdateBalanceViewResponse{
		Status:       dto.StatusSuccess,
		ErrorCode:    "",
		ErrorMessage: "",
		Success:      true,
	}, nil
}

func (h *handlerImpl) GetBalanceViewByUserID(ctx context.Context, userID uint64) (*entity.BalanceView, error) {
	balanceView, err := h.dataService.GetBalanceViewByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return balanceView, nil
}
