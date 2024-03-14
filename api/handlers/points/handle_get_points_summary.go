package points

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sebboness/yektaspoints/handlers"
	"github.com/sebboness/yektaspoints/models"
	apierr "github.com/sebboness/yektaspoints/util/error"
)

type getPointsSummaryHandlerRequest struct {
	UserID string `json:"-"`
}

type getPointsSummaryHandlerResponse struct {
	models.UserPoints
}

func (c *PointsController) GetPointsSummaryHandler(cgin *gin.Context) {

	authInfo := handlers.GetAuthorizerInfo(cgin)

	req := &getPointsSummaryHandlerRequest{
		UserID: authInfo.GetUserID(),
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

	now := time.Now().UTC()
	from := now.AddDate(0, -1, 0) // minus one month
	to := now

	attributes := []string{
		"id",
		"updated_on",
		"points",
		"balance",
		"request.decided_by_user_id",
		"request.decision",
		"request.parent_notes",
		"request.reason",
		"request.type",
	}

	filter := models.QueryPointsFilter{
		UpdatedOn:  *models.NewDateFilter().WithRange(from, to),
		Statuses:   []models.PointStatus{models.PointStatusWaiting},
		Attributes: attributes,
	}

	// get most recent points with waiting status (not approved yet by parents)
	pointsWaiting, err := c.pointsDB.GetPointsByUserID(context.Background(), req.UserID, filter)
	if err != nil {
		return resp, fmt.Errorf("failed to get recent point requests: %w", err)
	}

	filter = models.QueryPointsFilter{
		UpdatedOn:  *models.NewDateFilter().WithRange(from, to),
		Statuses:   []models.PointStatus{models.PointStatusSettled},
		Types:      []models.PointRequestType{models.PointRequestTypeAdd, models.PointRequestTypeSubtract},
		Attributes: attributes,
	}

	// get most recent points
	points, err := c.pointsDB.GetPointsByUserID(context.Background(), req.UserID, filter)
	if err != nil {
		return resp, fmt.Errorf("failed to get recent points: %w", err)
	}

	filter = models.QueryPointsFilter{
		UpdatedOn:  *models.NewDateFilter().WithRange(from, to),
		Statuses:   []models.PointStatus{models.PointStatusSettled},
		Types:      []models.PointRequestType{models.PointRequestTypeCashout},
		Attributes: attributes,
	}

	// get most recent cashouts
	cashouts, err := c.pointsDB.GetPointsByUserID(context.Background(), req.UserID, filter)
	if err != nil {
		return resp, fmt.Errorf("failed to get recent cashouts: %w", err)
	}

	resp.RecentRequests = models.ToPointSummaries(pointsWaiting[0:2])
	resp.RecentPoints = models.ToPointSummaries(points[0:2])
	resp.RecentCashouts = models.ToPointSummaries(cashouts[0:2])

	return resp, nil
}
