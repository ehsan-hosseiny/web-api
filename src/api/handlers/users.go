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
	err = h.service.SendOtp(req)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, -1, err))
		return
	}
	// Call Internal SMS Service
	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(nil, true, 0))

}
