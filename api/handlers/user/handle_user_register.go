package user

import (
	"context"
	"fmt"
	"net/http"
	"net/mail"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sebboness/yektaspoints/handlers"
	"github.com/sebboness/yektaspoints/models"
	"github.com/sebboness/yektaspoints/util"
	"github.com/sebboness/yektaspoints/util/auth"
	apierr "github.com/sebboness/yektaspoints/util/error"
	"github.com/sebboness/yektaspoints/util/log"
)

type userRegisterRequest struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Email           string `json:"email"`
	Name            string `json:"name"`
}

type userRegisterResponse struct {
	auth.UserRegisterResult
}

// UserRegisterHandler registers a new user
func (c *UserController) UserRegisterHandler(cgin *gin.Context) {

	var req userRegisterRequest

	// try to unmarshal from request body
	err := cgin.BindJSON(&req)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal json body: %w", err)
		cgin.JSON(http.StatusBadRequest, handlers.ErrorResult(err))
		return
	}

	resp, err := c.handleUserRegister(cgin.Request.Context(), &req)
	if err != nil {
		if apierr := apierr.IsApiError(err); apierr != nil {
			cgin.JSON(apierr.StatusCode(), handlers.ErrorResult(apierr))
			return
		}

		cgin.JSON(http.StatusInternalServerError, handlers.ErrorResult(err))
		return
	}

	cgin.JSON(http.StatusCreated, handlers.SuccessResult(resp))
}

func (c *UserController) handleUserRegister(ctx context.Context, req *userRegisterRequest) (userRegisterResponse, error) {
	resp := userRegisterResponse{}

	if err := validateUserRegister(req); err != nil {
		return resp, err
	}

	logger := log.Get().AddField("username", req.Username)

	// register user with cognito
	result, err := c.auth.Register(ctx, auth.UserRegisterRequest{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Name:     req.Name,
	})

	if err != nil {
		logger.WithContext(ctx).WithField("error", err.Error()).Errorf("failed to register user '%s'", req.Username)
		return resp, fmt.Errorf("failed to register user '%s': %w", req.Username, err)
	}

	// also add a separate user record in dynamodb
	user := models.User{
		UserID:       result.UserID,
		Username:     req.Username,
		Email:        req.Email,
		Name:         req.Name,
		Status:       models.UserStatusUnverified,
		CreatedOnStr: util.ToFormattedUTC(time.Now()),
		UpdatedOnStr: util.ToFormattedUTC(time.Now()),
		FamilyIDs:    []string{},
		Roles:        []string{},
	}

	if err := c.userDB.SaveUser(ctx, user); err != nil {
		logger.WithContext(ctx).WithFields(map[string]any{
			"user_id": result.UserID,
			"error":   err.Error(),
		}).Errorf("failed to store new user '%s'", req.Username)
		return resp, fmt.Errorf("failed to save new user '%s': %w", req.Username, err)
	}

	resp.UserRegisterResult = result
	return resp, nil
}

func validateUserRegister(req *userRegisterRequest) error {
	apierr := apierr.New(fmt.Errorf("%w: failed to validate request", apierr.InvalidInput))

	if len(req.Username) < 4 {
		apierr.AppendError("username must be at least 4 characters long")
	}

	if _, err := mail.ParseAddress(req.Email); err != nil {
		apierr.AppendError("email must be a valid email address")
	}

	if len(req.Name) < 2 {
		apierr.AppendError("name must be at least 2 characters long")
	}

	pwResult := auth.ValidatePassword(req.Password)
	if !pwResult.WithinLength {
		apierr.AppendError("password must be within 8 and 256 characters in length")
	}
	if !pwResult.Lower {
		apierr.AppendError("password must have at least one lower case letter")
	}
	if !pwResult.Upper {
		apierr.AppendError("password must have at least one upper case letter")
	}
	if !pwResult.Number {
		apierr.AppendError("password must have at least one digit")
	}
	if !pwResult.Special {
		apierr.AppendError("password must have at least one special character")
	}

	if req.Password != "" && req.Password != req.ConfirmPassword {
		apierr.AppendError("confirm password does not match password")
	}

	if len(apierr.Errors()) > 0 {
		return apierr
	}

	return nil
}
