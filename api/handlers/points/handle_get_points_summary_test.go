package points

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sebboness/yektaspoints/handlers"
	handlerMocks "github.com/sebboness/yektaspoints/mocks/handlers"
	mocks "github.com/sebboness/yektaspoints/mocks/storage"
	"github.com/sebboness/yektaspoints/models"
	"github.com/sebboness/yektaspoints/util/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Controller_GetPointsSummaryHandler(t *testing.T) {
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

			mockAuthContext := handlerMocks.NewMockAuthContext(t)
			mockUserDB := mocks.NewMockIUserStorage(t)
			mockPointsDB := mocks.NewMockIPointsStorage(t)

			ctrl := PointsController{
				BaseController: handlers.BaseController{
					AuthContext: mockAuthContext,
				},
				pointsDB: mockPointsDB,
				userDB:   mockUserDB,
			}

			points := []models.Point{
				{ID: "1", UserID: "a", Points: 1},
				{ID: "2", UserID: "a", Points: 1},
				{ID: "3", UserID: "a", Points: 1},
			}

			authInfo := handlers.AuthorizerInfo{
				Claims: handlers.DefaultMockAuthClaims,
			}

			if c.state.missingUser {
				authInfo = handlers.AuthorizerInfo{}
			}

			if !c.state.missingUser {
				mockAuthContext.EXPECT().GetAuthorizerInfo(mock.Anything).Return(authInfo)
				mockUserDB.EXPECT().ParentHasAccessToChild(mock.Anything, mock.Anything, mock.Anything).Return(true, nil).Once()
				mockPointsDB.EXPECT().GetPointsByUserID(mock.Anything, mock.Anything, mock.Anything).Return(points, c.state.err).Once()
			}

			w := httptest.NewRecorder()
			cgin, _ := gin.CreateTestContext(w)

			if !c.state.missingUser {
				cgin.AddParam("user_id", "a")
			}

			cgin.Request = httptest.NewRequest("GET", "/", nil)

			ctrl.GetPointsSummaryHandler(cgin)

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

func Test_Controller_handleGetPointsSummary(t *testing.T) {
	type state struct {
		missingUser  bool
		noAccess     bool
		noAccessErr  error
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
		{"fail - check access error", state{noAccessErr: errFail}, want{"failed to check access permissions"}},
		{"fail - no access error", state{noAccess: true}, want{"user does not have permissions to get points from user"}},
		{"fail - get points error", state{getPointsErr: errFail}, want{"failed to get points"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			ctx := context.Background()
			mockPointsDB := mocks.NewMockIPointsStorage(t)
			mockUserDB := mocks.NewMockIUserStorage(t)

			ctrl := PointsController{
				pointsDB: mockPointsDB,
				userDB:   mockUserDB,
			}

			if !c.state.missingUser {
				// setup some mock points
				now := time.Now().UTC()
				points := []models.Point{}
				daysBack := 0
				for {
					balance := int32(20 - daysBack)
					points = append(points, models.Point{
						Status:    "SETTLED",
						Points:    1,
						Balance:   &balance,
						UpdatedOn: now.AddDate(0, 0, -daysBack),
						Request: models.PointRequest{
							Type: "ADD",
						},
					})

					daysBack += 1

					if len(points) >= 10 {
						break
					}
				}

				mockUserDB.EXPECT().ParentHasAccessToChild(mock.Anything, mock.Anything, mock.Anything).Return(!c.state.noAccess, c.state.noAccessErr).Once()

				if !c.state.noAccess && c.state.noAccessErr == nil {
					mockPointsDB.EXPECT().GetPointsByUserID(mock.Anything, mock.Anything, mock.Anything).Return(points, c.state.getPointsErr).Once()
				}
			}

			req := &getPointsSummaryHandlerRequest{
				UserID: "1",
			}

			if c.state.missingUser {
				req.UserID = ""
			}

			res, err := ctrl.handleGetPointsSummary(ctx, req)

			tests.AssertError(t, err, c.want.err)
			if c.want.err == "" {
				assert.Equal(t, int32(20), res.Balance)
				assert.GreaterOrEqual(t, int32(8), res.PointsLast7Days)
				assert.Equal(t, int32(0), res.PointsLostLast7Days)
			}

			mockPointsDB.AssertExpectations(t)
			mockUserDB.AssertExpectations(t)
		})
	}
}

