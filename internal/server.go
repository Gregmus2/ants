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
	http.HandleFunc("/api/register", registerHandle)
	http.HandleFunc("/api/size", sizeHandle)
	http.HandleFunc("/api/get", getHandle)

	log.Println("Start server on port 12301")
	log.Fatal(http.ListenAndServe(":12301", nil))
}

func startHandle(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}

	namesString := r.PostFormValue("names")
	if namesString == "" {
		w.WriteHeader(400)
		_, _ = fmt.Fprint(w, "names have blank values")
		return
	}

	names := strings.Split(namesString, ",")
	id, err := prepareGame(names)
	if err != nil {
		w.WriteHeader(400)
		_, _ = fmt.Fprint(w, err)
		return
	}

	_, _ = fmt.Fprint(w, id)
}

func pipesHandle(w http.ResponseWriter, r *http.Request) {
	response := make([]string, 0, len(matches))
	for name, _ := range matches {
		response = append(response, name)
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
	_, err := fmt.Fprintf(w, strconv.Itoa(global.Config.AreaSize))
	if err != nil {
		log.Fatal(err)
	}
}

func getHandle(w http.ResponseWriter, r *http.Request) {
	_, okId := r.URL.Query()["id"]
	_, okPart := r.URL.Query()["part"]
	if !okId || !okPart {
		w.WriteHeader(400)
		_, _ = fmt.Fprint(w, "id, part query param must be exist")
		return
	}
	id := r.URL.Query()["id"][0]
	part := r.URL.Query()["part"][0]

	// @todo give pipes different names like alpha and other
	match, ok := matches[id]
	if !ok {
		w.WriteHeader(404)
		_, _ = fmt.Fprintf(w, "Not found")
		return
	}

	jsonResponse, err := json.Marshal(match.LoadRound(id, part))
	if err != nil {
		w.WriteHeader(500)
		_, _ = fmt.Fprint(w, err)
		return
	}

	_, _ = fmt.Fprintf(w, string(jsonResponse))
}
