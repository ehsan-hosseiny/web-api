package service_errors

type ServiceError struct {
	EndUserMessage   string `json:"endUserMessage"`
	TechniclaMessage string `json:"technicalMessage"`
	Err              error
}

func (s *ServiceError) Error() string {
	return s.EndUserMessage
}
