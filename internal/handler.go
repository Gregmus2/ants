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
var matches map[string]*game.Match

func init() {
	storage = global.NewBolt("ants")
	matches = make(map[string]*game.Match)
}

func prepareGame(names []string) (string, error) {
	users := make([]*global.User, 0, len(names))
	for i := 0; i < len(names); i++ {
		users = append(users, global.LoadUser(storage, names[i]))
	}

	builder, err := game.NewMatchBuilder(global.Config.AreaSize, users)
	if err != nil {
		return "", err
	}

	builder.BuildAnts()
	builder.BuildArea()
	builder.BuildFood(0.05, 0.07, len(names), true)

	id := strconv.Itoa(global.Random.Intn(1000))
	matches[id] = builder.BuildMatch(storage)
	go matches[id].Run(id)

	return id, nil
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

	algorithmsPath := os.ExpandEnv("$GOPATH/src/ants/algorithms")
	aFile, err := os.OpenFile(algorithmsPath+"/"+name+".so", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer aFile.Close()

	_, err = aFile.WriteAt(fileBytes, 0)
	if err != nil {
		return err
	}

	return err
}
