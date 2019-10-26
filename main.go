package main

import (
	"ants/internal"
	"ants/internal/global"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	global.InitConfig()

	internal.Serve()
}