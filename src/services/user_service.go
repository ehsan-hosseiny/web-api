package services

import (
	"github.com/ehsan-hosseiny/golang-web-api/api/dto"
	"github.com/ehsan-hosseiny/golang-web-api/common"
	"github.com/ehsan-hosseiny/golang-web-api/config"
	"github.com/ehsan-hosseiny/golang-web-api/constants"
	"github.com/ehsan-hosseiny/golang-web-api/data/db"
	"github.com/ehsan-hosseiny/golang-web-api/data/models"
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

func (s *UserSevice) existsByEmail(email string) (bool, error) {
	var exists bool
	if err := s.database.Model(&models.User{}).
		Select("count(*)>0").
		Where("email = ?", email).
		Find(&exists).Error; err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return false, err
	}
	return exists, nil
}

func (s *UserSevice) existsByUsername(username string) (bool, error) {
	var exists bool
	if err := s.database.Model(&models.User{}).
		Select("count(*)>0").
		Where("username = ?", username).
		Find(&exists).Error; err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return false, err
	}
	return exists, nil
}

func (s *UserSevice) existsByMobileNumber(mobileNumber string) (bool, error) {
	var exists bool
	if err := s.database.Model(&models.User{}).
		Select("count(*)>0").
		Where("mobile_number = ?", mobileNumber).
		Find(&exists).Error; err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return false, err
	}
	return exists, nil
}

func (s *UserSevice) getDefaultRole() (roleId int, err error) {

	if err := s.database.Model(&models.Role{}).
		Select("id").
		Where("name = ?", constants.DefaultRoleName).
		First(&roleId).Error; err != nil {
		return 0, err
	}
	return roleId, nil
}
