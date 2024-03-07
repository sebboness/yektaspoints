package family

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sebboness/yektaspoints/handlers"
	"github.com/sebboness/yektaspoints/models"
	apierr "github.com/sebboness/yektaspoints/util/error"
	"github.com/sebboness/yektaspoints/util/log"
)

type getFamilyHandlerRequest struct {
	FamilyID string
	UserID   string
}

type getFamilyHandlerResponse struct {
	Family models.Family `json:"family"`
}

func (c *FamilyController) GetFamilyHandler(cgin *gin.Context) {

	familyID, present := cgin.GetQuery("family_id")
	if !present {
		apiErr := apierr.New(fmt.Errorf("%w: family_id is a required query parameter", apierr.InvalidInput))
		cgin.JSON(apiErr.StatusCode(), handlers.ErrorResult(apiErr))
		return
	}

	authInfo := handlers.GetAuthorizerInfo(cgin)
	req := &getFamilyHandlerRequest{
		UserID:   authInfo.GetUserID(),
		FamilyID: familyID,
	}

	resp, err := c.handleGetFamily(cgin.Request.Context(), req)
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

func (c *FamilyController) handleGetFamily(ctx context.Context, req *getFamilyHandlerRequest) (getFamilyHandlerResponse, error) {
	resp := getFamilyHandlerResponse{}
	logger := log.Get().AddFields(map[string]any{
		"user_id":   req.UserID,
		"family_id": req.FamilyID,
	})

	familyUsers, err := c.familyDB.GetFamilyUsers(ctx, req.FamilyID)
	if err != nil {
		logger.WithFields(map[string]any{"error": err.Error()}).Errorf("failed to get family users")
		return resp, fmt.Errorf("failed to get family users: %w", err)
	}

	userIsPartOfFamily := false
	userIds := make([]string, len(familyUsers))
	for idx, fu := range familyUsers {
		userIds[idx] = fu.UserID
		if req.UserID == fu.UserID {
			userIsPartOfFamily = true
		}
	}

	if !userIsPartOfFamily {
		logger.Errorf("user is not part of family")
		return resp, apierr.New(fmt.Errorf("%w: user is not part of family", apierr.AccessDenied))
	}

	family, err := c.familyDB.GetFamilyMembersByUserIDs(ctx, req.FamilyID, userIds)
	if err != nil {
		logger.WithFields(map[string]any{"error": err.Error()}).Errorf("failed to get family")
		return resp, fmt.Errorf("failed to get family: %w", err)
	}

	resp.Family = family
	return resp, nil
}
