package main

import (
	"ants/internal"
	"ants/internal/config"
	"ants/internal/game"
	"ants/internal/storage"
	"ants/internal/user"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	cfg := config.NewConfig()
	s := storage.NewBolt("ants")
	userService := user.NewService(s, cfg)
	gameService := game.NewService(s, cfg, userService)
	infoService := config.NewService(s, cfg)

	server := internal.NewJSONServer()
	server.Handle("/api/pipes", gameService.MatchNamesAction)
	server.Handle("/api/size", infoService.SizeAction)
	server.Handle("/api/players", userService.PlayersAction)
	server.Handle("/api/register", userService.RegistrationAction)
	server.Handle("/api/start", gameService.StartAction)
	server.Handle("/api/get", gameService.GetMatchAction)

	server.Start(12301)
}