package global

import (
	"encoding/json"
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

func LoadUser(storage Storage, name string) *User {
	data, err := storage.Get(UserCollection, name)
	if err != nil {
		log.Print(err)
		return nil
	}

	if data == nil {
		return nil
	}

	user := &User{storage: storage, algorithm: loadAlgorithm(name)}
	err = json.Unmarshal(data, user)
	if err != nil {
		log.Print(err)
		return nil
	}

	return user
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

func loadAlgorithm(name string) pkg.Algorithm {
	path := "./algorithms/" + name + ".so"
	plug, err := plugin.Open(path)
	if err != nil {
		log.Println(err)
		return nil
	}

	symbol, err := plug.Lookup(strings.Title(name))
	if err != nil {
		log.Println(err)
		return nil
	}

	var algorithm pkg.Algorithm
	algorithm, ok := symbol.(pkg.Algorithm)
	if !ok {
		log.Println("Wrong symbol")
		return nil
	}

	return algorithm
}
