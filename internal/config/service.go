package config

import (
	"ants/internal/storage"
)

type Service struct {
	storage storage.Storage
	config  *Config
}

func NewService(storage storage.Storage, config *Config) *Service {
	return &Service{storage: storage, config: config}
}
