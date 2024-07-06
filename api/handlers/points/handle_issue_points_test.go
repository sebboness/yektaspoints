package points

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sebboness/yektaspoints/handlers"
	handlerMocks "github.com/sebboness/yektaspoints/mocks/handlers"
	mocks "github.com/sebboness/yektaspoints/mocks/storage"
	apierr "github.com/sebboness/yektaspoints/util/error"
	"github.com/sebboness/yektaspoints/util/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Controller_IssuePointsHandler(t *testing.T) {
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

			req := &issuePointsHandlerRequest{
				Points: 1,
				Reason: "I worked hard",
			}

			evtBody, _ := json.Marshal(req)
			evtBodyStr := string(evtBody)

			mockAuthContext := handlerMocks.NewMockAuthContext(t)
			mockPointsDB := mocks.NewMockIPointsStorage(t)

			authInfo := handlers.AuthorizerInfo{
				Claims: handlers.DefaultMockAuthClaims,
			}

			if !c.state.invalidBody {
				mockAuthContext.EXPECT().GetAuthorizerInfo(mock.Anything).Return(authInfo)
				mockPointsDB.EXPECT().SavePoint(mock.Anything, mock.Anything).Return(c.state.errSavePoint).Once()
			} else {
				evtBodyStr = `{"user_id":`
			}

			ctrl := PointsController{
				BaseController: handlers.BaseController{
					AuthContext: mockAuthContext,
				},
				pointsDB: mockPointsDB,
			}

			w := httptest.NewRecorder()
			cgin, _ := gin.CreateTestContext(w)
			cgin.Request = httptest.NewRequest("POST", "/points", bytes.NewReader([]byte(evtBodyStr)))

			ctrl.IssuePointsHandler(cgin)

			assert.Equal(t, c.want.code, w.Code)
			result := tests.AssertResult(t, w.Body)
			tests.AssertResultError(t, result, c.want.err)

			if c.want.code == 200 {
				assert.NotNil(t, result.Data)
			}

			mockAuthContext.AssertExpectations(t)
			mockPointsDB.AssertExpectations(t)
		})
	}
}

func Test_Controller_handleIssuePoints(t *testing.T) {
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
		t.Run(c.name, func(t *testing.T) {
			req := &issuePointsHandlerRequest{
				ToChildUserID: "123",
				Points:        1,
				Reason:        "I worked hard",
			}

			if c.state.validationError {
				req.Points = -1
			}

			mockPointsDB := mocks.NewMockIPointsStorage(t)

			saveCalled := map[bool]int{true: 0, false: 1}[c.state.validationError]
			if saveCalled > 0 {
				mockPointsDB.EXPECT().SavePoint(mock.Anything, mock.Anything).Return(c.state.errSavePoint).Times(saveCalled)
			}

			ctrl := PointsController{
				pointsDB: mockPointsDB,
			}

			ctx := context.Background()
			_, err := ctrl.handleIssuePoints(ctx, req)
			if err != nil {
				assert.Contains(t, err.Error(), c.want.err)
			}

			mockPointsDB.AssertExpectations(t)
		})
	}
}

func Test_validateIssuePoints(t *testing.T) {
	type state struct {
		invalidFromParentUserId bool
		invalidToChildUserId    bool
		pointsAreNegative       bool
		pointsAreZero           bool
		missingReason           bool
		tooShortReason          bool
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
		{"fail - invalid from parent user id", state{invalidFromParentUserId: true}, want{"missing from_parent_user_id"}},
		{"fail - invalid to child user id", state{invalidToChildUserId: true}, want{"missing to_child_user_id"}},
		{"fail - invalid points - negative", state{pointsAreNegative: true}, want{"failed to validate request: points must be a positive integer"}},
		{"fail - invalid points - zero", state{pointsAreZero: true}, want{"failed to validate request: points must be a positive integer"}},
		{"fail - missing reason", state{missingReason: true}, want{"failed to validate request: reason for requesting points must not be empty"}},
		{"fail - reason too short", state{tooShortReason: true}, want{"failed to validate request: reason for requesting points must not be empty"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := &issuePointsHandlerRequest{
				ToChildUserID: "123",
				Points:        1,
				Reason:        "I worked hard",
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
			if c.state.invalidFromParentUserId {
				req.FromParentUserID = ""
			}

			err := validateIssuePoints(req)
			if err != nil {
				assert.Contains(t, err.Error(), c.want.err)
			}
		})
	}
}
