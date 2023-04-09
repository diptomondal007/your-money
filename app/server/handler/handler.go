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
	"github.com/labstack/echo/v4"

	"github.com/diptomondal007/your-money/app/server/usecase"
)

type Handler struct {
	e  *echo.Echo
	uc usecase.UserUseCase
}

func NewHandler(e *echo.Echo, uc usecase.UserUseCase) Handler {
	h := Handler{e: e, uc: uc}

	// user group
	ug := e.Group("/users/:uid")

	ug.POST("/add", h.addBalance)
	ug.GET("/balance", h.checkBalance)
	ug.GET("/history", h.history)

	return h
}
