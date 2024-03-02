package userauth

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gin-gonic/gin"
	"github.com/sebboness/yektaspoints/handlers"
	mocks "github.com/sebboness/yektaspoints/mocks/auth"
	"github.com/sebboness/yektaspoints/util/tests"
	"github.com/stretchr/testify/assert"
)

func Test_UserRegisterHandler(t *testing.T) {
	type state struct {
		hasNoAuth bool
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
		{"happy path", state{}, want{"", 200}},
		{"fail - unauthorized", state{hasNoAuth: true}, want{"unauthorized", 401}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			mockAuther := mocks.NewMockAuthController(t)

			ctrl := UserAuthController{
				auth: mockAuther,
			}

			ctx := context.Background()

			evt := events.APIGatewayProxyRequest{
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: map[string]interface{}{
						"claims": map[string]string{
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

			ctx = handlers.PrepareAuthorizedContext(ctx, evt)

			w := httptest.NewRecorder()
			cgin, _ := gin.CreateTestContext(w)
			cgin.Request = httptest.NewRequest("GET", "/", nil).WithContext(ctx)

			ctrl.GetUserAuthHandler(cgin)

			assert.Equal(t, c.want.code, w.Code)
			result := tests.AssertResult(t, w.Body)
			tests.AssertResultError(t, result, c.want.err)

			if c.want.code == 200 {
				assert.Equal(t, "john", result.Data.(map[string]interface{})["username"])
			}

			mockAuther.AssertExpectations(t)
		})
	}
}
