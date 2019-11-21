package user

import "net/http"

func (s *Service) PlayersAction(r *http.Request) (interface{}, int) {
	players, err := s.storage.GetKeys(Collection)
	if err != nil {
		return err.Error(), http.StatusInternalServerError
	}

	return players, http.StatusOK
}

func (s *Service) RegistrationAction(r *http.Request) (interface{}, int) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return err.Error(), http.StatusBadRequest
	}

	name := r.FormValue("name")
	color := r.FormValue("color")
	if name == "" || color == "" {
		return "name or color have blank values", http.StatusBadRequest
	}

	file, _, err := r.FormFile("algorithm")
	if err != nil {
		return "Error Retrieving the File: " + err.Error(), http.StatusBadRequest
	}
	defer file.Close()

	err = s.Register(name, color, file)
	if err != nil {
		return err.Error(), http.StatusInternalServerError
	}

	return nil, http.StatusCreated
}