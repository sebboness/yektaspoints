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
	"github.com/segmentio/ksuid"
)

type issuePointsHandlerRequest struct {
	Points        int32                   `json:"points"`
	Reason        string                  `json:"reason"`
	ParentNotes   string                  `json:"parent_notes,omitempty"`
	ToChildUserID string                  `json:"to_user_id"`
	Type          models.PointRequestType `json:"type"`

	// set in code
	FromParentUserID string `json:"from_parent_user_id"`
}

type issuePointsHandlerResponse struct {
	Point   models.Point        `json:"point"`
	Summary models.PointSummary `json:"point_summary"`
}

// IssuePointsHandler is called by parent users to issue points. These requests can be point additions, subtractions,
// or cashouts.
func (c *PointsController) IssuePointsHandler(cgin *gin.Context) {

	var req issuePointsHandlerRequest

	// try to unmarshal from request body
	err := cgin.BindJSON(&req)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal json body: %w", err)
		cgin.JSON(http.StatusBadRequest, handlers.ErrorResult(err))
		return
	}

	authInfo := c.AuthContext.GetAuthorizerInfo(cgin)
	req.FromParentUserID = authInfo.GetUserID()

	resp, err := c.handleIssuePoints(cgin.Request.Context(), &req)
	if err != nil {
		if apierr := apierr.IsApiError(err); apierr != nil {
			cgin.JSON(apierr.StatusCode(), handlers.ErrorResult(apierr))
			return
		}

		logger.Errorf("failed to handle issue points: %v", err.Error())
		cgin.JSON(http.StatusInternalServerError, handlers.ErrorResult(err))
		return
	}

	cgin.JSON(http.StatusOK, handlers.SuccessResult(resp))
}

func (c *PointsController) handleIssuePoints(ctx context.Context, req *issuePointsHandlerRequest) (issuePointsHandlerResponse, error) {
	resp := issuePointsHandlerResponse{}

	logger.WithContext(ctx).WithFields(map[string]any{
		"parent_id": req.FromParentUserID,
		"child_id":  req.ToChildUserID,
	})

	if err := validateIssuePoints(req); err != nil {
		return resp, err
	}

	// Ensure parent has access to child points
	hasAccess, err := c.userDB.ParentHasAccessToChild(ctx, req.FromParentUserID, req.ToChildUserID)
	if err != nil {
		return resp, fmt.Errorf("failed to check user access permissions: %w", err)
	}
	if !hasAccess {
		return resp, apierr.New(apierr.AccessDenied).WithError("requesting user does not have permission to user's records")
	}

	// Get latest balance
	latestBalance, err := c.pointsDB.GetLatestBalance(ctx, req.ToChildUserID)
	if err != nil {
		return resp, fmt.Errorf("failed to get latest balance: %w", err)
	}

	now := util.ToFormattedUTC(time.Now())
	newBalance := latestBalance.Balance + req.Points
	if req.Type != models.PointRequestTypeAdd {
		newBalance = latestBalance.Balance - req.Points
	}

	point := models.Point{
		ID:      ksuid.New().String(),
		UserID:  req.ToChildUserID,
		Balance: &newBalance,
		Points:  req.Points,
		Status:  models.PointStatusSettled,
		Request: models.PointRequest{
			Decision:        models.PointRequestDecisionApprove,
			DecidedByUserID: req.FromParentUserID,
			DecidedOnStr:    now,
			ParentNotes:     req.ParentNotes,
			Reason:          req.Reason,
			Type:            req.Type,
		},
		CreatedOnStr: now,
		UpdatedOnStr: now,
	}

	err = c.pointsDB.SavePoint(ctx, point)
	if err != nil {
		return resp, fmt.Errorf("failed to issue points: %w", err)
	}

	point.ParseTimes()
	resp.Point = point
	resp.Summary = point.ToPointSummary()

	return resp, nil
}

func validateIssuePoints(req *issuePointsHandlerRequest) error {
	apierr := apierr.New(fmt.Errorf("%w: failed to validate request", apierr.InvalidInput))

	if req.ToChildUserID == "" {
		apierr.AppendError("missing to_user_id")
	}
	if req.FromParentUserID == "" {
		apierr.AppendError("missing from_parent_user_id")
	}

	if req.Points <= 0 {
		apierr.AppendError("points must be a positive integer")
	}

	if !slices.Contains(models.ValidPointRequestTypes, models.PointRequestType(req.Type)) {
		apierr.AppendError(fmt.Sprintf("invalid point request type %v", req.Type))
	}

	if len(req.ParentNotes) >= 500 {
		apierr.AppendError("parent notes should be no longer than 500 characters")
	}

	// Arbitrary check for some valid reason text
	// TODO: make it better
	if req.Reason == "" || len(req.Reason) <= 5 {
		apierr.AppendError("reason for issuing points must not be empty")
	}

	if len(apierr.Errors()) > 0 {
		return apierr
	}

	return nil
}
