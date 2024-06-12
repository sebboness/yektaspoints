package handlers

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	mocks "github.com/sebboness/yektaspoints/mocks/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_ApiAuthContext_GetAuthorizerInfo(t *testing.T) {
	type state struct {
		authInContext     bool
		missingAuthHeader bool
		getUserErr        error
	}
	type want struct {
		claims map[string]interface{}
	}
	type test struct {
		name string
		state
		want
	}

	cases := []test{
		{"happy path", state{}, want{map[string]interface{}{"sub": "123"}}},
		{"happy path - already in context", state{authInContext: true}, want{map[string]interface{}{"sub": "555"}}},
		{"empty - no auth header", state{missingAuthHeader: true}, want{}},
		{"empty - get claims err", state{getUserErr: errFail}, want{}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			ctx := context.Background()
			mockUserDB := mocks.NewMockJwtParser(t)

			claims := map[string]interface{}{
				"sub": "123",
			}

			if c.state.authInContext {
				ctx = context.WithValue(ctx, CtxKeyAuthInfo, AuthorizerInfo{
					Claims: map[string]any{
						"sub": "555",
					},
				})
			}

			if !c.state.missingAuthHeader && !c.state.authInContext {
				mockUserDB.EXPECT().GetJwtClaims(mock.Anything).Return(claims, c.state.getUserErr).Once()
			}

			w := httptest.NewRecorder()
			cgin, _ := gin.CreateTestContext(w)
			cgin.Request = httptest.NewRequest("GET", "/", nil).WithContext(ctx)

			if !c.state.missingAuthHeader {
				cgin.Request.Header.Set("Authorization", "Bearer 123")
			}

			authContext := &ApiAuthContext{
				JwtParser: mockUserDB,
			}

			res := authContext.GetAuthorizerInfo(cgin)

			assert.Equal(t, c.want.claims, res.Claims)

			mockUserDB.AssertExpectations(t)
		})
	}
}
