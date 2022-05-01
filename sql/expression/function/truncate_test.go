// Copyright 2020-2021 Dolthub, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either decress or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package function

import (
	"testing"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/expression"
	"github.com/stretchr/testify/require"
)

func TestTruncate(t *testing.T) {
	testCases := []struct {
		name     string
		rowType  sql.Type
		row      sql.Row
		decected interface{}
		err      bool
	}{
		{"Number and dec are nil", sql.Float64, sql.NewRow(nil, nil), nil, false},
		{"Number is nil", sql.Float64, sql.NewRow(2, nil), nil, false},
		{"Dec is nil", sql.Float64, sql.NewRow(nil, 2), nil, false},

		{"Number is 0", sql.Float64, sql.NewRow(0, 2), float64(0), false},
		{"Number and dec is 0", sql.Float64, sql.NewRow(0, 0), float64(0), false},
		{"Dec is 0", sql.Float64, sql.NewRow(2.123, 0), float64(2), false},
		{"Number is negative", sql.Float64, sql.NewRow(-298.0123, 2), float64(-298.01), false},
		{"Dec is negative", sql.Float64, sql.NewRow(2231.1, -2), float64(2200), false},
		{"Number and dec are invalid strings", sql.Float64, sql.NewRow("a", "b"), nil, true},
		{"Number and dec are valid strings", sql.Float64, sql.NewRow("232.333", "2"), float64(232.33), false},
	}
	for _, tt := range testCases {
		f := NewTruncate(
			expression.NewGetField(0, tt.rowType, "", false),
			expression.NewGetField(1, tt.rowType, "", false),
		)
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()
			require := require.New(t)
			ctx := sql.NewEmptyContext()

			v, err := f.Eval(ctx, tt.row)
			if tt.err {
				require.Error(err)
			} else {
				require.NoError(err)
				require.Equal(tt.decected, v)
			}
		})
	}
}
