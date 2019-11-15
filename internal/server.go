package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type JSONServer struct{}

func NewJSONServer() *JSONServer {
	return &JSONServer{}
}

func (s *JSONServer) Handle(path string, handler func(r *http.Request) (interface{}, int)) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		response, code := handler(r)

		res, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(w, err)
			return
		}

		w.WriteHeader(code)
		_, err = w.Write(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(w, err)
			return
		}
	})
}

func (s *JSONServer) Start(port uint16) {
	log.Printf("Start server on port %d", port)
	addr := fmt.Sprintf(":%d", port)
	err := http.ListenAndServe(addr, nil)
	log.Fatal(err)
}