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

package idal

import (
	"context"

	"github.com/samwang0723/jarvis/internal/app/entity"
)

type IUserDAL interface {
	CreateUser(ctx context.Context, obj *entity.User) error
	UpdateUser(ctx context.Context, obj *entity.User) error
	DeleteUserByID(ctx context.Context, id uint64) error
	DeleteSessionID(ctx context.Context, userID uint64) error
	GetUserByID(ctx context.Context, id uint64) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByPhone(ctx context.Context, phone string) (*entity.User, error)
	ListUsers(ctx context.Context, offset, limit int32) ([]*entity.User, int64, error)
	UpdateSessionID(ctx context.Context, obj *entity.User) error
}
