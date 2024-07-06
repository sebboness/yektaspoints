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
	"github.com/sebboness/yektaspoints/models"
	"github.com/sebboness/yektaspoints/util/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Controller_IssuePointsHandler(t *testing.T) {
	type state struct {
		validationError bool
		invalidBody     bool
		errSavePoint    error
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
		{"fail - validation error", state{validationError: true}, want{"invalid input", 400}},
		{"fail - internal server error", state{errSavePoint: errors.New("fail")}, want{"fail", 500}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			req := &issuePointsHandlerRequest{
				ToChildUserID: "child-1",
				ParentNotes:   "you did really well!",
				Points:        10,
				Reason:        "there is some reason",
				Type:          models.PointRequestTypeAdd,
			}

			if c.state.validationError {
				req.Reason = ""
			}

			evtBody, _ := json.Marshal(req)
			evtBodyStr := string(evtBody)

			mockAuthContext := handlerMocks.NewMockAuthContext(t)
			mockPointsDB := mocks.NewMockIPointsStorage(t)
			mockUserDB := mocks.NewMockIUserStorage(t)

			passedRequestChecks := !c.state.invalidBody
			passedInitialChecks := !c.state.invalidBody && !c.state.validationError

			if c.state.invalidBody {
				evtBodyStr = `{"user_id":`
			}

			authInfo := handlers.AuthorizerInfo{
				Claims: handlers.DefaultMockAuthClaims,
			}

			if passedRequestChecks {
				mockAuthContext.EXPECT().GetAuthorizerInfo(mock.Anything).Return(authInfo)
			}

			if passedInitialChecks {
				mockUserDB.EXPECT().ParentHasAccessToChild(mock.Anything, mock.Anything, mock.Anything).Return(true, nil).Once()
				mockPointsDB.EXPECT().GetLatestBalance(mock.Anything, mock.Anything).Return(models.PointBalance{}, nil).Once()
				mockPointsDB.EXPECT().SavePoint(mock.Anything, mock.Anything).Return(c.state.errSavePoint).Once()
			}

			ctrl := PointsController{
				BaseController: handlers.BaseController{
					AuthContext: mockAuthContext,
				},
				pointsDB: mockPointsDB,
				userDB:   mockUserDB,
			}

			ctx := context.Background()

			w := httptest.NewRecorder()

			cgin, _ := gin.CreateTestContext(w)

			cgin.Request = httptest.NewRequest("PUT", "/", bytes.NewReader([]byte(evtBodyStr))).WithContext(ctx)

			handlers.PrepareAuthorizedContext(ctx, handlers.MockApiGWEvent)

			ctrl.IssuePointsHandler(cgin)

			assert.Equal(t, c.want.code, w.Code)
			result := tests.AssertResult(t, w.Body)
			tests.AssertResultError(t, result, c.want.err)

			if c.want.code == 200 {
				assert.NotNil(t, result.Data)
			}

			mockAuthContext.AssertExpectations(t)
			mockPointsDB.AssertExpectations(t)
			mockUserDB.AssertExpectations(t)
		})
	}
}

