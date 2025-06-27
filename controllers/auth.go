package controllers

import (
	"net/http"
	"time"

	"github.com/hutamy/invoice-generator-backend/dto"
	"github.com/hutamy/invoice-generator-backend/services"
	"github.com/hutamy/invoice-generator-backend/utils"
	"github.com/hutamy/invoice-generator-backend/utils/errors"
	"github.com/labstack/echo/v4"
)

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// @Summary      User Sign Up
// @Description  Register a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.SignUpRequest  true  "Sign Up Request"
// @Success      201   {object}  utils.GenericResponse
// @Failure      400   {object}  utils.GenericResponse
// @Failure      409   {object}  utils.GenericResponse
// @Failure      500   {object}  utils.GenericResponse
// @Router       /v1/public/auth/sign-up [post]
func (c *AuthController) SignUp(ctx echo.Context) error {
	req := new(dto.SignUpRequest)
	if err := ctx.Bind(req); err != nil {
		return utils.Response(ctx, http.StatusBadRequest, errors.ErrBadRequest.Error(), nil)
	}

	if err := ctx.Validate(req); err != nil {
		return utils.Response(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	user, err := c.authService.SignUp(*req)
	if err != nil {
		if err == errors.ErrUserAlreadyExists {
			return utils.Response(ctx, http.StatusConflict, err.Error(), nil)
		}

		return utils.Response(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	accessToken, err := utils.GenerateJWT(user.ID, time.Hour*24) // Token valid for 24 hours
	if err != nil {
		return utils.Response(ctx, http.StatusInternalServerError, errors.ErrFailedGenerateToken.Error(), nil)
	}

	refreshToken, err := utils.GenerateJWT(user.ID, time.Hour*24*7) // Token valid for 7 days
	if err != nil {
		return utils.Response(ctx, http.StatusInternalServerError, errors.ErrFailedGenerateToken.Error(), nil)
	}

	return utils.Response(ctx, http.StatusCreated, "User created successfully", echo.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// @Summary      User Sign In
// @Description  Authenticate user and return JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.SignInRequest  true  "Sign In Request"
// @Success      200   {object}  utils.GenericResponse
// @Failure      400   {object}  utils.GenericResponse
// @Failure      401   {object}  utils.GenericResponse
// @Failure      500   {object}  utils.GenericResponse
// @Router       /v1/public/auth/sign-in [post]
func (c *AuthController) SignIn(ctx echo.Context) error {
	req := new(dto.SignInRequest)
	if err := ctx.Bind(req); err != nil {
		return utils.Response(ctx, http.StatusBadRequest, errors.ErrBadRequest.Error(), nil)
	}

	if err := ctx.Validate(req); err != nil {
		return utils.Response(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	user, err := c.authService.SignIn(req.Email, req.Password)
	if err != nil {
		if err == errors.ErrLoginFailed {
			return utils.Response(ctx, http.StatusUnauthorized, err.Error(), nil)
		}

		return utils.Response(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	accessToken, err := utils.GenerateJWT(user.ID, time.Hour*24) // Token valid for 24 hours
	if err != nil {
		return utils.Response(ctx, http.StatusInternalServerError, errors.ErrFailedGenerateToken.Error(), nil)
	}

	refreshToken, err := utils.GenerateJWT(user.ID, time.Hour*24*7) // Token valid for 7 days
	if err != nil {
		return utils.Response(ctx, http.StatusInternalServerError, errors.ErrFailedGenerateToken.Error(), nil)
	}

	return utils.Response(ctx, http.StatusOK, "Sign In successful", echo.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// @Summary      Get Current User
// @Description  Get details of the authenticated user
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200   {object}  utils.GenericResponse
// @Failure      401   {object}  utils.GenericResponse
// @Failure      404   {object}  utils.GenericResponse
// @Failure      500   {object}  utils.GenericResponse
// @Router       /v1/protected/me [get]
func (c *AuthController) Me(ctx echo.Context) error {
	userID, ok := ctx.Get("user_id").(uint)
	// Check if userID is present in the context
	if !ok {
		return utils.Response(ctx, http.StatusUnauthorized, errors.ErrUnauthorized.Error(), nil)
	}

	user, err := c.authService.GetUserByID(userID)
	if err != nil {
		return utils.Response(ctx, http.StatusInternalServerError, err.Error(), nil)
	}
	if user == nil {
		return utils.Response(ctx, http.StatusNotFound, errors.ErrNotFound.Error(), nil)
	}

	return utils.Response(ctx, http.StatusOK, "User retrieved successfully", echo.Map{
		"id":                  user.ID,
		"name":                user.Name,
		"email":               user.Email,
		"phone":               user.Phone,
		"address":             user.Address,
		"bank_name":           user.BankName,
		"bank_account_number": user.BankAccountNumber,
		"bank_account_name":   user.BankAccountName,
	})
}

// @Summary      Refresh Token
// @Description  Refresh access token using a valid refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.RefreshTokenRequest  true  "Refresh Token Request"
// @Success      200   {object}  utils.GenericResponse
// @Failure      400   {object}  utils.GenericResponse
// @Failure      401   {object}  utils.GenericResponse
// @Failure      500   {object}  utils.GenericResponse
// @Router       /v1/public/auth/refresh-token [post]
func (c *AuthController) RefreshToken(ctx echo.Context) error {
	req := new(dto.RefreshTokenRequest)
	if err := ctx.Bind(req); err != nil {
		return utils.Response(ctx, http.StatusBadRequest, errors.ErrBadRequest.Error(), nil)
	}

	if err := ctx.Validate(req); err != nil {
		return utils.Response(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	claims, err := utils.ParseJWT(req.RefreshToken)
	if err != nil {
		return utils.Response(ctx, http.StatusUnauthorized, errors.ErrUnauthorized.Error(), nil)
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return utils.Response(ctx, http.StatusUnauthorized, errors.ErrUnauthorized.Error(), nil)
	}

	user, err := c.authService.GetUserByID(uint(userID))
	if err != nil {
		return utils.Response(ctx, http.StatusInternalServerError, err.Error(), nil)
	}
	if user == nil {
		return utils.Response(ctx, http.StatusNotFound, errors.ErrNotFound.Error(), nil)
	}

	accessToken, err := utils.GenerateJWT(user.ID, time.Hour*24)
	if err != nil {
		return utils.Response(ctx, http.StatusInternalServerError, errors.ErrFailedGenerateToken.Error(), nil)
	}

	refreshToken, err := utils.GenerateJWT(user.ID, time.Hour*24*7)
	if err != nil {
		return utils.Response(ctx, http.StatusInternalServerError, errors.ErrFailedGenerateToken.Error(), nil)
	}

	return utils.Response(ctx, http.StatusOK, "Token refreshed successfully", echo.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// @Summary      Update User
// @Description  Update user details
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.UpdateUserRequest  true  "Update User Request"
// @Success      200   {object}  utils.GenericResponse
// @Failure      400   {object}  utils.GenericResponse
// @Failure      401   {object}  utils.GenericResponse
// @Failure      404   {object}  utils.GenericResponse
// @Failure      500   {object}  utils.GenericResponse
// @Router       /v1/protected/me [put]
func (c *AuthController) UpdateUser(ctx echo.Context) error {
	req := new(dto.UpdateUserRequest)
	if err := ctx.Bind(req); err != nil {
		return utils.Response(ctx, http.StatusBadRequest, errors.ErrBadRequest.Error(), nil)
	}

	if err := ctx.Validate(req); err != nil {
		return utils.Response(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	userID, ok := ctx.Get("user_id").(uint)
	if !ok {
		return utils.Response(ctx, http.StatusUnauthorized, errors.ErrUnauthorized.Error(), nil)
	}

	req.UserID = userID
	if err := c.authService.UpdateUser(*req); err != nil {
		return utils.Response(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.Response(ctx, http.StatusOK, "User updated successfully", nil)
}
