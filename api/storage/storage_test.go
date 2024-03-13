package storage

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/sebboness/yektaspoints/models"
	"github.com/sebboness/yektaspoints/util/tests"
	"github.com/stretchr/testify/assert"
)

func Test_dateFilterExpression(t *testing.T) {
	type state struct {
		from *time.Time
		to   *time.Time
	}
	type want struct {
	}
	type test struct {
		name string
		state
		want
	}

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	cases := []test{
		{"between", state{from: &from, to: &to}, want{}},
		{"from", state{from: &from}, want{}},
		{"to", state{to: &to}, want{}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			df := models.DateFilter{
				From: c.state.from,
				To:   c.state.to,
			}

			_ = dateFilterExpression("date", df)
		})
	}
}

func Test_valueInListExpression(t *testing.T) {
	type state struct {
		slice []any
	}
	type want struct {
	}
	type test struct {
		name string
		state
		want
	}

	cases := []test{
		{"none", state{[]any{}}, want{}},
		{"integers one", state{[]any{1}}, want{}},
		{"integers many", state{[]any{1, 2, 3, 4}}, want{}},
		{"strings", state{[]any{"1", "2", "3", "4", "abc"}}, want{}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_ = valueInListExpression("type_id", c.state.slice)
		})
	}
}

func Test_selectAttributesExpression(t *testing.T) {
	type state struct {
		names []string
	}
	type want struct {
		err string
		len int
	}
	type test struct {
		name string
		state
		want
	}

	cases := []test{
		{"none", state{[]string{}}, want{"unset parameter: ProjectionBuilder", 0}},
		{"one", state{[]string{"a"}}, want{"", 1}},
		{"many", state{[]string{"a", "b", "c", "d"}}, want{"", 4}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			projEx := selectAttributesExpression(c.state.names)
			expr, err := expression.NewBuilder().WithProjection(projEx).Build()

			tests.AssertError(t, err, c.want.err)
			assert.NotNil(t, expr)
			assert.Len(t, expr.Names(), c.want.len)
		})
	}
}
