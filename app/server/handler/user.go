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
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/diptomondal007/your-money/app/server/usecase"
	"github.com/diptomondal007/your-money/app/utils/response"
)

func validateAddReq(req *usecase.AddBalanceReq) error {
	if req.Amount <= 0 {
		return fmt.Errorf("not a valid amount. amount should be a positive value")
	}

	if req.TransactionID == "" {
		return fmt.Errorf("valid transaction id required")
	}

	if !strings.HasPrefix(req.TransactionID, "tx_") {
		return fmt.Errorf("valid transaction id should have prefix tx_")
	}

	return nil
}

func (h *Handler) addBalance(c echo.Context) error {
	userID := c.Param("uid")
	if userID == "" {
		return c.JSON(response.RespondError(response.ErrBadRequest, fmt.Errorf("not a valid user id")))
	}

	var req *usecase.AddBalanceReq
	err := c.Bind(&req)
	if err != nil {
		log.Println("bad request body", err)
		return c.JSON(response.RespondError(response.ErrBadRequest, fmt.Errorf("not a valid request body")))
	}

	err = validateAddReq(req)
	if err != nil {
		log.Println("bad request data, req: ", req)
		return c.JSON(response.RespondError(response.ErrBadRequest, err))
	}

	u, err := h.uc.AddBalance(userID, req)
	if err != nil {
		return c.JSON(response.RespondError(err))
	}

	return c.JSON(response.RespondSuccess(http.StatusAccepted, "transaction successful!", u))
}

func (h *Handler) checkBalance(c echo.Context) error {
	userID := c.Param("uid")
	if userID == "" {
		return c.JSON(response.RespondError(response.ErrBadRequest, fmt.Errorf("not a valid user id")))
	}

	ds, err := h.uc.CheckBalance(userID)
	if err != nil {
		return c.JSON(response.RespondError(err))
	}

	return c.JSON(response.RespondSuccess(http.StatusOK, "request successful!", ds))
}

func (h *Handler) history(c echo.Context) error {
	userID := c.Param("uid")
	if userID == "" {
		return c.JSON(response.RespondError(response.ErrBadRequest, fmt.Errorf("not a valid user id")))
	}

	pageSize := c.QueryParam("page_size")

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		return c.JSON(response.RespondError(response.ErrBadRequest, fmt.Errorf("page size should be a valid integer")))
	}

	cursor := c.QueryParam("page")

	ds, err := h.uc.ListHistory(userID, int64(pageSizeInt), cursor)
	if err != nil {
		return c.JSON(response.RespondError(err))
	}

	return c.JSON(response.RespondSuccess(http.StatusOK, "request successful!", ds))
}
