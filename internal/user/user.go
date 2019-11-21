package user

import (
	"encoding/json"
	"errors"
	"io"
	"log"

	pkg "github.com/gregmus2/ants-pkg"
)

type User struct {
	Name      string
	Color     string
	algorithm pkg.Algorithm
}

func (u *User) Algorithm() pkg.Algorithm {
	return u.algorithm
}

func (s *Service) CreateUser(name string, color string) {
	user := &User{
		Name:      name,
		Color:     color,
		algorithm: nil,
	}

	s.storage.CreateCollectionIfNotExist(Collection)

	s.SaveUser(user)
}

func (s *Service) SaveUser(u *User) {
	data, err := json.Marshal(u)
	if err != nil {
		log.Print(err)
		return
	}

	err = s.storage.Put(Collection, u.Name, data)
	if err != nil {
		log.Print(err)
		return
	}
}

func (s *Service) LoadUser(name string) (*User, error) {
	data, err := s.storage.Get(Collection, name)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, errors.New("user not found")
	}

	alg, err := s.LoadAlgorithm(name)
	if err != nil {
		return nil, err
	}

	user := &User{algorithm: alg}
	err = json.Unmarshal(data, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) Register(name string, color string, algorithmFile io.Reader) error {
	err := s.SaveCodeFile(algorithmFile, name)
	if err != nil {
		return err
	}

	_, err = s.LoadAlgorithm(name)
	if err != nil {
		return err
	}

	s.CreateUser(name, color)
	return nil
}

func (s *Service) GetUsersByNames(names []string) ([]*User, error) {
	users := make([]*User, 0, len(names))
	for i := 0; i < len(names); i++ {
		usr, err := s.LoadUser(names[i])
		if err != nil {
			return nil, err
		}

		users = append(users, usr)
	}

	return users, nil
}