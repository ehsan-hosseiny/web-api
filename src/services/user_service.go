package services

import (
	"github.com/ehsan-hosseiny/golang-web-api/api/dto"
	"github.com/ehsan-hosseiny/golang-web-api/common"
	"github.com/ehsan-hosseiny/golang-web-api/config"
	"github.com/ehsan-hosseiny/golang-web-api/data/db"
	"github.com/ehsan-hosseiny/golang-web-api/pkg/logging"
	"gorm.io/gorm"
)

type UserSevice struct {
	logger    logging.Logger
	cfg       *config.Config
	otpSevice *OtpSevice
	database  *gorm.DB
}

func NewUserService(cfg *config.Config) *UserSevice {
	database := db.GetDb()
	logger := logging.NewLogger(cfg)
	return &UserSevice{
		cfg:       cfg,
		database:  database,
		logger:    logger,
		otpSevice: NewOtpService(cfg),
	}
}

func (s *UserSevice) SendOtp(req *dto.GetOtpRequest) error {
	otp := common.GenerateOtp()
	err := s.otpSevice.SetOtp(req.MobileNumber, otp)
	if err != nil {
		return err
	}
	return nil
}
