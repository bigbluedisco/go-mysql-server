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
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sql

import (
	"encoding/json"

	"github.com/dolthub/vitess/go/sqltypes"
	"github.com/dolthub/vitess/go/vt/proto/query"
	"gopkg.in/src-d/go-errors.v1"
)

var ErrConvertingToJSON = errors.NewKind("value %v is not valid JSON")

var JSON JsonType = jsonType{}

type JsonType interface {
	Type
}

type jsonType struct{}

// Compare implements Type interface.
func (t jsonType) Compare(a interface{}, b interface{}) (int, error) {
	var err error
	if a, err = t.Convert(a); err != nil {
		return 0, err
	}
	if b, err = t.Convert(b); err != nil {
		return 0, err
	}
	// todo: making a context here is expensive
	return a.(JSONValue).Compare(NewEmptyContext(), b.(JSONValue))
}

// Convert implements Type interface.
func (t jsonType) Convert(v interface{}) (doc interface{}, err error) {
	switch v := v.(type) {
	case JSONValue:
		return v, nil
	case []byte:
		err = json.Unmarshal(v, &doc)
	case string:
		err = json.Unmarshal([]byte(v), &doc)
	case []string, []int, []uint16, []uint32, []uint64, []int16, []int32, []int64, []float32, []float64, []bool,
		map[string]string, map[string]int, map[string]uint16, map[string]uint32, map[string]uint64, map[string]int16, map[string]int32, map[string]int64, map[string]float32, map[string]float64:
		doc = v
	default:
		// if |v| can be marshalled, it contains
		// a valid JSON document representation
		if b, berr := json.Marshal(v); berr == nil {
			err = json.Unmarshal(b, &doc)
		}
	}
	if err != nil {
		return nil, err
	}
	return JSONDocument{Val: doc}, nil
}

// Promote implements the Type interface.
func (t jsonType) Promote() Type {
	return t
}

// SQL implements Type interface.
func (t jsonType) SQL(v interface{}) (sqltypes.Value, error) {
	if v == nil {
		return sqltypes.NULL, nil
	}

	js, ok := v.(JSONValue)
	if !ok {
		return sqltypes.NULL, nil
	}

	// todo: making a context here is expensive
	s, err := js.ToString(NewEmptyContext())
	if err != nil {
		return sqltypes.NULL, err
	}

	return sqltypes.MakeTrusted(sqltypes.TypeJSON, []byte(s)), nil
}

// String implements Type interface.
func (t jsonType) String() string {
	return "JSON"
}

// Type implements Type interface.
func (t jsonType) Type() query.Type {
	return sqltypes.TypeJSON
}

// Zero implements Type interface.
func (t jsonType) Zero() interface{} {
	// JSON Null
	return nil
}
