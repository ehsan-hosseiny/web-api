package api

import (
	"fmt"

	"github.com/ehsan-hosseiny/golang-web-api/api/middlewares"
	"github.com/ehsan-hosseiny/golang-web-api/api/routers"
	"github.com/ehsan-hosseiny/golang-web-api/config"
	"github.com/ehsan-hosseiny/golang-web-api/docs"
	"github.com/ehsan-hosseiny/golang-web-api/pkg/logging"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var logger = logging.NewLogger(config.GetConfig())

func InitServer(cfg *config.Config) {
	r := gin.New()

	RegisterSwagger(r, cfg)
	RegisterRoutes(r, cfg)

	r.Use(gin.Logger(), gin.CustomRecovery(middlewares.ErrorHandler))
	r.Use(middlewares.DefaultStructuredLogger(cfg))
	r.Use(middlewares.Cors(cfg))

	// api := r.Group("/api")

	// v1 := api.Group("/v1/")
	// {
	// 	health := v1.Group("/health", middlewares.LimitByRequest())
	// 	test_router := v1.Group("/test")

	// 	users := v1.Group("/users", middlewares.LimitByRequest())

	// 	routers.Health(health)
	// 	routers.TestRouter(test_router)
	// 	routers.User(users, cfg)
	// }

	r.Run(fmt.Sprintf(":%s", cfg.Server.Port))
	err := r.Run(fmt.Sprintf(":%s", cfg.Server.Port))
	if err != nil {
		logger.Fatal(logging.General, logging.Startup, err.Error(), nil)
	}
}

func RegisterRoutes(r *gin.Engine, cfg *config.Config) {
	api := r.Group("/api")

	v1 := api.Group("/v1")
	{
		// Test
		health := v1.Group("/health")
		test_router := v1.Group("/test")

		// User
		users := v1.Group("/users")

		// Base
		countries := v1.Group("/countries", middlewares.Authentication(cfg), middlewares.Authorization([]string{"admin"}))
		cities := v1.Group("/cities", middlewares.Authentication(cfg), middlewares.Authorization([]string{"admin"}))
		files := v1.Group("/files", middlewares.Authentication(cfg), middlewares.Authorization([]string{"admin"}))

		// Property
		properties := v1.Group("/properties", middlewares.Authentication(cfg), middlewares.Authorization([]string{"admin"}))
		propertyCategories := v1.Group("/property-categories", middlewares.Authentication(cfg), middlewares.Authorization([]string{"admin"}))

		// Test
		routers.Health(health)
		routers.TestRouter(test_router)

		// User
		routers.User(users, cfg)

		// Base
		routers.Country(countries, cfg)
		routers.City(cities, cfg)
		routers.File(files, cfg)

		// Property
		routers.Property(properties, cfg)
		routers.PropertyCategory(propertyCategories, cfg)

	}

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
