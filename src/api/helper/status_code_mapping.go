package helper

import (
	"net/http"

	"github.com/ehsan-hosseiny/golang-web-api/pkg/service_errors"
)

var StatusCodeMapping = map[string]int{

	// OTP
	service_errors.OtpExists:  409,
	service_errors.OtpUsed:    409,
	service_errors.OtpInvalid: 400,
}

func TranslateErrorToStatusCode(err error) int {
	value, ok := StatusCodeMapping[err.Error()]
	if !ok {
		return http.StatusInternalServerError
	}
	return value
}