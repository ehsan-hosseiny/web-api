package services

import (
	"github.com/ehsan-hosseiny/golang-web-api/api/dto"
	"github.com/ehsan-hosseiny/golang-web-api/common"
	"github.com/ehsan-hosseiny/golang-web-api/config"
	"github.com/ehsan-hosseiny/golang-web-api/constants"
	"github.com/ehsan-hosseiny/golang-web-api/data/db"
	"github.com/ehsan-hosseiny/golang-web-api/data/models"
	"github.com/ehsan-hosseiny/golang-web-api/pkg/logging"
	"github.com/ehsan-hosseiny/golang-web-api/pkg/service_errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserSevice struct {
	logger       logging.Logger
	cfg          *config.Config
	otpSevice    *OtpSevice
	tokenService *TokenService
	database     *gorm.DB
}

func NewUserService(cfg *config.Config) *UserSevice {
	database := db.GetDb()
	logger := logging.NewLogger(cfg)
	return &UserSevice{
		cfg:          cfg,
		database:     database,
		logger:       logger,
		otpSevice:    NewOtpService(cfg),
		tokenService: NewTokenService(cfg),
	}
}

func (s *UserSevice) RegisterByUsername(req *dto.RegisterUserByUsername) error {
	u := models.User{Username: req.Username, FirstName: req.FirstName, LastName: req.LastName,
		Email: req.Email}

	exists, err := s.existsByEmail(req.Email)
	if err != nil {
		return err
	}
	if exists {
		return &service_errors.ServiceError{EndUserMessage: service_errors.EmailExists}
	}

	exists, err = s.existsByUsername(req.Username)
	if err != nil {
		return err
	}
	if exists {
		return &service_errors.ServiceError{EndUserMessage: service_errors.UsernameExists}
	}

	bp := []byte(req.Password)
	hp, err := bcrypt.GenerateFromPassword(bp, bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error(logging.General, logging.HashPassword, err.Error(), nil)
		return err
	}
	u.Password = string(hp)
	roleId, err := s.getDefaultRole()
	if err != nil {
		s.logger.Error(logging.Postgres, logging.DefaultRoleNotFound, err.Error(), nil)
		return err
	}

	// Start of transaction
	tx := s.database.Begin()
	err = tx.Create(&u).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Rollback, err.Error(), nil)
		return err
	}
	err = tx.Create(&models.UserRole{RoleId: roleId, UserId: u.Id}).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Rollback, err.Error(), nil)
		return err
	}
	tx.Commit()
	return nil

}

func (s *UserSevice) LoginByUserName(req *dto.LoginByUsernameRequest) (*dto.TokenDetail, error) {
	var user models.User
	err := s.database.Model(&models.User{}).
		Where("username = ?", req.Username).
		Preload("UserRoles", func(tx *gorm.DB) *gorm.DB {
			return tx.Preload("Role")
		}).
		Find(&user).Error
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, err
	}

	tdto := tokenDto{UserId: user.Id, FirstName: user.FirstName, LastName: user.LastName,
		Email: user.Email, MobileNumber: user.MobileNumber}
	if len(*user.UserRoles) > 0 {
		for _, userRole := range *user.UserRoles {
			tdto.Roles = append(tdto.Roles, userRole.Role.Name)
		}
	}
	token, err := s.tokenService.GenerateToken(&tdto)
	if err != nil {
		return nil, err
	}
	return token, nil

}

func (s *UserSevice) RegisterLoginByMobileNumber(req *dto.RegisterLoginByMobileRequest) (*dto.TokenDetail, error) {
	err := s.otpSevice.ValidatOtp(req.MobileNumber, req.Otp)
	if err != nil {
		return nil, err
	}
	exists, err := s.existsByMobileNumber(req.MobileNumber)
	if err != nil {
		return nil, err
	}

	u := models.User{MobileNumber: req.MobileNumber, Username: req.MobileNumber}

	if exists {
		// Login method
		var user models.User
		err = s.database.Model(&models.User{}).
			Where("username = ?", u.Username).
			Preload("UserRoles", func(tx *gorm.DB) *gorm.DB {
				return tx.Preload("Role")
			}).
			Find(&user).Error
		if err != nil {
			return nil, err
		}
		tdto := tokenDto{UserId: user.Id, FirstName: user.FirstName, LastName: user.LastName,
			Email: user.Email, MobileNumber: user.MobileNumber}
		if len(*user.UserRoles) > 0 {
			for _, userRole := range *user.UserRoles {
				tdto.Roles = append(tdto.Roles, userRole.Role.Name)
			}
		}
		token, err := s.tokenService.GenerateToken(&tdto)
		if err != nil {
			return nil, err
		}
		return token, nil

	}
	// Register method
	bp := []byte(common.GeneratePassword())
	hp, err := bcrypt.GenerateFromPassword(bp, bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error(logging.General, logging.HashPassword, err.Error(), nil)
		return nil, err
	}
	u.Password = string(hp)
	roleId, err := s.getDefaultRole()
	if err != nil {
		s.logger.Error(logging.Postgres, logging.DefaultRoleNotFound, err.Error(), nil)
		return nil, err
	}
	// Start of transaction
	tx := s.database.Begin()
	err = tx.Create(&u).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Rollback, err.Error(), nil)
		return nil, err
	}
	err = tx.Create(&models.UserRole{RoleId: roleId, UserId: u.Id}).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Rollback, err.Error(), nil)
		return nil, err
	}
	tx.Commit()

	var user models.User
	err = s.database.Model(&models.User{}).
		Where("username = ?", u.Username).
		Preload("UserRoles", func(tx *gorm.DB) *gorm.DB {
			return tx.Preload("Role")
		}).
		Find(&user).Error
	if err != nil {
		return nil, err
	}
	tdto := tokenDto{UserId: user.Id, FirstName: user.FirstName, LastName: user.LastName,
		Email: user.Email, MobileNumber: user.MobileNumber}
	if len(*user.UserRoles) > 0 {
		for _, userRole := range *user.UserRoles {
			tdto.Roles = append(tdto.Roles, userRole.Role.Name)
		}
	}
	token, err := s.tokenService.GenerateToken(&tdto)
	if err != nil {
		return nil, err
	}
	return token, nil

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
