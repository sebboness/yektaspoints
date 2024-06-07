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
	mocks "github.com/sebboness/yektaspoints/mocks/storage"
	"github.com/sebboness/yektaspoints/models"
	"github.com/sebboness/yektaspoints/util/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Controller_ApprovePointsHandler(t *testing.T) {
	type state struct {
		validationError bool
		invalidBody     bool
		missingUserId   bool
		missingPointId  bool
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
		{"fail - missing user id param", state{missingUserId: true}, want{"user_id is a required query parameter", 400}},
		{"fail - missing point id param", state{missingPointId: true}, want{"point_id is a required query parameter", 400}},
		{"fail - invalid body", state{invalidBody: true}, want{"failed to unmarshal json body", 400}},
		{"fail - validation error", state{validationError: true}, want{"invalid input", 400}},
		{"fail - internal server error", state{errSavePoint: errors.New("fail")}, want{"fail", 500}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			req := &approvePointsRequest{
				PointID:  "point-1",
				Decision: "APPROVE",
			}

			if c.state.validationError {
				req.Decision = "blah"
			}

			evtBody, _ := json.Marshal(req)
			evtBodyStr := string(evtBody)

			mockPointsDB := mocks.NewMockIPointsStorage(t)
			mockUserDB := mocks.NewMockIUserStorage(t)

			passedInitialChecks := !c.state.invalidBody && !c.state.missingPointId && !c.state.missingUserId && !c.state.validationError

			if c.state.invalidBody {
				evtBodyStr = `{"user_id":`
			}

			if passedInitialChecks {
				point := models.Point{UserID: "child-1"}
				mockUserDB.EXPECT().ParentHasAccessToChild(mock.Anything, mock.Anything, mock.Anything).Return(true, nil).Once()
				mockPointsDB.EXPECT().GetPointByID(mock.Anything, mock.Anything, mock.Anything).Return(point, nil).Once()
				mockPointsDB.EXPECT().SavePoint(mock.Anything, mock.Anything).Return(c.state.errSavePoint).Once()
			}

			ctrl := PointsController{
				pointsDB: mockPointsDB,
				userDB:   mockUserDB,
			}

			ctx := handlers.PrepareAuthorizedContext(context.Background(), handlers.MockApiGWEvent)

			w := httptest.NewRecorder()

			cgin, _ := gin.CreateTestContext(w)

			if !c.state.missingUserId {
				cgin.AddParam("user_id", "child-1")
			}
			if !c.state.missingPointId {
				cgin.AddParam("point_id", "point-1")
			}

			cgin.Request = httptest.NewRequest("POST", "/", bytes.NewReader([]byte(evtBodyStr))).WithContext(ctx)

			handlers.PrepareAuthorizedContext(ctx, handlers.MockApiGWEvent)

			ctrl.ApprovePointsHandler(cgin)

			assert.Equal(t, c.want.code, w.Code)
			result := tests.AssertResult(t, w.Body)
			tests.AssertResultError(t, result, c.want.err)

			if c.want.code == 200 {
				assert.NotNil(t, result.Data)
			}

			mockPointsDB.AssertExpectations(t)
			mockUserDB.AssertExpectations(t)
		})
	}
}

func Test_Controller_handleApprovePoints(t *testing.T) {
	type state struct {
		validationError     bool
		noAccess            bool
		pointUserIdMismatch bool
		errHasAccess        error
		errGetPoint         error
		errSavePoint        error
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
		{"fail - has access err", state{errHasAccess: errFail}, want{"failed to check user access permissions: fail"}},
		{"fail - no access", state{noAccess: true}, want{"requesting user does not have permission to user's records"}},
		{"fail - get point err", state{errGetPoint: errFail}, want{"failed to get point point-1: fail"}},
		{"fail - point user id mistmatch", state{pointUserIdMismatch: true}, want{"point user id user-1 does not match request child-1"}},
		{"fail - save points", state{errSavePoint: errFail}, want{"failed to approve point request: fail"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := &approvePointsRequest{
				ParentID: "parent-1",
				PointID:  "point-1",
				ChildID:  "child-1",
				Decision: "APPROVE",
			}

			if c.state.validationError {
				req.Decision = "blah"
			}

			mockPointsDB := mocks.NewMockIPointsStorage(t)
			mockUserDB := mocks.NewMockIUserStorage(t)

			if !c.state.validationError {
				mockUserDB.EXPECT().ParentHasAccessToChild(mock.Anything, mock.Anything, mock.Anything).Return(!c.state.noAccess, c.state.errHasAccess).Once()
			}

			if !c.state.validationError && c.state.errHasAccess == nil && !c.state.noAccess {
				point := models.Point{UserID: "child-1"}
				if c.state.pointUserIdMismatch {
					point.UserID = "user-1"
				}
				mockPointsDB.EXPECT().GetPointByID(mock.Anything, mock.Anything, mock.Anything).Return(point, c.state.errGetPoint).Once()
			}

			if !c.state.validationError && c.state.errGetPoint == nil && c.state.errHasAccess == nil && !c.state.noAccess && !c.state.pointUserIdMismatch {
				mockPointsDB.EXPECT().SavePoint(mock.Anything, mock.Anything).Return(c.state.errSavePoint).Once()
			}

			ctrl := PointsController{
				pointsDB: mockPointsDB,
				userDB:   mockUserDB,
			}

			ctx := context.Background()
			res, err := ctrl.handleApprovePoints(ctx, req)
			tests.AssertError(t, err, c.want.err)
			if err == nil {
				assert.Equal(t, "parent-1", res.Point.Request.DecidedByUserID)
				assert.Equal(t, models.PointRequestDecisionApprove, res.Point.Request.Decision)
				assert.Equal(t, "parent-1", res.Summary.DecidedByUserID)
				assert.Equal(t, models.PointRequestDecisionApprove, res.Summary.Decision)
			}

			mockPointsDB.AssertExpectations(t)
			mockUserDB.AssertExpectations(t)
		})
	}
}

func Test_validateApprovePoints(t *testing.T) {
	type state struct {
		missingDecision bool
		missingParentID bool
		noNotes         bool
		invalidDecision bool
		tooLongNotes    bool
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
		{"happy path - no notes", state{noNotes: true}, want{}},
		{"fail - missing decision", state{invalidDecision: true}, want{"failed to validate request: invalid decision"}},
		{"fail - invalid decision", state{invalidDecision: true}, want{"failed to validate request: invalid decision"}},
		{"fail - missing parent id", state{missingParentID: true}, want{"failed to validate request: missing parent id"}},
		{"fail - notes too long", state{tooLongNotes: true}, want{"failed to validate request: parent notes should be no longer than 500 characters"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := &approvePointsRequest{
				ParentID:    "123",
				Decision:    "APPROVE",
				ParentNotes: "I worked hard",
			}

			if c.state.noNotes {
				req.ParentNotes = ""
			}
			if c.state.missingParentID {
				req.ParentID = ""
			}
			if c.state.tooLongNotes {
				req.ParentNotes = "hello"
			}
			if c.state.missingDecision {
				req.Decision = ""
			}
			if c.state.invalidDecision {
				req.Decision = "blah"
			}

			err := validateApprovePoints(req)
			if err != nil {
				assert.Contains(t, err.Error(), c.want.err)
			}
		})
	}
}
