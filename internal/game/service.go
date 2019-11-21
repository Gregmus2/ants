package game

import (
	"ants/internal/config"
	"ants/internal/storage"
	"ants/internal/user"
	"math/rand"
	"strconv"
)

type Service struct {
	storage     storage.Storage
	config      *config.Config
	userService *user.Service
	matches     map[string]*Match
}

func NewService(storage storage.Storage, config *config.Config, userService *user.Service) *Service {
	return &Service{
		storage:     storage,
		config:      config,
		userService: userService,
		matches:     make(map[string]*Match),
	}
}

func (s *Service) RunGame(names []string) (string, error) {
	users, err := s.userService.GetUsersByNames(names)
	if err != nil {
		return "", err
	}

	id := strconv.Itoa(rand.Intn(1000))
	state, err := NewMatchState(s.config.AreaSize, users)
	if err != nil {
		return "", err
	}

	BuildArea(state)
	BuildAnts(state)
	BuildFood(state, 0.01, 0.03, len(names), true)

	s.matches[id] = CreateMatch(s, state, id)
	go s.matches[id].Run()

	return id, nil
}