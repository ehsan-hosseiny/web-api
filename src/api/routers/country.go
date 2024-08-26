package routers

import (
	"github.com/ehsan-hosseiny/golang-web-api/api/handlers"
	"github.com/ehsan-hosseiny/golang-web-api/config"
	"github.com/gin-gonic/gin"
)

func Counntry(r *gin.RouterGroup, cfg *config.Config) {
	h := handlers.NewCountryHandler(cfg)

	r.POST("/", h.Create)
	r.PUT("/:id", h.Update)
	r.DELETE("/:id", h.Delete)
	r.GET("/:id", h.GetById)
	r.POST("/get-by-filter", h.GetByFilter)
}
