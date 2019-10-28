package internal

import (
	"ants/internal/global"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func Serve() {
	http.HandleFunc("/api/start", startHandle)
	http.HandleFunc("/api/pipes", pipesHandle)
	http.HandleFunc("/api/players", playersHandle)
	http.HandleFunc("/api/register", registerHandle)
	http.HandleFunc("/api/size", sizeHandle)
	http.HandleFunc("/api/get", getHandle)

	log.Println("Start server on port 12301")
	log.Fatal(http.ListenAndServe(":12301", nil))
}

func startHandle(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		badRequest(w, err)
		return
	}

	namesString := r.PostFormValue("names")
	if namesString == "" {
		badRequest(w, "names have blank values")
		return
	}

	n := strings.Split(namesString, ",")
	id, err := prepareGame(n)
	if err != nil {
		badRequest(w, err)
		return
	}

	_, _ = fmt.Fprint(w, id)
}

func pipesHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	response := make([]string, 0, len(matches))
	for name := range matches {
		response = append(response, name)
	}

	res, err := json.Marshal(response)
	if err != nil {
		serverError(w, err)
		return
	}

	_, err = w.Write(res)
	if err != nil {
		serverError(w, err)
		return
	}
}

func playersHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	names, err := global.GetNames(storage)
	if err != nil {
		serverError(w, err)
		return
	}

	res, err := json.Marshal(names)
	if err != nil {
		serverError(w, err)
		return
	}

	_, err = w.Write(res)
	if err != nil {
		serverError(w, err)
		return
	}
}

func registerHandle(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		badRequest(w, err)
		return
	}

	name := r.FormValue("name")
	color := r.FormValue("color")
	if name == "" || color == "" {
		badRequest(w, "name or color have blank values")
		return
	}

	file, _, err := r.FormFile("algorithm")
	if err != nil {
		badRequest(w, "Error Retrieving the File: "+err.Error())
		return
	}
	defer file.Close()

	err = registration(name, color, file)
	if err != nil {
		serverError(w, err)
		return
	}
}

func sizeHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	_, err := fmt.Fprint(w, strconv.Itoa(global.Config.AreaSize))
	if err != nil {
		log.Fatal(err)
	}
}

func getHandle(w http.ResponseWriter, r *http.Request) {
	_, okID := r.URL.Query()["id"]
	_, okPart := r.URL.Query()["part"]
	if !okID || !okPart {
		badRequest(w, "id, part query param must be exist")
		return
	}
	id := r.URL.Query()["id"][0]
	part := r.URL.Query()["part"][0]

	// todo give pipes different names like alpha and other
	match, ok := matches[id]
	if !ok {
		notFound(w)
		return
	}

	res, err := json.Marshal(match.LoadRound(id, part))
	if err != nil {
		serverError(w, err)
		return
	}

	_, _ = fmt.Fprint(w, string(res))
}

func serverError(w http.ResponseWriter, msg interface{}) {
	w.WriteHeader(500)
	_, _ = fmt.Fprint(w, msg)
}

func notFound(w http.ResponseWriter) {
	w.WriteHeader(404)
	_, _ = fmt.Fprintf(w, "Not found")
}

func badRequest(w http.ResponseWriter, msg interface{}) {
	w.WriteHeader(400)
	_, _ = fmt.Fprint(w, msg)
}