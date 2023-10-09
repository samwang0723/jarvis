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

	"github.com/samwang0723/jarvis/internal/app/dto"
	"github.com/samwang0723/jarvis/internal/app/entity"
)

func (s *serviceImpl) CreateTransactions(ctx context.Context, objs []*entity.Transaction) error {
	err := s.dal.CreateTransactions(ctx, objs)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to create transaction")

		return err
	}

	return nil
}

func (s *serviceImpl) GetTransactionByID(ctx context.Context, id uint64) (*entity.Transaction, error) {
	transaction, err := s.dal.GetTransactionByID(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get transaction by id")

		return nil, err
	}

	return transaction, nil
}

func (s *serviceImpl) ListTransactions(
	ctx context.Context,
	req *dto.ListTransactionsRequest,
) (objs []*entity.Transaction, totalCount int64, err error) {
	transactions, totalCount, err := s.dal.ListTransactions(ctx, req.UserID, req.Limit, req.Offset)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to list transactions")

		return nil, 0, err
	}

	return transactions, totalCount, nil
}
