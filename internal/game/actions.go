package game

import (
	"net/http"
	"strings"
)

func (s *Service) MatchNamesAction(r *http.Request) (interface{}, int) {
	names := make([]string, 0, len(s.matches))
	for name := range s.matches {
		names = append(names, name)
	}

	return names, http.StatusOK
}

func (s *Service) StartAction(r *http.Request) (interface{}, int) {
	err := r.ParseForm()
	if err != nil {
		return err, http.StatusBadRequest
	}

	namesString := r.PostFormValue("names")
	if namesString == "" {
		return "names have blank values", http.StatusBadRequest
	}

	n := strings.Split(namesString, ",")
	id, err := s.RunGame(n)
	if err != nil {
		return err, http.StatusBadRequest
	}

	return id, http.StatusOK
}

func (s *Service) GetMatchAction(r *http.Request) (interface{}, int) {
	_, okID := r.URL.Query()["id"]
	_, okPart := r.URL.Query()["part"]
	if !okID || !okPart {
		return "id, part query param must be exist", http.StatusBadRequest
	}
	id := r.URL.Query()["id"][0]
	part := r.URL.Query()["part"][0]

	// todo give pipes different names like alpha and other
	match, ok := s.matches[id]
	if !ok {
		return nil, http.StatusNotFound
	}

	return match.LoadRound(id, part), http.StatusOK
}