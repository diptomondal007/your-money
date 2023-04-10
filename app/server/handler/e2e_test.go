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

package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"

	"github.com/diptomondal007/your-money/app/server/repository"
	"github.com/diptomondal007/your-money/app/server/usecase"
	"github.com/diptomondal007/your-money/infrastructure/conn"
)

type e2eTestSuite struct {
	suite.Suite
	repo repository.UserRepository
}

func (s *e2eTestSuite) SetupSuite() {
	if testing.Short() {
		return
	}

	err := conn.ConnectDB()
	s.NoError(err)

	ur := repository.NewUserRepo(conn.GetDB().DB)
	s.repo = ur
}

func TestE2ETestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, &e2eTestSuite{})
}

func (s *e2eTestSuite) TestE2EGetBalance() {
	user, err := s.repo.CheckBalance("6d7750a1-c3f2-4765-bf8f-33bc80f3f809")
	s.NoError(err)

	response, err := s.req(echo.GET, fmt.Sprintf("http://localhost:%d/users/6d7750a1-c3f2-4765-bf8f-33bc80f3f809/balance", 8080), nil)
	s.NoError(err)
	s.Equal(http.StatusOK, response.StatusCode)

	byteBody, err := io.ReadAll(response.Body)
	s.NoError(err)

	s.Equal(fmt.Sprintf(`{"success":true,"message":"request successful!","status_code":200,"data":{"balance":%d}}`, int(user.Balance)), strings.Trim(string(byteBody), "\n"))

	response.Body.Close()
}

func (s *e2eTestSuite) TestE2EAddBalance() {
	user, err := s.repo.CheckBalance("6d7750a1-c3f2-4765-bf8f-33bc80f3f809")
	s.NoError(err)

	amount := 10
	body := usecase.AddBalanceReq{
		TransactionID: "tx_1as4ndakda",
		Amount:        float64(amount),
	}

	response, err := s.req(echo.POST, fmt.Sprintf("http://localhost:%d/users/6d7750a1-c3f2-4765-bf8f-33bc80f3f809/add", 8080), body)

	s.NoError(err)
	s.Equal(http.StatusAccepted, response.StatusCode)

	byteBody, err := io.ReadAll(response.Body)
	s.NoError(err)

	s.Equal(fmt.Sprintf(`{"success":true,"message":"transaction successful!","status_code":202,"data":{"current_balance":%d}}`, int(user.Balance+float64(amount))), strings.Trim(string(byteBody), "\n"))

	response.Body.Close()
}

func (s *e2eTestSuite) TestE2EAddBalanceConcurrent() {
	user, err := s.repo.CheckBalance("6d7750a1-c3f2-4765-bf8f-33bc80f3f809")
	s.NoError(err)

	ids := []string{
		"6022714107",
		"4177041956",
		"2555431161",
		"4032143776",
		"2633417522",
	}

	wg := new(sync.WaitGroup)
	for i := 0; i < 5; i++ {
		wg.Add(1)
		tx := ids[i]
		go func(wg *sync.WaitGroup, id string) {
			amount := 10
			body := usecase.AddBalanceReq{
				TransactionID: fmt.Sprintf("tx_%s", id),
				Amount:        float64(amount),
			}

			_, _ = s.req(echo.POST, fmt.Sprintf("http://localhost:%d/users/6d7750a1-c3f2-4765-bf8f-33bc80f3f809/add", 8080), body)

			wg.Done()
		}(wg, tx)

	}

	wg.Wait()

	u, err := s.repo.CheckBalance("6d7750a1-c3f2-4765-bf8f-33bc80f3f809")
	s.Equal(user.Balance+(5*10), u.Balance)
}

func (s *e2eTestSuite) TestE2ETransactionList() {
	response, err := s.req(echo.GET, fmt.Sprintf("http://localhost:%d/users/6d7750a1-c3f2-4765-bf8f-33bc80f3f809/history?page_size=10", 8080), nil)

	s.NoError(err)
	s.Equal(http.StatusOK, response.StatusCode)

	response.Body.Close()
}

func (s *e2eTestSuite) req(method, url string, body interface{}) (*http.Response, error) {
	var buf bytes.Buffer

	if body != nil {
		err := json.NewEncoder(&buf).Encode(body)
		s.NoError(err)
	}

	req, err := http.NewRequest(method, url, &buf)
	s.NoError(err)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	client := http.Client{}
	response, err := client.Do(req)
	return response, err
}
