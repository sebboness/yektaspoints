package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sebboness/yektaspoints/handlers"
	apierr "github.com/sebboness/yektaspoints/util/error"
)

type userRegisterConfirmRequest struct {
	Username string `json:"username"`
	Code     string `json:"code"`
}

// UserRegisterConfirmHandler confirms a user registration by providing a code that was emailed/SMSed to them
func (c *UserController) UserRegisterConfirmHandler(cgin *gin.Context) {

	var req userRegisterConfirmRequest

	// try to unmarshal from request body
	err := cgin.BindJSON(&req)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal json body: %w", err)
		cgin.JSON(http.StatusBadRequest, handlers.ErrorResult(err))
		return
	}

	err = c.handleUserRegisterConfirm(cgin.Request.Context(), &req)
	if err != nil {
		if apierr := apierr.IsApiError(err); apierr != nil {
			cgin.JSON(apierr.StatusCode(), handlers.ErrorResult(apierr))
			return
		}

		cgin.JSON(http.StatusInternalServerError, handlers.ErrorResult(err))
		return
	}

	cgin.JSON(http.StatusOK, handlers.SuccessResult(nil))
}

func (c *UserController) handleUserRegisterConfirm(ctx context.Context, req *userRegisterConfirmRequest) error {

	if err := validateUserRegisterConfirm(req); err != nil {
		return err
	}

	err := c.auth.ConfirmRegistration(ctx, req.Username, req.Code)

	if err != nil {
		return fmt.Errorf("failed to confirm user registration for '%s': %w", req.Username, err)
	}

	return nil
}

func validateUserRegisterConfirm(req *userRegisterConfirmRequest) error {
	apierr := apierr.New(fmt.Errorf("%w: failed to validate request", apierr.InvalidInput))

	if req.Username == "" {
		apierr.AppendError("missing username")
	}

	if req.Code == "" {
		apierr.AppendError("missing code")
	}

	if len(apierr.Errors()) > 0 {
		return apierr
	}

	return nil
}