func Test_Controller_mapPointsToSummaries(t *testing.T) {
	type test struct {
		name string
	}

	cases := []test{
		{"happy path"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			ctrl := PointsController{}

			bal := func(v int32) *int32 {
				return &v
			}

			// setup mock points
			now := time.Now()
			from := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, -7)

			p0 := models.Point{ID: "0", Status: "SETTLED", Points: 1, Balance: bal(18), UpdatedOn: now.AddDate(0, 0, 0), Request: models.PointRequest{Type: "ADD"}}
			p1 := models.Point{ID: "1", Status: "SETTLED", Points: 1, Balance: bal(17), UpdatedOn: now.AddDate(0, 0, -1), Request: models.PointRequest{Type: "ADD"}}
			p2 := models.Point{ID: "2", Status: "SETTLED", Points: -1, Balance: bal(16), UpdatedOn: now.AddDate(0, 0, -2), Request: models.PointRequest{Type: "SUBTRACT"}}
			p3 := models.Point{ID: "3", Status: "SETTLED", Points: 1, Balance: bal(17), UpdatedOn: now.AddDate(0, 0, -3), Request: models.PointRequest{Type: "ADD"}}
			p4 := models.Point{ID: "4", Status: "SETTLED", Points: 1, Balance: bal(16), UpdatedOn: now.AddDate(0, 0, -4), Request: models.PointRequest{Type: "ADD"}}
			p5 := models.Point{ID: "5", Status: "SETTLED", Points: -1, Balance: bal(15), UpdatedOn: now.AddDate(0, 0, -5), Request: models.PointRequest{Type: "CASHOUT"}}
			p6 := models.Point{ID: "6", Status: "SETTLED", Points: 1, Balance: bal(16), UpdatedOn: now.AddDate(0, 0, -6), Request: models.PointRequest{Type: "ADD"}}
			p7 := models.Point{ID: "7", Status: "WAITING", Points: 1, UpdatedOn: now.AddDate(0, 0, -7), Request: models.PointRequest{Type: "ADD"}}
			p8 := models.Point{ID: "8", Status: "SETTLED", Points: -1, Balance: bal(15), UpdatedOn: now.AddDate(0, 0, -8), Request: models.PointRequest{Type: "SUBTRACT"}}
			p9 := models.Point{ID: "9", Status: "SETTLED", Points: -1, Balance: bal(16), UpdatedOn: now.AddDate(0, 0, -9), Request: models.PointRequest{Type: "CASHOUT"}}

			points := []models.Point{p0, p1, p2, p3, p4, p5, p6, p7, p8, p9}

			up := &models.UserPoints{}
			ctrl.mapPointsToSummaries(up, from, points)

			assert.Equal(t, int32(18), up.Balance)
			assert.Equal(t, int32(4), up.PointsLast7Days)
			assert.Equal(t, int32(-1), up.PointsLostLast7Days)

			assert.Len(t, up.RecentPoints, 3)
			assert.Len(t, up.RecentRequests, 1)
			assert.Len(t, up.RecentCashouts, 2)

			// assert recent points
			assert.Equal(t, "0", up.RecentPoints[0].ID)
			assert.Equal(t, "1", up.RecentPoints[1].ID)
			assert.Equal(t, "2", up.RecentPoints[2].ID)

			// assert open requests
			assert.Equal(t, "7", up.RecentRequests[0].ID)

			// assert cashouts
			assert.Equal(t, "5", up.RecentCashouts[0].ID)
			assert.Equal(t, "9", up.RecentCashouts[1].ID)
		})
	}
}
