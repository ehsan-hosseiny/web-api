package handlers

import (
	"net/http"

	"github.com/ehsan-hosseiny/golang-web-api/api/helper"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck godoc
// @Summary Heaalth Check
// @Description Health Check
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} helper.BaseHttpResponse{} "Success"
// @Router /v1/health/ [get]
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, helper.GenerateBaseResponse("working!", true, 0))
}
