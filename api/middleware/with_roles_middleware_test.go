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

func Test_WithRolesAny(t *testing.T) {
	type state struct {
		mwRoles   []string
		authRoles string
		noAuth    bool
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
		{"happy path - 1mw role - multiple auth roles", state{mwRoles: []string{"parent"}, authRoles: "admin,parent"}, want{code: 200}},
		{"happy path - 1mw role - just parent role", state{mwRoles: []string{"parent"}, authRoles: "parent"}, want{code: 200}},
		{"happy path - 2mw roles - multiple auth roles", state{mwRoles: []string{"admin", "parent"}, authRoles: "admin,parent"}, want{code: 200}},
		{"happy path - 2mw roles - just parent role", state{mwRoles: []string{"admin", "parent"}, authRoles: "parent"}, want{code: 200}},
		{"fail - no matching role", state{mwRoles: []string{"parent"}, authRoles: "child"}, want{code: 401, err: "user not a parent"}},
		{"fail - no roles", state{mwRoles: []string{"parent"}, authRoles: ""}, want{code: 401, err: "user not a parent"}},
		{"fail - no auth", state{mwRoles: []string{"parent"}, noAuth: true}, want{code: 401, err: "user not a parent"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := handlers.PrepareAuthorizedContext(context.Background(), handlers.GetMockApiGWEvent(true, map[string]interface{}{
				"cognito:groups": c.state.authRoles,
			}))

			if c.state.noAuth {
				ctx = context.Background()
			}

			cgin := gin.Default()
			cgin.GET("/test", WithRolesAny(c.state.mwRoles), func(cgin *gin.Context) {
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
