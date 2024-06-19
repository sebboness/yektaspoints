package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sebboness/yektaspoints/handlers"
	"github.com/sebboness/yektaspoints/util/tests"
	"github.com/stretchr/testify/assert"
)

func Test_WithAuthorizedUser(t *testing.T) {
	type state struct {
		noAuth           bool
		noUserID         bool
		emailNotVerified bool
	}
	type want struct {
		err  string
		code int
	}
	type test struct {
		name string
		state
		want
	}

	cases := []test{
		{"happy path", state{}, want{code: 200}},
		{"fail - no user id", state{noUserID: true}, want{code: 401, err: "unknown user ID"}},
		{"fail - email not verified", state{emailNotVerified: true}, want{code: 401, err: "unverified user"}},
		{"fail - no auth", state{noAuth: true}, want{code: 401, err: "unauthorized"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			claims := map[string]interface{}{
				"cognito:username": "john",
				"email":            "john@info.co",
				"email_verified":   "true",
				"name":             "John",
				"sub":              "123",
			}

			if c.state.noUserID {
				delete(claims, "sub")
			}
			if c.state.emailNotVerified {
				delete(claims, "email_verified")
			}

			ctx := handlers.PrepareAuthorizedContext(context.Background(), handlers.GetMockApiGWEvent(false, claims))

			if c.state.noAuth {
				ctx = context.Background()
			}

			cgin := gin.Default()
			cgin.GET("/test", WithAuthorizedUser(), func(cgin *gin.Context) {
				cgin.JSON(http.StatusOK, handlers.SuccessResult(1))
			})

			req, _ := http.NewRequest("GET", "/test", nil)
			req = req.WithContext(ctx)

			testHTTPResponse(cgin, req, func(w *httptest.ResponseRecorder) {
				assert.Equal(t, c.want.code, w.Code)
				result := tests.AssertResult(t, w.Body)
				tests.AssertResultError(t, result, c.want.err)
			})
		})
	}
}
