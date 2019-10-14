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
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.HandleFunc("/home", homeHandle)
	http.HandleFunc("/start", startHandle)
	http.HandleFunc("/pipes", pipesHandle)
	http.HandleFunc("/register", registerHandle)
	http.HandleFunc("/size", sizeHandle)
	http.HandleFunc("/get", getHandle)
	http.HandleFunc("/game", gameHandle)

	log.Println("Start server on port 12301")
	log.Fatal(http.ListenAndServe(":12301", nil))
}

func gameHandle(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/get.html")
}

func homeHandle(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/register.html")
}

func startHandle(w http.ResponseWriter, r *http.Request) {
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

func pipesHandle(w http.ResponseWriter, r *http.Request) {
	response := make([]string, len(pipes))
	for i, _ := range pipes {
		response = append(response, string(i))
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(500)
		_, _ = fmt.Fprint(w, err)
		return
	}

	_, err = w.Write(responseJson)
	if err != nil {
		w.WriteHeader(500)
		_, _ = fmt.Fprint(w, err)
		return
	}
}

func registerHandle(w http.ResponseWriter, r *http.Request) {
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

func sizeHandle(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, os.Getenv("AREA_SIZE"))
	if err != nil {
		log.Fatal(err)
	}
}

func getHandle(w http.ResponseWriter, r *http.Request) {
	if _, ok := r.URL.Query()["id"]; !ok {
		w.WriteHeader(400)
		_, _ = fmt.Fprint(w, "id query param must be exist")
		return
	}

	id, err := strconv.ParseInt(r.URL.Query()["id"][0], 10, 32)
	if err != nil {
		w.WriteHeader(500)
		_, _ = fmt.Fprint(w, "id parse error")
		return
	}

	// @todo give pipes different names like alpha and other
	if len(pipes) <= int(id) {
		w.WriteHeader(404)
		_, _ = fmt.Fprintf(w, string("Not found"))
		return
	}
	pipe := pipes[id]

	// @todo we need buffer to return current state of game thought it.
	// Because now every user pop values from pipe and other users can't see this values
	jsonResponse, err := json.Marshal(<-pipe)
	if err != nil {
		w.WriteHeader(500)
		_, _ = fmt.Fprint(w, err)
		return
	}

	_, _ = fmt.Fprintf(w, string(jsonResponse))
}
