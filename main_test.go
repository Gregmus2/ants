package main

import (
	"ants/internal"
	"ants/internal/config"
	"ants/internal/game"
	"ants/internal/storage"
	"ants/internal/user"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

const basePath string = "http://127.0.0.1:12302"

func setup() *config.Config {
	rand.Seed(time.Now().UTC().UnixNano())

	cfg := config.NewConfig()
	s := storage.NewBolt("test")
	userService := user.NewService(s, cfg)
	gameService := game.NewService(s, cfg, userService)
	infoService := config.NewService(s, cfg)

	server := internal.NewJSONServer()
	server.Handle("/api/pipes", gameService.MatchNamesAction)
	server.Handle("/api/size", infoService.SizeAction)
	server.Handle("/api/players", userService.PlayersAction)
	server.Handle("/api/register", userService.RegistrationAction)
	server.Handle("/api/start", gameService.StartAction)
	server.Handle("/api/get", gameService.GetMatchAction)

	go server.Start(12302)

	return cfg
}

func TestServe(t *testing.T) {
	cfg := setup()
	time.Sleep(1000 * time.Millisecond)

	pipes := make([]string, 0)
	JSONDecode(t, get(t, basePath+"/api/pipes"), &pipes)
	if len(pipes) != 0 {
		t.Error("pipes must be empty by start")
	}

	size := string(get(t, basePath+"/api/size"))
	if size != strconv.Itoa(cfg.AreaSize) {
		t.Error("size endpoint must return size from config")
	}

	registrationTestRequest(t, "Greg", "blue")
	registrationTestRequest(t, "Greg2", "green")

	players := make([]string, 0)
	JSONDecode(t, get(t, basePath+"/api/players"), &players)
	if len(players) != 2 || players[0] != "Greg" || players[1] != "Greg2" {
		t.Error("wrong players")
	}

	var id string
	JSONDecode(t, startTestRequest(t), &id)
	if id == "" {
		t.Error("empty id from start request")
	}

	time.Sleep(1000 * time.Millisecond)
	area := getTestRequest(t, id)
	if len(area) != cfg.Match.PartSize {
		t.Errorf("wrong batch size %d", len(area))
	}

	_ = os.Remove("test.db")
}

func JSONDecode(t *testing.T, body []byte, data interface{}) {
	err := json.Unmarshal(body, data)
	if err != nil {
		t.Error(err, string(body))
	}
}

func get(t *testing.T, path string) []byte {
	res, err := (&http.Client{}).Get(path)
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}

	return body
}

func registrationTestRequest(t *testing.T, name string, color string) {
	file, _ := os.Open("./testdata/" + name + ".zip")
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("algorithm", filepath.Base(file.Name()))
	if err != nil {
		t.Fatal(err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatal(err)
	}

	err = writer.WriteField("name", name)
	if err != nil {
		t.Fatal(err)
	}

	err = writer.WriteField("color", color)
	if err != nil {
		t.Fatal(err)
	}

	err = writer.Close()
	if err != nil {
		t.Fatal(err)
	}

	res, err := (&http.Client{}).Post(basePath+"/api/register", writer.FormDataContentType(), body)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusCreated {
		t.Error("Wrong status code: " + strconv.Itoa(res.StatusCode))
		body, _ := ioutil.ReadAll(res.Body)
		t.Error("Body: " + string(body))
	}
}

func startTestRequest(t *testing.T) []byte {
	form := url.Values{}
	form.Add("names", "Greg,Greg2")

	contentType := "application/x-www-form-urlencoded"
	res, err := (&http.Client{}).Post(basePath+"/api/start", contentType, strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 200 {
		t.Error("Wrong status code: " + strconv.Itoa(res.StatusCode))
		body, _ := ioutil.ReadAll(res.Body)
		t.Error("Body: " + string(body))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}

	return body
}

func getTestRequest(t *testing.T, id string) [][][]string {
	params := url.Values{}
	params.Add("id", id)
	params.Add("part", "1")

	res, err := (&http.Client{}).Get(basePath + "/api/get?" + params.Encode())
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 200 {
		t.Error("Wrong status code: " + strconv.Itoa(res.StatusCode))
		body, _ := ioutil.ReadAll(res.Body)
		t.Error("Body: " + string(body))
	}

	area := make([][][]string, 0)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(body, &area)
	if err != nil {
		t.Error(err, string(body))
	}

	return area
}
