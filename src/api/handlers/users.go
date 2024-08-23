package handlers

import (
	"net/http"

	"github.com/ehsan-hosseiny/golang-web-api/api/dto"
	"github.com/ehsan-hosseiny/golang-web-api/api/helper"
	"github.com/ehsan-hosseiny/golang-web-api/config"
	"github.com/ehsan-hosseiny/golang-web-api/services"
	"github.com/gin-gonic/gin"
)

type UsersHandler struct {
	service *services.UserSevice
}

func NewUsersHandler(cfg *config.Config) *UsersHandler {
	service := services.NewUserService(cfg)
	return &UsersHandler{service: service}
}

// LoginByUsername godoc
// @Summary LoginByUsername
// @Description Login by user name
// @Tags Users
// @Accept json
// @Produce json
// @Param Request body dto.LoginByUsernameRequest true "LoginByUsernameRequest"
// @Success 200 {object} helper.BaseHttpResponse{} "Success"
// @Success 400 {object} helper.BaseHttpResponse{} "Failed"
// @Success 409 {object} helper.BaseHttpResponse{} "Failed"
// @Router /v1/users/login-by-username [post]
func (h *UsersHandler) LoginByUsername(c *gin.Context) {

	req := new(dto.LoginByUsernameRequest)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, -1, err))
		return
	}
	token, err := h.service.LoginByUserName(req)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, -1, err))
		return
	}
	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(token, true, 0))
}

// RegisterByUsername godoc
// @Summary RegisterByUsername
// @Description Login by user name
// @Tags Users
// @Accept json
// @Produce json
// @Param Request body dto.RegisterUserByUsername true "RegisterUserByUsername"
// @Success 200 {object} helper.BaseHttpResponse{} "Success"
// @Success 400 {object} helper.BaseHttpResponse{} "Failed"
// @Success 409 {object} helper.BaseHttpResponse{} "Failed"
// @Router /v1/users/register-by-username [post]
func (h *UsersHandler) RegisterByUsername(c *gin.Context) {
	req := new(dto.RegisterUserByUsername)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, -1, err))
		return
	}
	err = h.service.RegisterByUsername(req)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, -1, err))
		return
	}
	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(nil, true, 0))

}

// RegisterLoginByMobileNumber godoc
// @Summary RegisterLoginByMobileNumber
// @Description RegisterLoginByMobileNumber
// @Tags Users
// @Accept json
// @Produce json
// @Param Request body dto.RegisterLoginByMobileRequest true "RegisterLoginByMobileRequest"
// @Success 200 {object} helper.BaseHttpResponse{} "Success"
// @Success 400 {object} helper.BaseHttpResponse{} "Failed"
// @Success 409 {object} helper.BaseHttpResponse{} "Failed"
// @Router /v1/users/login-by-mobile [post]
func (h *UsersHandler) RegisterLoginByMobileNumber(c *gin.Context) {
	req := new(dto.RegisterLoginByMobileRequest)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, -1, err))
		return
	}
	token, err := h.service.RegisterLoginByMobileNumber(req)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, -1, err))
		return
	}
	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(token, true, 0))

}

// SendOtp godoc
// @Summary Send otp to user
// @Description Send opt to user
// @Tags Users
// @Accept json
// @Produce json
// @Param Request body dto.GetOtpRequest true "GetOtpRequest"
// @Success 200 {object} helper.BaseHttpResponse{} "Success"
// @Success 400 {object} helper.BaseHttpResponse{} "Failed"
// @Success 409 {object} helper.BaseHttpResponse{} "Failed"
// @Router /v1/users/send-otp [post]
func (h *UsersHandler) SendOtp(c *gin.Context) {
	req := new(dto.GetOtpRequest)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, -1, err))
		return
	}
	otp, err := h.service.SendOtp(req)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, -1, err))
		return
	}
	// Call Internal SMS Service
	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(otp, true, 0))

}
