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
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestParseCursor(t *testing.T) {
	type args struct {
		cursor string
	}
	tests := []struct {
		name    string
		args    args
		want    *Cursor
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "t-01",
			args:    args{cursor: "MQ=="},
			want:    &Cursor{ID: 1},
			wantErr: false,
		},
		{
			name:    "t-02",
			args:    args{cursor: "M"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCursor(tt.args.cursor)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCursor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseCursor() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCursorToBase64String(t *testing.T) {
	t.Run("good", func(t *testing.T) {
		c := Cursor{ID: 1}
		bStr := c.ToBase64String()

		assert.Equal(t, "MQ==", bStr)
	})
}
