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

package response

import (
	b64 "encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

// Cursor is a structure maintaing information about the cursor
type Cursor struct {
	ID uint
}

// ToBase64String converts a cursor to a base64 encoded string
func (c *Cursor) ToBase64String() string {
	cursor := fmt.Sprintf("%d", c.ID)
	return b64.StdEncoding.EncodeToString([]byte(cursor))
}

// ParseCursor returns a cursor structure from a base64 encoded string
func ParseCursor(cursor string) (*Cursor, error) {
	fromID, err := b64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, fmt.Errorf("error decoding pagination cursor")
	}

	fields := strings.Split(string(fromID), "/")
	if len(fields) < 1 {
		return nil, fmt.Errorf("invalid pagination cursor")
	}

	cID, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("invalid ID")
	}

	return &Cursor{ID: uint(cID)}, nil
}
