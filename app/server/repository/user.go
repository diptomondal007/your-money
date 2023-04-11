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

package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"

	"github.com/diptomondal007/your-money/app/server/model"
	"github.com/diptomondal007/your-money/app/utils/response"
)

// userRepository ...
type userRepository struct {
	db *sqlx.DB
}

// UserRepository ...
type UserRepository interface {
	AddBalance(userID string, transactionID string, amount float64) (*model.User, error)
	GetUserInfo(userID string) (*model.User, error)
	GetHistoryList(userID string, pageSize int64, cursor string) ([]*model.Transaction, error)
	GetHistoryCount(userID string) (int64, error)
}

// NewUserRepo returns a new user repo instance
func NewUserRepo(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (u userRepository) AddBalance(userID string, transactionID string, amount float64) (*model.User, error) {
	updatedUser := &model.User{}

	tx := u.db.MustBegin()
	defer tx.Rollback()

	var user model.User
	q, _, err := goqu.From(goqu.T(model.TableUsers).As("u")).
		Select("u.*").
		Where(goqu.Ex{"id": goqu.Op{"eq": userID}}).
		ForUpdate(exp.Wait). // add for update clause to lock the row so that other process or connection can't read the dirty value
		ToSQL()
	if err != nil {
		return nil, err
	}

	if err := tx.Get(&user, q); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, response.WrapError(fmt.Errorf("user not found"), http.StatusNotFound, "")
		}
		return nil, err
	}

	tr := model.Transaction{}

	q, _, err = goqu.From(goqu.T(model.TableTransactions).As("t")).
		Select("t.*").
		Where(goqu.Ex{"transaction_id": goqu.Op{"eq": transactionID}}).
		ToSQL()
	if err != nil {
		return nil, err
	}

	if err = tx.Get(&tr, q); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	if err == nil {
		return nil, response.WrapError(fmt.Errorf("transaction was already processed"), http.StatusUnprocessableEntity, "")
	}

	q, _, err = goqu.Update(model.TableUsers).
		Set(map[string]interface{}{"balance": goqu.L("balance + ?", amount)}).
		Where(goqu.Ex{"id": goqu.Op{"eq": userID}}).ToSQL()
	if err != nil {
		return nil, err
	}
	tx.MustExec(q)

	q, _, err = goqu.Insert(goqu.T(model.TableTransactions)).Rows(model.Transaction{
		CreatedAt:     time.Now().UTC(),
		Amount:        amount,
		UserID:        userID,
		TransactionID: transactionID,
	}).ToSQL()
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(q)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	q, _, err = goqu.
		From(goqu.T(model.TableUsers).As("u")).
		Select("u.*").
		Where(goqu.Ex{"id": goqu.Op{"eq": userID}}).ToSQL()
	if err != nil {
		return nil, err
	}

	if err = tx.Get(updatedUser, q); err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

// GetUserInfo fetches user info for a user from db
func (u userRepository) GetUserInfo(userID string) (*model.User, error) {
	user := &model.User{}

	q, _, err := goqu.From(goqu.T(model.TableUsers).As("u")).
		Select("u.*").
		Where(goqu.Ex{"id": goqu.Op{"eq": userID}}).
		ToSQL()
	if err != nil {
		return nil, err
	}

	if err := u.db.Get(user, q); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, response.WrapError(fmt.Errorf("user not found"), http.StatusNotFound, "")
		}
		return nil, err
	}
	return user, nil
}

// GetHistoryList returns the transaction history list for a user
func (u userRepository) GetHistoryList(userID string, pageSize int64, cursor string) ([]*model.Transaction, error) {
	res := make([]*model.Transaction, 0)

	d := goqu.From(goqu.T(model.TableTransactions).As("t")).
		Select("t.*")

	d = d.Where(goqu.Ex{"user_id": goqu.Op{"eq": userID}})

	if cursor != "" {
		c, err := response.ParseCursor(cursor)
		if err != nil {
			return res, response.WrapError(fmt.Errorf("invalid pagination cursor"), http.StatusBadRequest, "")
		}

		d = d.Where(goqu.Ex{
			"id": goqu.Op{"lt": c.ID},
		})
	}

	d = d.Order(goqu.I("t.id").Desc())
	d = d.Limit(uint(pageSize))

	sql, _, err := d.ToSQL()
	if err != nil {
		return nil, err
	}

	if err = u.db.Select(&res, sql); err != nil {
		return nil, err
	}
	return res, nil
}

func (u userRepository) GetHistoryCount(userID string) (int64, error) {
	var count int64

	d := goqu.From(goqu.T(model.TableTransactions).As("t")).
		Select(goqu.COUNT("*"))

	d = d.Where(goqu.Ex{"user_id": goqu.Op{"eq": userID}})

	sql, _, err := d.ToSQL()
	if err != nil {
		return 0, err
	}

	if err = u.db.Get(&count, sql); err != nil {
		return 0, err
	}
	return count, nil
}
