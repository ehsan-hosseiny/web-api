package main

import (
	"github.com/ehsan-hosseiny/golang-web-api/api"
	"github.com/ehsan-hosseiny/golang-web-api/config"
	"github.com/ehsan-hosseiny/golang-web-api/data/cache"
	"github.com/ehsan-hosseiny/golang-web-api/data/db"
	"github.com/ehsan-hosseiny/golang-web-api/data/db/migrations"
	"github.com/ehsan-hosseiny/golang-web-api/pkg/logging"
)

// @securityDefinitions.apikey AuthBearer
// @in header
// @name Authorization
func main() {
	cfg := config.GetConfig()
	logger := logging.NewLogger(cfg)

	err := cache.InitRedis(cfg)
	defer cache.CloseRedis()
	if err != nil {
		logger.Fatal(logging.Redis, logging.Startup, err.Error(), nil)
	}

	err = db.InitDb(cfg)
	defer db.Closedb()
	if err != nil {
		logger.Fatal(logging.Postgres, logging.Startup, err.Error(), nil)
	}
	migrations.Up_1()

	api.InitServer(cfg)

}
