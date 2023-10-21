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

func (h *handlerImpl) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.CreateUserResponse, error) {
	user := &entity.User{
		Email:     req.Email,
		Phone:     req.Phone,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  req.Password,
	}

	err := h.dataService.CreateUser(ctx, user)
	if err != nil {
		return &dto.CreateUserResponse{
			Status:       dto.StatusError,
			ErrorCode:    "",
			ErrorMessage: err.Error(),
			Success:      false,
		}, err
	}

	return &dto.CreateUserResponse{
		Status:       dto.StatusSuccess,
		ErrorCode:    "",
		ErrorMessage: "",
		Success:      true,
	}, nil
}

func (h *handlerImpl) ListUsers(ctx context.Context, req *dto.ListUsersRequest) (*dto.ListUsersResponse, error) {
	entries, totalCount, err := h.dataService.WithUserID(ctx).ListUsers(ctx, req)
	if err != nil {
		return nil, err
	}

	return &dto.ListUsersResponse{
		Offset:     req.Offset,
		Limit:      req.Limit,
		Entries:    entries,
		TotalCount: totalCount,
	}, nil
}
