package internal

import (
	"ants/internal/game"
	"ants/internal/global"
	"errors"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
)

var storage global.Storage
var matches map[string]*game.Match

func init() {
	storage = global.NewBolt("ants")
	matches = make(map[string]*game.Match)
}

// todo wrap errors with github.com/pkg/errors
func prepareGame(names []string) (string, error) {
	users := make([]*global.User, 0, len(names))
	for i := 0; i < len(names); i++ {
		user, err := global.LoadUser(storage, names[i])
		if err != nil {
			return "", err
		}

		users = append(users, user)
	}

	id := strconv.Itoa(rand.Intn(1000))
	b, err := game.NewMatchBuilder(id, global.Config.AreaSize, users)
	if err != nil {
		return "", err
	}

	b.BuildArea()
	b.BuildAnts()
	b.BuildFood(0.01, 0.03, len(names), true)

	matches[id] = b.BuildMatch(storage)
	go matches[id].Run()

	return id, nil
}

func registration(name string, color string, algorithmFile io.Reader) error {
	err := saveCodeFile(algorithmFile, name)
	if err != nil {
		return err
	}

	_, err = global.LoadAlgorithm(name)
	if err != nil {
		return err
	}

	global.CreateUser(name, color, storage)
	return nil
}

func saveCodeFile(file io.Reader, name string) error {
	// read all of the contents of our uploaded file into a byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	codePath := global.Config.BasePath + "/algorithms/" + name + ".go"
	aFile, err := os.OpenFile(codePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	_, err = aFile.WriteAt(fileBytes, 0)
	if err != nil {
		return err
	}

	err = aFile.Close()
	if err != nil {
		return err
	}

	outputPath := global.Config.BasePath + "/algorithms/" + name + ".so"
	cmd := exec.Command("/usr/local/go/bin/go", "build", "-buildmode=plugin", "-o", outputPath, codePath)

	s, err := cmd.Output()
	if err != nil {
		return errors.New(err.Error() + string(s))
	}

	return os.Remove(codePath)
}
