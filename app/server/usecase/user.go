// Licensed to Dipto Mondal under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Dipto Mondal licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package usecase

import (
	"github.com/diptomondal007/your-money/app/common/response"
	"github.com/diptomondal007/your-money/app/server/model"
	"github.com/diptomondal007/your-money/app/server/repository"
	"log"
	"time"
)

type AddBalanceReq struct {
	TransactionID string  `json:"transaction_id"`
	Amount        float64 `json:"amount"`
}

type AddBalanceResp struct {
	Balance float64 `json:"current_balance"`
}

type CheckBalanceResp struct {
	Balance float64 `json:"balance"`
}

type ListHistory struct {
	Total     int64     `json:"total"`
	PageSize  int64     `json:"page_size"`
	NextPage  string    `json:"next_page"`
	Histories []History `json:"histories"`
}

type History struct {
	CreatedAt     time.Time `json:"created_at"`
	Amount        float64   `json:"amount"`
	TransactionID string    `json:"transaction_id"`
}

// UserUseCase ...
type userUseCase struct {
	repo repository.UserRepository
}

// UserUseCase is interface for user use case
type UserUseCase interface {
	AddBalance(userID string, req *AddBalanceReq) (*AddBalanceResp, error)
	CheckBalance(userID string) (*CheckBalanceResp, error)
	ListHistory(userID string, pageSize int64, cursor string) (*ListHistory, error)
}

// NewUserUseCase returns a new user use case instance
func NewUserUseCase(repo repository.UserRepository) UserUseCase {
	return &userUseCase{repo: repo}
}

func (u *userUseCase) AddBalance(userID string, req *AddBalanceReq) (*AddBalanceResp, error) {
	us, err := u.repo.AddBalance(userID, req.TransactionID, req.Amount)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return toAddBalanceResp(us), nil
}

func (u *userUseCase) CheckBalance(userID string) (*CheckBalanceResp, error) {
	us, err := u.repo.CheckBalance(userID)
	if err != nil {
		return nil, err
	}

	return &CheckBalanceResp{Balance: us.Balance}, nil
}

func (u *userUseCase) ListHistory(userID string, pageSize int64, cursor string) (*ListHistory, error) {
	ts, err := u.repo.GetHistoryList(userID, pageSize, cursor)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	count, err := u.repo.GetHistoryCount(userID)
	if err != nil {
		return nil, err
	}

	var lastID uint

	histories := make([]History, 0)
	for i := range ts {
		histories = append(histories, History{
			CreatedAt:     ts[i].CreatedAt,
			Amount:        ts[i].Amount,
			TransactionID: ts[i].TransactionID,
		})

		if i == len(ts)-1 {
			lastID = ts[i].ID
		}
	}

	nPage := &response.Cursor{ID: lastID}

	return &ListHistory{
		Total:     count,
		PageSize:  pageSize,
		NextPage:  nPage.ToBase64String(),
		Histories: histories,
	}, nil
}

func toAddBalanceResp(info *model.User) *AddBalanceResp {
	return &AddBalanceResp{Balance: info.Balance}
}
