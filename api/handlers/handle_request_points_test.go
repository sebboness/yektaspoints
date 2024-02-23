package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	mocks "github.com/sebboness/yektaspoints/mocks/storage"
	apierr "github.com/sebboness/yektaspoints/util/error"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var errFail = errors.New("fail")

func Test_Controller_RequestPointsHandler(t *testing.T) {
	type state struct {
		invalidBody  bool
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
		{"fail - invalid body", state{invalidBody: true}, want{"failed to unmarshal json body", 400}},
		{"fail - validation error", state{errSavePoint: apierr.New(apierr.InvalidInput)}, want{"invalid input", 400}},
		{"fail - unauthorized", state{errSavePoint: apierr.New(apierr.Unauthorized)}, want{"unauthorized", 401}},
		{"fail - internal server error", state{errSavePoint: errors.New("fail")}, want{"fail", 500}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			req := &pointsHandlerRequest{
				Points: 1,
				Reason: "I worked hard",
			}

			evtBody, _ := json.Marshal(req)
			evtBodyStr := string(evtBody)

			mockPointsDB := mocks.NewMockIPointsStorage(t)

			if !c.state.invalidBody {
				mockPointsDB.EXPECT().SavePoint(mock.Anything, mock.Anything).Return(c.state.errSavePoint).Once()
			} else {
				evtBodyStr = `{"user_id":`
			}

			ctrl := LambdaController{
				pointsDB: mockPointsDB,
			}

			ctx := PrepareAuthorizedContext(context.Background(), mockApiGWEvent)

			w := httptest.NewRecorder()
			cgin, _ := gin.CreateTestContext(w)
			cgin.Request = httptest.NewRequest("POST", "/points", bytes.NewReader([]byte(evtBodyStr))).WithContext(ctx)
			PrepareAuthorizedContext(ctx, mockApiGWEvent)

			ctrl.RequestPointsHandler(cgin)

			assert.Equal(t, c.want.code, w.Code)
			body, err := io.ReadAll(w.Body)
			assert.Nil(t, err, "reading response body should have no error")

			if c.want.err == "" {
				assert.Contains(t, string(body), `"points":0`)
			} else {
				assert.Contains(t, string(body), c.want.err)
			}

			mockPointsDB.AssertExpectations(t)
		})
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
		{"fail - save points", state{errSavePoint: errFail}, want{"failed to save points: fail"}},
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

		ctrl := LambdaController{
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
