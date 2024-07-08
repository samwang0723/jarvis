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

package services

import (
	"context"
	"errors"

	"github.com/samwang0723/jarvis/internal/app/domain"
	"github.com/samwang0723/jarvis/internal/app/dto"
	"golang.org/x/crypto/bcrypt"
)

var (
	errUserNotFound         = errors.New("user not found")
	errUserPasswordNotMatch = errors.New("user password not match")
)

func (s *serviceImpl) GetUser(ctx context.Context) (obj *domain.User, err error) {
	obj, err = s.dal.GetUserByID(ctx, s.currentUserID)
	if err != nil {
		return nil, err
	}

	if obj == nil {
		return nil, errUserNotFound
	}

	return obj, nil
}

func (s *serviceImpl) CreateUser(ctx context.Context, obj *domain.User) (err error) {
	// Hash the password with a default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(obj.Password), bcrypt.DefaultCost)
	obj.Password = string(hashedPassword)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to hash password")

		return err
	}

	err = s.dal.CreateUser(ctx, obj)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to create user")

		return err
	}

	return nil
}

func (s *serviceImpl) UpdateUser(ctx context.Context, obj *domain.User) (err error) {
	err = s.dal.UpdateUser(ctx, obj)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to update user")

		return err
	}

	return nil
}

func (s *serviceImpl) DeleteUser(ctx context.Context) (err error) {
	err = s.dal.DeleteUserByID(ctx, s.currentUserID)
	if err != nil {
		return err
	}

	return nil
}

func (s *serviceImpl) ListUsers(
	ctx context.Context,
	req *dto.ListUsersRequest,
) (objs []*domain.User, totalCount int64, err error) {
	objs, totalCount, err = s.dal.ListUsers(ctx, req.Offset, req.Limit)
	if err != nil {
		return nil, 0, err
	}

	return objs, totalCount, nil
}

func (s *serviceImpl) GetUserByEmail(
	ctx context.Context,
	email string,
) (obj *domain.User, err error) {
	obj, err = s.dal.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if obj == nil {
		return nil, errUserNotFound
	}

	return obj, nil
}

func (s *serviceImpl) GetUserByPhone(
	ctx context.Context,
	phone string,
) (obj *domain.User, err error) {
	obj, err = s.dal.GetUserByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}

	if obj == nil {
		return nil, errUserNotFound
	}

	return obj, nil
}
