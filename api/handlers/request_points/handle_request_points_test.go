package request_points

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	mocks "github.com/sebboness/yektaspoints/mocks/storage"
	apierr "github.com/sebboness/yektaspoints/util/error"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var fail = errors.New("fail")

func Test_Controller_RequestPointsHandler(t *testing.T) {
	type state struct {
		errSavePoint error
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
		{"fail - validation error", state{errSavePoint: apierr.New(apierr.InvalidInput)}, want{"invalid input", 400}},
		{"fail - unauthorized", state{errSavePoint: apierr.New(apierr.Unauthorized)}, want{"unauthorized", 401}},
		{"fail - internal server error", state{errSavePoint: errors.New("fail")}, want{"fail", 500}},
	}

	for _, c := range cases {

		req := &pointsHandlerRequest{
			APIGatewayProxyRequest: events.APIGatewayProxyRequest{
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: map[string]interface{}{
						"claims": map[string]interface{}{
							"cognito:username": "123",
						},
					},
				},
			},
			Points: 1,
			Reason: "I worked hard",
		}

		mockPointsDB := mocks.NewMockIPointsStorage(t)
		mockPointsDB.EXPECT().SavePoint(mock.Anything, mock.Anything).Return(c.state.errSavePoint).Once()

		ctrl := RequestPointsController{
			pointsDB: mockPointsDB,
		}

		ctx := context.Background()
		resp, err := ctrl.RequestPointsHandler(ctx, req)
		if err != nil {
			assert.Contains(t, err.Error(), c.want.err)
		}

		assert.Equal(t, c.want.code, resp.StatusCode)
	}
}

func Test_Controller_handleRequestPoints(t *testing.T) {
	type state struct {
		validationError bool
		errSavePoint    error
	}
	type want struct {
		err string
	}
	type test struct {
		name string
		state
		want
	}

	cases := []test{
		{"happy path", state{}, want{}},
		{"fail - validation error", state{validationError: true}, want{"invalid input: failed to validate request"}},
		{"fail - save points", state{errSavePoint: fail}, want{"failed to save points: fail"}},
	}

	for _, c := range cases {

		req := &pointsHandlerRequest{
			UserID: "123",
			Points: 1,
			Reason: "I worked hard",
		}

		if c.state.validationError {
			req.Points = -1
		}

		mockPointsDB := mocks.NewMockIPointsStorage(t)

		saveCalled := map[bool]int{true: 0, false: 1}[c.state.validationError]
		if saveCalled > 0 {
			mockPointsDB.EXPECT().SavePoint(mock.Anything, mock.Anything).Return(c.state.errSavePoint).Times(saveCalled)
		}

		ctrl := RequestPointsController{
			pointsDB: mockPointsDB,
		}

		ctx := context.Background()
		_, err := ctrl.handleRequestPoints(ctx, req)
		if err != nil {
			assert.Contains(t, err.Error(), c.want.err)
		}

		mockPointsDB.AssertExpectations(t)
	}
}

func Test_validateRequestPoints(t *testing.T) {
	type state struct {
		invalidUserId     bool
		pointsAreNegative bool
		pointsAreZero     bool
		missingReason     bool
		tooShortReason    bool
	}
	type want struct {
		err string
	}
	type test struct {
		name string
		state
		want
	}

	cases := []test{
		{"happy path", state{}, want{}},
		{"fail - invalid user id", state{invalidUserId: true}, want{"unauthorized: missing user ID"}},
		{"fail - invalid points - negative", state{pointsAreNegative: true}, want{"failed to validate request: points must be a positive integer"}},
		{"fail - invalid points - zero", state{pointsAreZero: true}, want{"failed to validate request: points must be a positive integer"}},
		{"fail - missing reason", state{missingReason: true}, want{"failed to validate request: reason for requesting points must not be empty"}},
		{"fail - reason too short", state{tooShortReason: true}, want{"failed to validate request: reason for requesting points must not be empty"}},
	}

	for _, c := range cases {

		req := &pointsHandlerRequest{
			UserID: "123",
			Points: 1,
			Reason: "I worked hard",
		}

		if c.state.pointsAreZero {
			req.Points = 0
		}
		if c.state.pointsAreNegative {
			req.Points = -1
		}
		if c.state.missingReason {
			req.Reason = ""
		}
		if c.state.tooShortReason {
			req.Reason = "hello"
		}
		if c.state.invalidUserId {
			req.UserID = ""
		}

		err := validateRequestPoints(req)
		if err != nil {
			assert.Contains(t, err.Error(), c.want.err)
		}
	}
}
