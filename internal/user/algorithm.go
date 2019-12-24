package user

import (
	user "ants/internal/util"
	"errors"
	"io"
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
	codePath := s.config.BasePath + "/algorithms/" + name
	err := user.Unzip(file, codePath)
	if err != nil {
		return err
	}

	outputPath := s.config.BasePath + "/algorithms/" + name + ".so"

	cmd := exec.Command("/usr/local/go/bin/go", "build", "-buildmode=plugin", "-o", outputPath)
	cmd.Dir = codePath

	out, err := cmd.Output()
	if err != nil {
		return errors.New(err.Error() + string(out))
	}

	return os.RemoveAll(codePath)
}
