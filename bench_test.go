package main

import (
	"ants/internal/config"
	"ants/internal/game"
	"ants/internal/storage"
	"ants/internal/user"
	"math/rand"
	"os"
	"strconv"
	"testing"
)

func BenchmarkMatch(b *testing.B) {
	rand.Seed(666)

	cfg := config.NewConfig()
	s := storage.NewBolt("bench")
	userService := user.NewService(s, cfg)
	gameService := game.NewService(s, cfg, userService)

	registerUser(b, userService, "Greg", "blue")
	registerUser(b, userService, "Greg2", "red")

	names := []string{"Greg", "Greg2"}
	users, err := userService.GetUsersByNames(names)
	if err != nil {
		b.Fatal(err)
	}

	id := strconv.Itoa(rand.Intn(1000))
	state, err := game.NewMatchState(cfg.AreaSize, users)
	if err != nil {
		b.Fatal(err)
	}

	game.BuildArea(state)
	game.BuildAnts(state)
	game.BuildFood(state, 0.01, 0.03, len(names), true)
	match := game.CreateMatch(gameService, state, id)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		match.Run()
	}
	b.StopTimer()

	_ = os.Remove("bench.db")
}

func registerUser(b *testing.B, userService *user.Service, name string, color string) {
	file, err := os.Open("./testdata/bench/" + name + ".go")
	if err != nil {
		b.Fatal(err)
	}

	err = userService.Register(name, color, file)
	if err != nil {
		b.Fatal(err)
	}
}
