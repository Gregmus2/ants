package user

import (
	"ants/internal/config"
	"ants/internal/storage"
)

type Service struct {
	storage storage.Storage
	config  *config.Config
}

const Collection string = "Users"

func NewService(storage storage.Storage, config *config.Config) *Service {
	return &Service{storage: storage, config: config}
}