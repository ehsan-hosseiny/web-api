package routers

import (
	"github.com/ehsan-hosseiny/golang-web-api/api/handlers"
	"github.com/gin-gonic/gin"
)

func Health(r *gin.RouterGroup) {
	h := handlers.NewHealthHandler()
	r.GET("/", h.Health)
}
