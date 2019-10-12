package internal

import (
	"ants/internal/game"
	"ants/internal/global"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

var storage global.Storage
var pipes []chan [][]string

func init() {
	storage = global.NewBolt("ants")
	pipes = make([]chan [][]string, 0)
}

func prepareGame(names []string) (int, error) {
	users := make([]*global.User, len(names))
	for i := 0; i < len(names); i++ {
		users[i] = global.LoadUser(storage, names[i])
	}

	size64, err := strconv.ParseInt(os.Getenv("AREA_SIZE"), 10, 64)
	size := int(size64)
	if err != nil {
		return 0, err
	}

	builder, err := game.NewMatchBuilder(size, users)
	if err != nil {
		return 0, err
	}

	builder.BuildAnts()
	builder.BuildArea()
	builder.BuildFood(0.05, 0.07, len(names), true)

	pipe := make(chan [][]string, 100)
	// @todo we need to recreate pipes for some time, because of length grow
	pipes = append(pipes, pipe)
	go builder.BuildMatch().Run(pipe)

	return len(pipes) - 1, nil
}

func registration(name string, color string, algorithmFile io.Reader) error {
	err := saveAlgorithmFile(algorithmFile, name)
	if err != nil {
		return err
	}

	global.CreateUser(name, color, storage)
	return nil
}

func saveAlgorithmFile(file io.Reader, name string) error {
	// read all of the contents of our uploaded file into a byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	// Create a file within our algorithm directory that follows a particular naming pattern
	err = ioutil.WriteFile("algorithms/"+name+".so", fileBytes, 0644)
	if err != nil {
		return err
	}

	return err
}
