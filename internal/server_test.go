package internal

import (
	"ants/internal/global"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestServe(t *testing.T) {
	if len(pipesTestRequest(t)) != 0 {
		t.Error("pipes must be empty by start")
	}

	if sizeTestRequest(t) != strconv.Itoa(global.Config.AreaSize) {
		t.Error("size endpoint must return size from config")
	}

	registerTestRequest(t, "Greg", "blue")
	registerTestRequest(t, "Greg2", "green")

	players := playersTestRequest(t)
	if len(players) != 2 || players[0] != "Greg" || players[1] != "Greg2" {
		t.Error("wrong players")
	}

	id := startTestRequest(t)
	if id == "" {
		t.Error("empty id from start request")
	}

	time.Sleep(1000 * time.Millisecond)
	area := getTestRequest(t, id)
	if len(area) != global.Config.MatchPartSize {
		t.Error("wrong batch size")
	}

	_ = os.Remove("ants.db")
}

func pipesTestRequest(t *testing.T) []string {
	ts := httptest.NewServer(http.HandlerFunc(pipesHandle))
	defer ts.Close()

	res, err := ts.Client().Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	pipes := make([]string, 0)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(body, &pipes)
	if err != nil {
		t.Error(err)
	}

	return pipes
}

func playersTestRequest(t *testing.T) []string {
	ts := httptest.NewServer(http.HandlerFunc(playersHandle))
	defer ts.Close()

	res, err := ts.Client().Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	players := make([]string, 0)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(body, &players)
	if err != nil {
		t.Error(err)
	}

	return players
}

func sizeTestRequest(t *testing.T) string {
	ts := httptest.NewServer(http.HandlerFunc(sizeHandle))
	defer ts.Close()

	res, err := ts.Client().Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}

	return string(body)
}

func registerTestRequest(t *testing.T, name string, color string) {
	ts := httptest.NewServer(http.HandlerFunc(registerHandle))
	defer ts.Close()

	file, _ := os.Open("../test/" + name + ".go")
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

	res, err := ts.Client().Post(ts.URL, writer.FormDataContentType(), body)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 200 {
		t.Error("Wrong status code: " + strconv.Itoa(res.StatusCode))
		body, _ := ioutil.ReadAll(res.Body)
		t.Error("Body: " + string(body))
	}
}

func startTestRequest(t *testing.T) string {
	ts := httptest.NewServer(http.HandlerFunc(startHandle))
	defer ts.Close()

	form := url.Values{}
	form.Add("names", "Greg,Greg2")

	res, err := ts.Client().Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
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

	return string(body)
}

func getTestRequest(t *testing.T, id string) [][][]string {
	ts := httptest.NewServer(http.HandlerFunc(getHandle))
	defer ts.Close()

	params := url.Values{}
	params.Add("id", id)
	params.Add("part", "1")

	res, err := ts.Client().Get(ts.URL + "?" + params.Encode())
	if err != nil {
		t.Fatal(err)
	}

	area := make([][][]string, 0)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(body, &area)
	if err != nil {
		t.Error(err)
	}

	return area
}
