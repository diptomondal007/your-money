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

package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/diptomondal007/your-money/app/server/repository"
	"github.com/diptomondal007/your-money/app/server/usecase"
	"github.com/diptomondal007/your-money/app/utils"
	"github.com/diptomondal007/your-money/app/utils/response"
)

func TestAddBalanceBadRequest(t *testing.T) {
	s := echo.New()

	h, mock, err := newTest(s)
	if err != nil {
		panic(err)
	}

	reqBody := `{"amount": 0}`

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	c := s.NewContext(req, rec)
	c.SetPath("/users/:uid/add")

	// params
	c.SetParamNames("uid")
	c.SetParamValues("6d7750a1-c3f2-4765-bf8f-33bc80f3f809")

	mock.ExpectBegin()

	if assert.NoError(t, h.addBalance(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestAddBalanceSuccessful(t *testing.T) {
	s := echo.New()

	res := `{"success":true,"message":"transaction successful!","status_code":202,"data":{"current_balance":110.1}}`
	purchaseBody := `{
    					"amount": 10,
						"transaction_id": "tx_1as4ndakda"
					}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(purchaseBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	c := s.NewContext(req, rec)
	c.SetPath("/users/:uid/add")
	// params
	c.SetParamNames("uid")
	c.SetParamValues("6d7750a1-c3f2-4765-bf8f-33bc80f3f809")

	db, mock := utils.MockSqlxDB()
	ur := repository.NewUserRepo(db)

	us := usecase.NewUserUseCase(ur)

	uRows := sqlmock.NewRows([]string{"id", "name", "balance"}).AddRow("6d7750a1-c3f2-4765-bf8f-33bc80f3f809", "Test", 100.10)

	mock.ExpectBegin()
	query := `SELECT "u".* FROM "users" AS "u" WHERE ("id" = '6d7750a1-c3f2-4765-bf8f-33bc80f3f809')`
	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(uRows)

	query = `SELECT "t".* FROM "transactions" AS "t" WHERE ("transaction_id" = 'tx_1as4ndakda')`
	//uRows = sqlmock.NewRows([]string{"id", "created_at", "amount", "transaction_id"}).AddRow(1, time.Now().UTC(), 100.10, "tx_1as4ndakda")
	//mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(uRows)
	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnError(sql.ErrNoRows)

	query = `UPDATE "users" SET "balance"=balance + 10 WHERE ("id" = '6d7750a1-c3f2-4765-bf8f-33bc80f3f809')`
	mock.ExpectExec(regexp.QuoteMeta(query)).WillReturnResult(sqlmock.NewResult(0, 1))

	query = `INSERT INTO "transactions" ("amount", "created_at", "transaction_id", "user_id")`
	mock.ExpectExec(regexp.QuoteMeta(query)).WillReturnResult(sqlmock.NewErrorResult(nil))

	query = `SELECT "u".* FROM "users" AS "u" WHERE ("id" = '6d7750a1-c3f2-4765-bf8f-33bc80f3f809')`
	uRows = sqlmock.NewRows([]string{"id", "name", "balance"}).AddRow("6d7750a1-c3f2-4765-bf8f-33bc80f3f809", "Test", 110.10)
	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(uRows)

	mock.ExpectCommit()

	h := NewHandler(s, us)

	if assert.NoError(t, h.addBalance(c)) {
		assert.Equal(t, http.StatusAccepted, rec.Code)
		assert.Equal(t, res+"\n", rec.Body.String())
	}
}

func TestAddBalanceUnSuccessfulUserNotFound(t *testing.T) {
	s := echo.New()

	res := `{"success":false,"message":"user not found","status_code":404}`

	body := `{"amount": 10, "transaction_id": "tx_1as4ndakda"}`

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	c := s.NewContext(req, rec)
	c.SetPath("/users/:uid/add")

	// params
	c.SetParamNames("uid")
	c.SetParamValues("6d7750a1-c3f2-4765-bf8f-33bc80f3f809")

	h, mock, err := newTest(s)
	if err != nil {
		panic(err)
	}

	// uRows := sqlmock.NewRows([]string{"id", "name", "cash_balance"}).AddRow(0, "Test", 10.10)

	mock.ExpectBegin()
	query := `SELECT "u".* FROM "users" AS "u" WHERE ("id" = '6d7750a1-c3f2-4765-bf8f-33bc80f3f809') FOR UPDATE`
	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnError(response.WrapError(fmt.Errorf("user not found"), http.StatusNotFound, ""))

	mock.ExpectRollback()

	if assert.NoError(t, h.addBalance(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, res+"\n", rec.Body.String())
	}
}

func TestAddBalanceUnSuccessfulTransactionIDExists(t *testing.T) {
	s := echo.New()

	res := `{"success":false,"message":"transaction was already processed","status_code":422}`
	body := `{
    					"amount": 10,
						"transaction_id": "tx_1as4ndakda"
					}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	c := s.NewContext(req, rec)
	c.SetPath("/users/:uid/add")
	// params
	c.SetParamNames("uid")
	c.SetParamValues("6d7750a1-c3f2-4765-bf8f-33bc80f3f809")

	db, mock := utils.MockSqlxDB()
	ur := repository.NewUserRepo(db)

	us := usecase.NewUserUseCase(ur)

	uRows := sqlmock.NewRows([]string{"id", "name", "balance"}).AddRow("6d7750a1-c3f2-4765-bf8f-33bc80f3f809", "Test", 100.10)

	mock.ExpectBegin()
	query := `SELECT "u".* FROM "users" AS "u" WHERE ("id" = '6d7750a1-c3f2-4765-bf8f-33bc80f3f809')`
	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(uRows)

	query = `SELECT "t".* FROM "transactions" AS "t" WHERE ("transaction_id" = 'tx_1as4ndakda')`
	uRows = sqlmock.NewRows([]string{"id", "created_at", "amount", "transaction_id"}).AddRow(1, time.Now().UTC(), 100.10, "tx_1as4ndakda")
	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(uRows)
	//mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnError(sql.ErrNoRows)

	query = `UPDATE "users" SET "balance"=balance + 10 WHERE ("id" = '6d7750a1-c3f2-4765-bf8f-33bc80f3f809')`
	mock.ExpectExec(regexp.QuoteMeta(query)).WillReturnResult(sqlmock.NewResult(0, 1))

	query = `INSERT INTO "transactions" ("amount", "created_at", "transaction_id", "user_id")`
	mock.ExpectExec(regexp.QuoteMeta(query)).WillReturnResult(sqlmock.NewErrorResult(nil))

	query = `SELECT "u".* FROM "users" AS "u" WHERE ("id" = '6d7750a1-c3f2-4765-bf8f-33bc80f3f809')`
	uRows = sqlmock.NewRows([]string{"id", "name", "balance"}).AddRow("6d7750a1-c3f2-4765-bf8f-33bc80f3f809", "Test", 110.10)
	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(uRows)

	mock.ExpectCommit()

	h := NewHandler(s, us)

	if assert.NoError(t, h.addBalance(c)) {
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		assert.Equal(t, res+"\n", rec.Body.String())
	}
}

func newTest(e *echo.Echo) (Handler, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		return Handler{}, nil, err
	}

	dbp := sqlx.NewDb(db, "postgres")

	// repos
	ur := repository.NewUserRepo(dbp)

	// use cases
	us := usecase.NewUserUseCase(ur)

	h := NewHandler(e, us)
	return h, mock, nil
}
