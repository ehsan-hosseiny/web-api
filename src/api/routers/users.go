package routers

import (
	"github.com/ehsan-hosseiny/golang-web-api/api/handlers"
	"github.com/ehsan-hosseiny/golang-web-api/config"
	"github.com/gin-gonic/gin"
)

func User(router *gin.RouterGroup, cfg *config.Config) {
	h := handlers.NewUsersHandler(cfg)
	router.POST("/send-otp", h.SendOtp)
}
