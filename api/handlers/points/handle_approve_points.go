package points

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sebboness/yektaspoints/handlers"
	"github.com/sebboness/yektaspoints/models"
	"github.com/sebboness/yektaspoints/util"
	apierr "github.com/sebboness/yektaspoints/util/error"
)

type approvePointsRequest struct {
	// From request
	Decision    string `json:"decision"`
	ParentNotes string `json:"parent_notes,omitempty"`
	ChildID     string `json:"user_id"`

	// Set in code
	PointID  string `json:"-"`
	ParentID string `json:"-"`
}

type approvePointsResponse struct {
	Point   models.Point        `json:"point"`
	Summary models.PointSummary `json:"point_summary"`
}

func (c *PointsController) ApprovePointsHandler(cgin *gin.Context) {

	var req approvePointsRequest

	// try to unmarshal from request body
	err := cgin.BindJSON(&req)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal json body: %w", err)
		cgin.JSON(http.StatusBadRequest, handlers.ErrorResult(err))
		return
	}

	// get point id from request
	req.PointID = cgin.Param("point_id")
	if req.PointID == "" {
		apiErr := apierr.New(apierr.InvalidInput).WithError("point_id is a required query parameter")
		cgin.JSON(apiErr.StatusCode(), handlers.ErrorResult(apiErr))
		return
	}

	authInfo := c.AuthContext.GetAuthorizerInfo(cgin)
	req.ParentID = authInfo.GetUserID()

	resp, err := c.handleApprovePoints(cgin.Request.Context(), &req)
	if err != nil {
		logger.Errorf("failed to handle approving points: %v", err.Error())

		if apierr := apierr.IsApiError(err); apierr != nil {
			cgin.JSON(apierr.StatusCode(), handlers.ErrorResult(apierr))
			return
		}

		cgin.JSON(http.StatusInternalServerError, handlers.ErrorResult(err))
		return
	}

	cgin.JSON(http.StatusOK, handlers.SuccessResult(resp))
}

func (c *PointsController) handleApprovePoints(ctx context.Context, req *approvePointsRequest) (approvePointsResponse, error) {
	resp := approvePointsResponse{}

	logger.WithContext(ctx).WithFields(map[string]any{
		"parent_id": req.ParentID,
		"child_id":  req.ChildID,
	})

	if err := validateApprovePoints(req); err != nil {
		return resp, err
	}

	hasAccess, err := c.userDB.ParentHasAccessToChild(ctx, req.ParentID, req.ChildID)
	if err != nil {
		return resp, fmt.Errorf("failed to check user access permissions: %w", err)
	}
	if !hasAccess {
		return resp, apierr.New(apierr.AccessDenied).WithError("requesting user does not have permission to user's records")
	}

	point, err := c.pointsDB.GetPointByID(ctx, req.ChildID, req.PointID)
	if err != nil {
		return resp, fmt.Errorf("failed to get point %v: %w", req.PointID, err)
	}

	if point.UserID != req.ChildID {
		return resp, apierr.New(apierr.BadRequest).WithError(fmt.Sprintf("point user id %v does not match request %v", point.UserID, req.ChildID))
	}
	if point.Status != models.PointStatusWaiting {
		return resp, apierr.New(apierr.BadRequest).WithError(fmt.Sprintf("invalid point status %v", point.Status))
	}

	latestBalance, err := c.pointsDB.GetLatestBalance(ctx, req.ChildID)
	if err != nil {
		return resp, fmt.Errorf("failed to get latest balance: %w", err)
	}

	logger.Infof("latestBalance is %v", latestBalance)

	point.Request.Decision = models.PointRequestDecision(req.Decision)
	point.Request.DecidedByUserID = req.ParentID
	point.Request.DecidedOnStr = util.ToFormattedUTC(time.Now())
	point.Request.ParentNotes = req.ParentNotes
	point.Status = models.PointStatusSettled
	point.UpdatedOnStr = util.ToFormattedUTC(time.Now())

	// Update latest point balance for user
	if point.Request.Decision == models.PointRequestDecisionApprove {
		var newBalance int32

		// TODO
		// Don't allow cashouts if new balance goes below 0

		if models.IsSubtractType(point.Request.Type) {
			newBalance = latestBalance.Balance - point.Points
		} else {
			newBalance = latestBalance.Balance + point.Points
		}

		point.Balance = &newBalance
	}

	if err := c.pointsDB.SavePoint(ctx, point); err != nil {
		return resp, fmt.Errorf("failed to approve point request: %w", err)
	}

	point.ParseTimes()

	resp.Point = point
	resp.Summary = point.ToPointSummary()

	return resp, nil
}

func validateApprovePoints(req *approvePointsRequest) error {
	apierr := apierr.New(fmt.Errorf("%w: failed to validate request", apierr.InvalidInput))

	if req.ParentID == "" {
		apierr.AppendError("missing parent id")
	}

	if !slices.Contains(models.ValidPointRequestDecisions, models.PointRequestDecision(req.Decision)) {
		apierr.AppendError(fmt.Sprintf("invalid decision %v", req.Decision))
	}

	if len(req.ParentNotes) >= 500 {
		apierr.AppendError("parent notes should be no longer than 500 characters")
	}

	if len(apierr.Errors()) > 0 {
		return apierr
	}

	return nil
}
