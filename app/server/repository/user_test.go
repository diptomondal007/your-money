// Licensed to Dipto Mondal under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Dipto Mondal licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
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
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/diptomondal007/your-money/app/utils"
)

func TestCheckBalance(t *testing.T) {
	db, mock := utils.MockSqlxDB()
	defer db.Close()

	ur := NewUserRepo(db)

	id := "6d7750a1-c3f2-4765-bf8f-33bc80f3f809"

	query := `SELECT "u".* FROM "users" AS "u" WHERE ("id" = '6d7750a1-c3f2-4765-bf8f-33bc80f3f809')`

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "balance"}).
		AddRow("6d7750a1-c3f2-4765-bf8f-33bc80f3f809", time.Now().UTC(), time.Now().UTC(), "Test", 10).
		AddRow("6d7750a1-c3f2-4765-bf8f-33bc80f3f80a", time.Now().UTC(), time.Now().UTC(), "Other", 10)

	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)

	user, err := ur.GetUserInfo(id)

	assert.NoError(t, err)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		assert.NoError(t, err)
	}

	assert.Equal(t, id, user.ID)
	assert.Equal(t, "Test", user.Name)
}

func TestGetHistoryCount(t *testing.T) {
	db, mock := utils.MockSqlxDB()
	defer db.Close()

	ur := NewUserRepo(db)

	id := "6d7750a1-c3f2-4765-bf8f-33bc80f3f809"

	query := `SELECT COUNT(*) FROM "transactions" AS "t" WHERE ("user_id" = '6d7750a1-c3f2-4765-bf8f-33bc80f3f809')`

	rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)

	count, err := ur.GetHistoryCount(id)

	assert.NoError(t, err)
	assert.Equal(t, count, int64(2))
}

func TestGetHistoryList(t *testing.T) {
	db, mock := utils.MockSqlxDB()
	defer db.Close()

	ur := NewUserRepo(db)

	id := "6d7750a1-c3f2-4765-bf8f-33bc80f3f809"
	pageSize := 10
	cursor := ""

	query := `SELECT "t".* FROM "transactions" AS "t" WHERE ("user_id" = '6d7750a1-c3f2-4765-bf8f-33bc80f3f809') ORDER BY "t"."id" DESC LIMIT 10`

	rows := sqlmock.NewRows([]string{"id", "created_at", "amount", "transaction_id", "user_id"}).
		AddRow(1, time.Now().UTC(), 10, "tx_asdasfa", "6d7750a1-c3f2-4765-bf8f-33bc80f3f809")

	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)

	list, err := ur.GetHistoryList(id, int64(pageSize), cursor)

	assert.NoError(t, err)
	assert.Equal(t, len(list), 1)
}

func TestAddBalance(t *testing.T) {
	db, mock := utils.MockSqlxDB()
	defer db.Close()

	ur := NewUserRepo(db)

	id := "6d7750a1-c3f2-4765-bf8f-33bc80f3f809"
	//transactionID := "tx_1as4ndakda"
	//amount := 10

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

	user, err := ur.AddBalance(id, "tx_1as4ndakda", 10)

	assert.NoError(t, err)
	assert.Equal(t, user.Balance, 110.1)
}
