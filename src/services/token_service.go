package services

import (
	"time"

	"github.com/ehsan-hosseiny/golang-web-api/api/dto"
	"github.com/ehsan-hosseiny/golang-web-api/config"
	"github.com/ehsan-hosseiny/golang-web-api/constants"
	"github.com/ehsan-hosseiny/golang-web-api/pkg/logging"
	"github.com/ehsan-hosseiny/golang-web-api/pkg/service_errors"
	"github.com/golang-jwt/jwt"
)

type TokenService struct {
	logger logging.Logger
	cfg    *config.Config
}

type tokenDto struct {
	UserId       int
	FirstName    string
	LastName     string
	UserName     string
	MobileNumber string
	Email        string
	Roles        []string
}

func NewTokenService(cfg *config.Config) *TokenService {

	logger := logging.NewLogger(cfg)
	return &TokenService{
		cfg:    cfg,
		logger: logger,
	}
}

func (s *TokenService) GenerateToken(token *tokenDto) (*dto.TokenDetail, error) {
	td := &dto.TokenDetail{}
	td.AccessTokenExpireTime = time.Now().Add(s.cfg.JWT.AccessTokenExpireTimeDuration * time.Minute).Unix()
	td.RefreshTokenExpireTime = time.Now().Add(s.cfg.JWT.RefreshTokenExpireTimeDuration * time.Minute).Unix()

	// generate access token
	accessTokenClaims := jwt.MapClaims{}
	accessTokenClaims[constants.UserIdKey] = token.UserId
	accessTokenClaims[constants.FirstNameKey] = token.FirstName
	accessTokenClaims[constants.LastNameKey] = token.LastName
	accessTokenClaims[constants.UsernameKey] = token.UserName
	accessTokenClaims[constants.EmailKey] = token.Email
	accessTokenClaims[constants.MobileNumberKey] = token.MobileNumber
	accessTokenClaims[constants.RolesKey] = token.Roles
	accessTokenClaims[constants.ExpireTimeKey] = td.AccessTokenExpireTime

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	var err error
	td.AccessToken, err = accessToken.SignedString([]byte(s.cfg.JWT.Secret))

	if err != nil {
		return nil, err
	}

	// generate refresh token
	refreshTokenClaims := jwt.MapClaims{}
	refreshTokenClaims["user_id"] = token.UserId
	refreshTokenClaims["exp"] = td.AccessTokenExpireTime

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	td.RefreshToken, err = rt.SignedString([]byte(s.cfg.JWT.RefreshSecret))

	if err != nil {
		return nil, err
	}
	return td, nil
}

func (s *TokenService) VerifyToken(token string) (*jwt.Token, error) {
	ac, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, &service_errors.ServiceError{EndUserMessage: service_errors.UnExpectedError}
		}
		return []byte(s.cfg.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	return ac, nil
}

// get user info by this method like user id , user name, email ...
func (s *TokenService) GetClaims(token string) (claimMap map[string]interface{}, err error) {
	claimMap = map[string]interface{}{}

	verifyToken, err := s.VerifyToken(token)
	if err != nil {
		return nil, err
	}
	claims, ok := verifyToken.Claims.(jwt.MapClaims)
	if ok && verifyToken.Valid {
		for k, v := range claims {
			claimMap[k] = v
		}
		return claimMap, nil
	}
	return nil, &service_errors.ServiceError{EndUserMessage: service_errors.ClaimNotFound}
}