func Test_Controller_handleIssuePoints(t *testing.T) {
	type state struct {
		pointsType      string
		validationError bool
		noAccess        bool
		errHasAccess    error
		errGetBalance   error
		errSavePoint    error
	}
	type want struct {
		bal int32
		err string
	}
	type test struct {
		name string
		state
		want
	}

	cases := []test{
		{"happy path - add", state{pointsType: "ADD"}, want{bal: 20}},
		{"happy path - subtract", state{pointsType: "SUBTRACT"}, want{bal: 0}},
		{"happy path - cashout", state{pointsType: "CASHOUT"}, want{bal: 0}},
		{"fail - validation error", state{pointsType: "ADD", validationError: true}, want{err: "invalid input: failed to validate request"}},
		{"fail - has access err", state{pointsType: "ADD", errHasAccess: errFail}, want{err: "failed to check user access permissions: fail"}},
		{"fail - no access", state{pointsType: "ADD", noAccess: true}, want{err: "requesting user does not have permission to user's records"}},
		{"fail - get balance err", state{pointsType: "ADD", errGetBalance: errFail}, want{err: "failed to get latest balance"}},
		{"fail - save points", state{pointsType: "ADD", errSavePoint: errFail}, want{err: "failed to issue points: fail"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := &issuePointsHandlerRequest{
				FromParentUserID: "parent-1",
				Reason:           "not a short reason",
				Points:           10,
				ToChildUserID:    "child-1",
				Type:             models.PointRequestType(c.state.pointsType),
			}

			if c.state.validationError {
				req.Type = models.PointRequestType("blah")
			}

			mockPointsDB := mocks.NewMockIPointsStorage(t)
			mockUserDB := mocks.NewMockIUserStorage(t)

			if !c.state.validationError {
				mockUserDB.EXPECT().ParentHasAccessToChild(mock.Anything, mock.Anything, mock.Anything).Return(!c.state.noAccess, c.state.errHasAccess).Once()
			}

			if !c.state.validationError && c.state.errHasAccess == nil && !c.state.noAccess {
				bal := models.PointBalance{Balance: 10}

				mockPointsDB.EXPECT().GetLatestBalance(mock.Anything, mock.Anything).Return(bal, c.state.errGetBalance).Once()

				if c.state.errGetBalance == nil {
					mockPointsDB.EXPECT().SavePoint(mock.Anything, mock.Anything).Return(c.state.errSavePoint).Once()
				}
			}

			ctrl := PointsController{
				pointsDB: mockPointsDB,
				userDB:   mockUserDB,
			}

			ctx := context.Background()
			res, err := ctrl.handleIssuePoints(ctx, req)
			tests.AssertError(t, err, c.want.err)
			if err == nil {
				assert.Equal(t, c.want.bal, *res.Point.Balance)
				assert.Equal(t, "parent-1", res.Point.Request.DecidedByUserID)
				assert.Equal(t, models.PointRequestDecisionApprove, res.Point.Request.Decision)
				assert.Equal(t, models.PointStatusSettled, res.Point.Status)
				assert.Equal(t, int32(10), res.Summary.Points)
				assert.Equal(t, "parent-1", res.Summary.DecidedByUserID)
				assert.Equal(t, models.PointRequestDecisionApprove, res.Summary.Decision)
			}

			mockPointsDB.AssertExpectations(t)
			mockUserDB.AssertExpectations(t)
		})
	}
}

func Test_validateIssuePoints(t *testing.T) {
	type state struct {
		invalidFromParentUserId bool
		invalidToChildUserId    bool
		invalidType             bool
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
		{"fail - invalid to child user id", state{invalidToChildUserId: true}, want{"missing to_user_id"}},
		{"fail - invalid type", state{invalidType: true}, want{"invalid point request type blah"}},
		{"fail - invalid points - negative", state{pointsAreNegative: true}, want{"failed to validate request: points must be a positive integer"}},
		{"fail - invalid points - zero", state{pointsAreZero: true}, want{"failed to validate request: points must be a positive integer"}},
		{"fail - missing reason", state{missingReason: true}, want{"failed to validate request: reason for issuing points must not be empty"}},
		{"fail - reason too short", state{tooShortReason: true}, want{"failed to validate request: reason for issuing points must not be empty"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := &issuePointsHandlerRequest{
				FromParentUserID: "parent-1",
				ToChildUserID:    "child-1",
				ParentNotes:      "some notes",
				Points:           10,
				Reason:           "I worked hard",
				Type:             models.PointRequestTypeAdd,
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
			if c.state.invalidToChildUserId {
				req.ToChildUserID = ""
			}
			if c.state.invalidType {
				req.Type = models.PointRequestType("blah")
			}

			err := validateIssuePoints(req)
			tests.AssertError(t, err, c.want.err)
		})
	}
}
