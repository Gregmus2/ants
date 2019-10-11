package internal

import (
	"ants/internal/global"
	"ants/pkg"
	"errors"
	"io"
	"io/ioutil"
	"math"
	"os"
	"strconv"
)

var storage global.Storage
var pipes []chan [][]string

func init() {
	storage = global.NewBolt("ants")
	pipes = make([]chan [][]string, 0)
}

// @todo separate by methods
func start(names []string) (int, error) {
	if len(names) != 2 {
		return 0, errors.New("names must be 2")
	}

	users := make([]*global.User, len(names))
	for i := 0; i < len(names); i++ {
		users[i] = global.LoadUser(storage, names[i])
	}

	size64, err := strconv.ParseInt(os.Getenv("AREA_SIZE"), 10, 64)
	size := int(size64)
	usize := uint(size)
	if err != nil {
		return 0, err
	}

	quartSize := uint(math.Round(float64(size / 4)))
	halfSize := uint(math.Round(float64(size / 2)))

	// @todo more ants, more automation
	ants := make([]*global.Ant, len(users))
	ants[0] = &global.Ant{
		Pos:    [2]uint{quartSize, halfSize},
		User:   users[0],
		IsDead: false,
	}
	ants[1] = &global.Ant{
		Pos:    [2]uint{usize - quartSize, halfSize},
		User:   users[1],
		IsDead: false,
	}

	area := make([][]*global.Object, size)
	for x := 0; x < size; x++ {
		area[x] = make([]*global.Object, size)
		for y := 0; y < size; y++ {
			if x == 0 || x == size-1 || y == 0 || y == size-1 {
				area[x][y] = global.CreateWall()
			} else {
				area[x][y] = global.CreateEmptyObject()
			}
		}
	}

	area[ants[0].Pos.X()][ants[0].Pos.Y()] = global.CreateAnt(ants[0])
	area[ants[1].Pos.X()][ants[1].Pos.Y()] = global.CreateAnt(ants[1])

	foodCount := int(float64(size) * 0.5)
	if foodCount < len(ants) {
		foodCount = len(ants)
	}
	for i := 0; i < foodCount; i += 2 {
		x := global.Random.Intn(int(halfSize))
		y := global.Random.Intn(size)
		if area[x][y].Type == pkg.AntField {
			x = global.Random.Intn(int(halfSize))
			y = global.Random.Intn(size)
		}
		area[x][y] = global.CreateFood()
		area[x+int(halfSize)][y] = global.CreateFood()
	}

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
