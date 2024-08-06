package api

import (
	"fmt"

	"github.com/ehsan-hosseiny/golang-web-api/api/middlewares"
	"github.com/ehsan-hosseiny/golang-web-api/api/routers"
	"github.com/ehsan-hosseiny/golang-web-api/config"
	"github.com/ehsan-hosseiny/golang-web-api/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitServer(cfg *config.Config) {
	r := gin.New()

	RegisterSwagger(r, cfg)

	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middlewares.DefaultStructuredLogger(cfg))
	r.Use(middlewares.Cors(cfg))

	api := r.Group("/api")

	v1 := api.Group("/v1/")
	{
		health := v1.Group("/health", middlewares.LimitByRequest())
		test_router := v1.Group("/test")

		routers.Health(health)
		routers.TestRouter(test_router)
	}

	r.Run(fmt.Sprintf(":%s", cfg.Server.Port))
}

func RegisterSwagger(r *gin.Engine, cfg *config.Config) {
	docs.SwaggerInfo.Title = "golang web api documentation"
	docs.SwaggerInfo.Description = "golang web api documentation"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", cfg.Server.Port)
	docs.SwaggerInfo.Schemes = []string{"http"}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
