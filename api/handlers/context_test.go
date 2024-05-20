package handlers

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_PrepareAuthorizedContext(t *testing.T) {
	type state struct {
		hasNoClaims bool
		hasNoAuth   bool
	}
	type want struct {
	}
	type test struct {
		name string
		state
		want
	}

	cases := []test{
		{"happy path", state{}, want{}},
		{"without claims", state{hasNoClaims: true}, want{}},
		{"without authorizer", state{hasNoAuth: true}, want{}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			ctx := context.Background()

			evt := events.APIGatewayProxyRequest{
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: map[string]interface{}{
						"claims": map[string]interface{}{
							"cognito:username": "john",
							"email":            "john@info.co",
							"email_verified":   "true",
							"name":             "John",
							"sub":              "123",
						},
					},
				},
			}

			if c.state.hasNoAuth {
				evt.RequestContext.Authorizer = nil
			}
			if c.state.hasNoClaims {
				evt.RequestContext.Authorizer = map[string]interface{}{"bleep": "bloop"}
			}

			ctx = PrepareAuthorizedContext(ctx, evt)
			assert.NotNil(t, ctx.Value(CtxKeyAuthInfo))
			assert.IsType(t, AuthorizerInfo{}, ctx.Value(CtxKeyAuthInfo))
		})
	}
}

func Test_GetAuthorizerInfo(t *testing.T) {
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

			res := GetAuthorizerInfo(cgin)
			assert.Equal(t, c.want.hasInfo, res.HasInfo())
		})
	}
}
