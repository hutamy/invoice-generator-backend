package controllers

import (
	"net/http"

	"github.com/hutamy/invoice-generator/dto"
	"github.com/hutamy/invoice-generator/services"
	"github.com/hutamy/invoice-generator/utils"
	"github.com/hutamy/invoice-generator/utils/errors"
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

	err := c.authService.SignUp(*req)
	if err != nil {
		if err == errors.ErrUserAlreadyExists {
			return utils.Response(ctx, http.StatusConflict, err.Error(), nil)
		}

		return utils.Response(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.Response(ctx, http.StatusCreated, "User created successfully", nil)
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

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return utils.Response(ctx, http.StatusInternalServerError, errors.ErrFailedGenerateToken.Error(), nil)
	}

	return utils.Response(ctx, http.StatusOK, "Sign In successful", echo.Map{
		"token": token,
		"user": echo.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
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
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}
