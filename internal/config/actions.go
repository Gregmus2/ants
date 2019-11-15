package config

import "net/http"

func (s *Service) SizeAction(r *http.Request) (interface{}, int) {
	return s.config.AreaSize, http.StatusOK
}
