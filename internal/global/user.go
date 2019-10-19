package global

import (
	"encoding/json"
	"errors"
	pkg "github.com/gregmus2/ants-pkg"
	"log"
	"plugin"
	"strings"
)

type User struct {
	Name      string
	Color     string
	algorithm pkg.Algorithm
	storage   Storage
}

func (u *User) Algorithm() pkg.Algorithm {
	return u.algorithm
}

const UserCollection string = "Users"

func CreateUser(name string, color string, storage Storage) {
	user := &User{
		Name:      name,
		Color:     color,
		algorithm: nil,
		storage:   storage,
	}

	storage.CreateCollectionIfNotExist(UserCollection)

	user.Save()
}

func LoadUser(storage Storage, name string) (*User, error) {
	data, err := storage.Get(UserCollection, name)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, errors.New("user not found")
	}

	alg, err := LoadAlgorithm(name)
	if err != nil {
		return nil, err
	}

	user := &User{storage: storage, algorithm: alg}
	err = json.Unmarshal(data, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetNames(storage Storage) ([]string, error) {
	return storage.GetKeys(UserCollection)
}

func (u *User) Save() {
	data, err := json.Marshal(u)
	if err != nil {
		log.Print(err)
		return
	}
	err = u.storage.Put(UserCollection, u.Name, data)
	if err != nil {
		log.Print(err)
		return
	}
}

func LoadAlgorithm(name string) (pkg.Algorithm, error) {
	path := Config.BasePath + "/algorithms/" + name + ".so"
	plug, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}

	symbol, err := plug.Lookup(strings.Title(name))
	if err != nil {
		return nil, err
	}

	var algorithm pkg.Algorithm
	algorithm, ok := symbol.(pkg.Algorithm)
	if !ok {
		return nil, errors.New("wrong symbol")
	}

	return algorithm, nil
}
