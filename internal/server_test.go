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
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestServe(t *testing.T) {
	if len(pipesTestRequest(t)) != 0 {
		t.Error("pipes must be empty by start")
	}

	if sizeTestRequest(t) != strconv.Itoa(global.Config.AreaSize) {
		t.Error("size endpoint must return size from config")
	}

	registerTestRequest(t)
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

func registerTestRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(registerHandle))
	defer ts.Close()

	file, _ := os.Open("../algorithms/test.go")
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

	err = writer.WriteField("name", "Greg")
	if err != nil {
		t.Fatal(err)
	}

	err = writer.WriteField("color", "#000000")
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
