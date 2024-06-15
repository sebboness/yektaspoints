package points

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sebboness/yektaspoints/handlers"
	"github.com/sebboness/yektaspoints/models"
	"github.com/sebboness/yektaspoints/util"
	apierr "github.com/sebboness/yektaspoints/util/error"
	"github.com/sebboness/yektaspoints/util/log"
)

type getPointsSummaryHandlerRequest struct {
	// UserID is for the user that owns the points (child)
	UserID string `json:"-"`

	// User ID that makes the request (parent)
	RequestingUserID string `json:"-"`
}

type getPointsSummaryHandlerResponse struct {
	models.UserPoints
}

func (c *PointsController) GetPointsSummaryHandler(cgin *gin.Context) {

	userID := cgin.Param("user_id")
	if userID == "" {
		apiErr := apierr.New(apierr.InvalidInput).WithError("user_id is a required query parameter")
		cgin.JSON(apiErr.StatusCode(), handlers.ErrorResult(apiErr))
		return
	}

	authInfo := c.AuthContext.GetAuthorizerInfo(cgin)

	req := &getPointsSummaryHandlerRequest{
		UserID:           userID,
		RequestingUserID: authInfo.GetUserID(),
	}

	resp, err := c.handleGetPointsSummary(cgin.Request.Context(), req)
	if err != nil {
		if apierr := apierr.IsApiError(err); apierr != nil {
			cgin.JSON(apierr.StatusCode(), handlers.ErrorResult(apierr))
			return
		}

		cgin.JSON(http.StatusInternalServerError, handlers.ErrorResult(err))
		return
	}

	cgin.JSON(http.StatusOK, handlers.SuccessResult(resp))
}

func (c *PointsController) handleGetPointsSummary(ctx context.Context, req *getPointsSummaryHandlerRequest) (getPointsSummaryHandlerResponse, error) {
	resp := getPointsSummaryHandlerResponse{}

	if req.UserID == "" {
		return resp, apierr.New(apierr.AccessDenied).WithError("missing user id")
	}

	// check that the requested user (a parent) has access to data owned by user (a child)
	hasAccess, err := c.userDB.ParentHasAccessToChild(ctx, req.RequestingUserID, req.UserID)
	if err != nil {
		return resp, fmt.Errorf("failed to check access permissions: %w", err)
	}
	if !hasAccess {
		logger.WithFields(map[string]any{
			"requesting_user_id": req.RequestingUserID,
			"user_id":            req.UserID,
		})
		return resp, fmt.Errorf("user does not have permissions to get points from user: %w", err)
	}

	now := time.Now().UTC()
	from := now.AddDate(0, 0, -14)   // minus two weeks
	weekAgo := now.AddDate(0, 0, -7) // minus one week
	to := now

	attributes := []string{
		"id",
		"updated_on",
		"points",
		"balance",
		"status",
		"request.decided_by_user_id",
		"request.decision",
		"request.parent_notes",
		"request.reason",
		"request.type",
	}

	filter := models.QueryPointsFilter{
		UpdatedOn: *models.NewDateFilter().WithRange(from, to),
		Statuses: []models.PointStatus{
			models.PointStatusSettled,
			models.PointStatusWaiting,
		},
		Types: []models.PointRequestType{
			models.PointRequestTypeCashout,
			models.PointRequestTypeAdd,
			models.PointRequestTypeSubtract,
		},
		Attributes: attributes,
	}

	// get all points with filters applied from 2 weeks ago
	points, err := c.pointsDB.GetPointsByUserID(ctx, req.UserID, filter)
	if err != nil {
		return resp, fmt.Errorf("failed to get points: %w", err)
	}

	// map all points to user point summaries
	// the weekAgo date will summarize point amounts from last 7 days.
	c.mapPointsToSummaries(&resp.UserPoints, weekAgo, points)

	logger := log.Get()
	logger.WithContext(ctx).WithFields(map[string]any{
		"dt_from":     util.ToFormatted(from),
		"dt_to":       util.ToFormatted(to),
		"dt_weekago":  util.ToFormatted(weekAgo),
		"user_points": resp.UserPoints,
		"points_len":  len(points),
		"user_id":     req.UserID,
	}).Infof("retrieved user point summaries")

	return resp, nil
}

func (c *PointsController) mapPointsToSummaries(up *models.UserPoints, recentFromDate time.Time, points []models.Point) {

	unsettled := []models.PointSummary{}
	settled := []models.PointSummary{}
	cashouts := []models.PointSummary{}

	// user's points is the balance value in the most recent settled point object
	for _, p := range points {
		// first settled point is latest
		if up.Balance == 0 && p.Status == models.PointStatusSettled && p.Balance != nil {
			up.Balance = *p.Balance
		}

		// sum up points after given recentFromDate
		if p.UpdatedOn.Compare(recentFromDate) >= 0 &&
			p.Status == models.PointStatusSettled &&
			p.Request.Type != models.PointRequestTypeCashout {

			up.PointsLast7Days += p.Points

			if p.Points < 0 {
				up.PointsLostLast7Days += p.Points
			}
		}

		// unsettled points
		if len(unsettled) < 3 && p.Status == models.PointStatusWaiting {
			unsettled = append(unsettled, p.ToPointSummary())
		}
		// settled points
		if len(settled) < 3 && p.Status == models.PointStatusSettled && p.Request.Type != models.PointRequestTypeCashout {
			settled = append(settled, p.ToPointSummary())
		}
		// cashouts
		if len(cashouts) < 3 && p.Status == models.PointStatusSettled && p.Request.Type == models.PointRequestTypeCashout {
			cashouts = append(cashouts, p.ToPointSummary())
		}
	}

	up.RecentRequests = unsettled
	up.RecentPoints = settled
	up.RecentCashouts = cashouts
}
