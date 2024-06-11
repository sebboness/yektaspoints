package handlers

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_LambdaAuthContext_GetAuthorizerInfo(t *testing.T) {
	type state struct {
		setupCtxWithInfo bool
	}
	type want struct {
		hasInfo bool
	}
	type test struct {
		name string
		state
		want
	}

	cases := []test{
		{"has info", state{setupCtxWithInfo: true}, want{hasInfo: true}},
		{"no info", state{setupCtxWithInfo: false}, want{hasInfo: false}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			ctx := context.Background()

			if c.state.setupCtxWithInfo {
				ctx = context.WithValue(ctx, CtxKeyAuthInfo, AuthorizerInfo{
					Claims: map[string]any{
						"sub": "123",
					},
				})
			}

			w := httptest.NewRecorder()
			cgin, _ := gin.CreateTestContext(w)
			cgin.Request = httptest.NewRequest("GET", "/", nil).WithContext(ctx)

			authContext, err := NewLambdaAuthContext()
			assert.Nil(t, err)

			res := authContext.GetAuthorizerInfo(cgin)
			assert.Equal(t, c.want.hasInfo, res.HasInfo())
		})
	}
}
