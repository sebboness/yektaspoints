package points

import (
	"context"
	"net/http"
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

func Test_Controller_GetUserPointsHandler(t *testing.T) {
	type state struct {
		missingUser bool
		err         error
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
		{"happy path", state{}, want{"", http.StatusOK}},
		{"fail - missing user", state{missingUser: true}, want{"user_id is a required query parameter", http.StatusBadRequest}},
		{"fail - internal server error", state{err: errFail}, want{"fail", http.StatusInternalServerError}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			pointsDB := mocks.NewMockIPointsStorage(t)

			ctrl := PointsController{
				pointsDB: pointsDB,
			}

			points := []models.Point{
				{ID: "1", UserID: "a", Points: 1},
				{ID: "2", UserID: "a", Points: 1},
				{ID: "3", UserID: "a", Points: 1},
			}

			if !c.state.missingUser {
				pointsDB.EXPECT().GetPointsByUserID(mock.Anything, mock.Anything, mock.Anything).Return(points, c.state.err).Once()
			}

			ctx := context.Background()

			evt := handlers.MockApiGWEvent

			if c.state.missingUser {
				evt.RequestContext.Authorizer = nil
			}

			ctx = handlers.PrepareAuthorizedContext(ctx, evt)

			w := httptest.NewRecorder()
			cgin, _ := gin.CreateTestContext(w)

			if !c.state.missingUser {
				cgin.AddParam("user_id", "a")
			}

			cgin.Request = httptest.NewRequest("GET", "/", nil).WithContext(ctx)

			ctrl.GetUserPointsHandler(cgin)

			assert.Equal(t, c.want.code, w.Code)
			result := tests.AssertResult(t, w.Body)
			tests.AssertResultError(t, result, c.want.err)

			if c.want.code == 200 {
				assert.NotNil(t, result.Data)
			}

			pointsDB.AssertExpectations(t)
		})
	}
}

func Test_Controller_handleGetUserPoints(t *testing.T) {
	type state struct {
		missingUser  bool
		getPointsErr error
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
		{"fail - missing user ID", state{missingUser: true}, want{"missing user id"}},
		{"fail - get points error", state{getPointsErr: errFail}, want{"failed to get points"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			ctx := context.Background()
			pointsDB := mocks.NewMockIPointsStorage(t)

			ctrl := PointsController{
				pointsDB: pointsDB,
			}

			points := []models.Point{
				{
					ID:     "1",
					UserID: "a",
				},
				{
					ID:     "2",
					UserID: "a",
				},
			}

			if !c.state.missingUser {
				pointsDB.EXPECT().GetPointsByUserID(mock.Anything, mock.Anything, mock.Anything).Return(points, c.state.getPointsErr).Once()
			}

			req := &getUserPointsHandlerRequest{
				UserID: "1",
			}

			if c.state.missingUser {
				req.UserID = ""
			}

			res, err := ctrl.handleGetUserPoints(ctx, req)

			tests.AssertError(t, err, c.want.err)
			if c.want.err == "" {
				assert.Len(t, res.Points, 2)
			}

			pointsDB.AssertExpectations(t)
		})
	}
}
