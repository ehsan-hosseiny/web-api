package routers

import (
	"github.com/ehsan-hosseiny/golang-web-api/api/handlers"
	"github.com/ehsan-hosseiny/golang-web-api/config"
	"github.com/gin-gonic/gin"
)

func User(router *gin.RouterGroup, cfg *config.Config) {
	h := handlers.NewUsersHandler(cfg)
	// router.POST("/send-otp", middlewares.OtpLimiter(cfg), h.SendOtp)
	router.POST("/send-otp", h.SendOtp)
	router.POST("/login-by-username", h.LoginByUsername)
	router.POST("/register-by-username", h.RegisterByUsername)
	router.POST("/login-by-mobile", h.RegisterLoginByMobileNumber)
}
