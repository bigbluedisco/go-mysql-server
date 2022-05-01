package function

import (
	"fmt"
	"math"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/expression"
)

// Truncate is a function that returns a number truncated to certain decimal places.
type Truncate struct {
	expression.BinaryExpression
}

var _ sql.FunctionExpression = (*Truncate)(nil)

// NewTruncate creates a new Truncate expression.
func NewTruncate(e1, e2 sql.Expression) sql.Expression {
	return &Truncate{
		expression.BinaryExpression{
			Left:  e1,
			Right: e2,
		},
	}
}

// FunctionName implements sql.FunctionExpression
func (p *Truncate) FunctionName() string {
	return "truncate"
}

// Description implements sql.FunctionExpression
func (p *Truncate) Description() string {
	return "returns the value of N truncated to D decimal places."
}

// Type implements the Expression interface.
func (p *Truncate) Type() sql.Type { return sql.Float64 }

// IsNullable implements the Expression interface.
func (p *Truncate) IsNullable() bool { return p.Left.IsNullable() || p.Right.IsNullable() }

func (p *Truncate) String() string {
	return fmt.Sprintf("truncate(%s, %s)", p.Left, p.Right)
}

// WithChildren implements the Expression interface.
func (p *Truncate) WithChildren(children ...sql.Expression) (sql.Expression, error) {
	if len(children) != 2 {
		return nil, sql.ErrInvalidChildrenNumber.New(p, len(children), 2)
	}
	return NewTruncate(children[0], children[1]), nil
}

// Eval implements the Expression interface.
func (p *Truncate) Eval(ctx *sql.Context, row sql.Row) (interface{}, error) {
	left, err := p.Left.Eval(ctx, row)
	if err != nil {
		return nil, err
	}

	if left == nil {
		return nil, nil
	}

	left, err = sql.Float64.Convert(left)
	if err != nil {
		return nil, err
	}

	right, err := p.Right.Eval(ctx, row)
	if err != nil {
		return nil, err
	}

	if right == nil {
		return nil, nil
	}

	right, err = sql.Float64.Convert(right)
	if err != nil {
		return nil, err
	}

	d := math.Pow(10, -right.(float64))
	return math.Trunc(left.(float64)/d) * d, nil
}
