package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func Serve() {
	http.HandleFunc("/start", start)
	http.HandleFunc("/register", register)
	http.HandleFunc("/size", size)
	http.HandleFunc("/get", get)

	err := http.ListenAndServe(":12301", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func start(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")

	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}

	namesString := r.FormValue("names")
	if namesString == "" {
		w.WriteHeader(400)
		_, _ = fmt.Fprint(w, "names have blank values")
		return
	}

	names := strings.Split(namesString, ",")
	pipeNum, err := prepareGame(names)
	if err != nil {
		w.WriteHeader(400)
		_, _ = fmt.Fprint(w, err)
		return
	}

	_, _ = fmt.Fprint(w, pipeNum)
}

func register(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}

	name := r.FormValue("name")
	color := r.FormValue("color")
	if name == "" || color == "" {
		w.WriteHeader(400)
		_, _ = fmt.Fprint(w, "name or color have blank values")
		return
	}

	file, _, err := r.FormFile("algorithm")
	if err != nil {
		w.WriteHeader(400)
		_, _ = fmt.Fprint(w, "Error Retrieving the File")
		_, _ = fmt.Fprint(w, err)
		return
	}
	defer file.Close()

	err = registration(name, color, file)
	if err != nil {
		w.WriteHeader(500)
		_, _ = fmt.Fprint(w, err)
		return
	}
}

func size(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")

	_, err := fmt.Fprintf(w, os.Getenv("AREA_SIZE"))
	if err != nil {
		log.Fatal(err)
	}
}

func get(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")

	id, err := strconv.ParseInt(r.URL.Query()["id"][0], 10, 32)
	if err != nil {
		w.WriteHeader(400)
		_, _ = fmt.Fprint(w, "id query param must be exist")
		return
	}

	pipe := pipes[id]
	if pipe == nil {
		w.WriteHeader(400)
		_, _ = fmt.Fprint(w, "Wrong channel")
		return
	}

	jsonResponse, err := json.Marshal(<-pipe)
	if err != nil {
		w.WriteHeader(500)
		_, _ = fmt.Fprint(w, err)
		return
	}

	_, _ = fmt.Fprintf(w, string(jsonResponse))
}
