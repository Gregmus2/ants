package user

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"plugin"
	"strings"

	pkg "github.com/gregmus2/ants-pkg"
)

func (s *Service) LoadAlgorithm(name string) (pkg.Algorithm, error) {
	p := s.config.BasePath + "/algorithms/" + name + ".so"
	plug, err := plugin.Open(p)
	if err != nil {
		return nil, err
	}

	sym, err := plug.Lookup(strings.Title(name))
	if err != nil {
		return nil, err
	}

	var algorithm pkg.Algorithm
	algorithm, ok := sym.(pkg.Algorithm)
	if !ok {
		return nil, errors.New("wrong symbol")
	}

	return algorithm, nil
}

func (s *Service) SaveCodeFile(file io.Reader, name string) error {
	// read all of the contents of our uploaded file into a byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	codePath := s.config.BasePath + "/algorithms/" + name + ".go"
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

	outputPath := s.config.BasePath + "/algorithms/" + name + ".so"
	cmd := exec.Command("/usr/local/go/bin/go", "build", "-buildmode=plugin", "-o", outputPath, codePath)

	out, err := cmd.Output()
	if err != nil {
		return errors.New(err.Error() + string(out))
	}

	return os.Remove(codePath)
}
