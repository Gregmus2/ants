package internal

import (
	"ants/internal/global"
	"errors"
	"io"
	"io/ioutil"
)

var storage global.Storage
var pipes []chan [][]string

func init() {
	storage = global.NewBolt("ants")
	pipes = make([]chan [][]string, 0)
}

func start(names []string) (int, error) {
	if len(names) != 2 {
		return 0, errors.New("names must be 2")
	}

	users := make([]*global.User, len(names))
	for i := 0; i < len(names); i++ {
		users[i] = global.LoadUser(storage, names[i])
	}

	ants := make([]*global.Ant, len(users))
	ants[0] = &global.Ant{
		Pos:    [2]uint{2, 5},
		User:   users[0],
		IsDead: false,
	}
	ants[1] = &global.Ant{
		Pos:    [2]uint{7, 5},
		User:   users[1],
		IsDead: false,
	}

	area := make([][]*global.Object, 10)
	for x := 0; x < 10; x++ {
		area[x] = make([]*global.Object, 10)
		for y := 0; y < 10; y++ {
			if x == 0 || x == 9 || y == 0 || y == 9 {
				area[x][y] = global.CreateWall()
			} else {
				area[x][y] = global.CreateEmptyObject()
			}
		}
	}

	area[2][5] = global.CreateAnt(ants[0])
	area[7][5] = global.CreateAnt(ants[1])

	game := CreateGame(users, ants, area)
	pipe := make(chan [][]string, 100)
	pipes = append(pipes, pipe)
	go game.Run(pipe)

	return len(pipes) - 1, nil
}

func register(name string, color string, algorithmFile io.Reader) error {
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
